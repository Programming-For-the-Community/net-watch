-- Hourly stats for a set of checks, per network, merging the two sources:
--   * check_rollups_hourly for hours before the rollup watermark
--   * check_results (aggregated on the fly) for hours after it
--
-- Building block for get_latency_series, compare_networks and
-- get_network_health. Hour-granular: the first bucket starts at the top of
-- the hour containing p_from.
--
-- Latency figures only include successful results.
-- Buckets are aligned to UTC hours regardless of the session time zone.

CREATE OR REPLACE FUNCTION hourly_check_stats(
    p_check_ids bigint[],
    p_from      timestamptz,
    p_to        timestamptz
)
RETURNS TABLE (
    check_id       bigint,
    network_id     bigint,
    bucket_start   timestamptz,
    total_count    bigint,
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
    v_watermark timestamptz := rollup_watermark();
    v_from      timestamptz := date_trunc('hour', p_from, 'UTC');
BEGIN
    RETURN QUERY
    -- Hours already rolled up
    SELECT r.check_id,
           r.network_id,
           r.bucket_start,
           r.total_count::bigint,
           r.success_count::bigint,
           r.latency_min,
           r.latency_avg,
           r.latency_max,
           r.latency_p95
      FROM check_rollups_hourly r
     WHERE r.check_id = ANY (p_check_ids)
       AND r.bucket_start >= v_from
       AND r.bucket_start <  LEAST(p_to, v_watermark)

    UNION ALL

    -- Hours not rolled up yet
    SELECT cr.check_id,
           cr.network_id,
           date_trunc('hour', cr.checked_at, 'UTC'),
           count(*),
           count(*) FILTER (WHERE cr.success),
           min(cr.latency_ms) FILTER (WHERE cr.success),
           (avg(cr.latency_ms) FILTER (WHERE cr.success))::real,
           max(cr.latency_ms) FILTER (WHERE cr.success),
           (percentile_cont(0.95) WITHIN GROUP (ORDER BY cr.latency_ms)
                FILTER (WHERE cr.success))::real
      FROM check_results cr
     WHERE cr.check_id = ANY (p_check_ids)
       AND cr.checked_at >= GREATEST(v_from, v_watermark)
       AND cr.checked_at <  p_to
     GROUP BY cr.check_id, cr.network_id, date_trunc('hour', cr.checked_at, 'UTC');
END;
$$;
