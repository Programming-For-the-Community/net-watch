-- Latency chart points for one check over [p_from, p_to), grouped into
-- buckets of p_bucket (e.g. '5 minutes', '1 hour', '1 day').
--
-- Buckets under an hour are computed from raw results, so they only reach
-- back as far as retention.raw_days. Buckets of an hour or more use
-- hourly_check_stats (rollups + recent raw), so they reach back as far as
-- retention.rollup_days.
--
-- For buckets of an hour or more, latency_avg is weighted by successful
-- samples, and latency_p95 is a weighted average of hourly p95s (an
-- approximation; exact percentiles can't be combined from rollups).
--
-- p_bucket must not contain months or years (a date_bin restriction).

CREATE OR REPLACE FUNCTION get_latency_series(
    p_check_id   bigint,
    p_from       timestamptz,
    p_to         timestamptz,
    p_bucket     interval DEFAULT interval '5 minutes',
    p_network_id bigint   DEFAULT NULL
)
RETURNS TABLE (
    bucket_start   timestamptz,
    sample_count   bigint,
    success_count  bigint,
    latency_min    real,
    latency_avg    real,
    latency_max    real,
    latency_p95    real
)
LANGUAGE plpgsql
STABLE
AS $$
#variable_conflict use_column
DECLARE
    c_origin CONSTANT timestamptz := '2000-01-01 00:00:00+00';
BEGIN
    IF p_bucket < interval '1 hour' THEN
        RETURN QUERY
        SELECT date_bin(p_bucket, cr.checked_at, c_origin) AS b,
               count(*),
               count(*) FILTER (WHERE cr.success),
               min(cr.latency_ms) FILTER (WHERE cr.success),
               (avg(cr.latency_ms) FILTER (WHERE cr.success))::real,
               max(cr.latency_ms) FILTER (WHERE cr.success),
               (percentile_cont(0.95) WITHIN GROUP (ORDER BY cr.latency_ms)
                    FILTER (WHERE cr.success))::real
          FROM check_results cr
         WHERE cr.check_id = p_check_id
           AND (p_network_id IS NULL OR cr.network_id = p_network_id)
           AND cr.checked_at >= p_from
           AND cr.checked_at <  p_to
         GROUP BY b
         ORDER BY b;
    ELSE
        RETURN QUERY
        SELECT date_bin(p_bucket, h.bucket_start, c_origin) AS b,
               sum(h.total_count)::bigint,
               sum(h.success_count)::bigint,
               min(h.latency_min),
               (sum(h.latency_avg * h.success_count)
                    / NULLIF(sum(h.success_count), 0)::double precision)::real,
               max(h.latency_max),
               (sum(h.latency_p95 * h.success_count)
                    / NULLIF(sum(h.success_count), 0)::double precision)::real
          FROM hourly_check_stats(ARRAY[p_check_id], p_from, p_to) h
         WHERE (p_network_id IS NULL OR h.network_id = p_network_id)
         GROUP BY b
         ORDER BY b;
    END IF;
END;
$$;
