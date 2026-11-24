-- Uptime for one check over [p_from, p_to), optionally limited to one
-- network.
--
-- Exact at the edges: whole hours come from check_rollups_hourly, and the
-- partial hours at either end (plus anything after the rollup watermark)
-- come from check_results.
--
-- Note: raw results older than retention.raw_days are purged, so partial
-- edge hours that far back have no data. Use hour-aligned ranges for long
-- look-backs.

CREATE OR REPLACE FUNCTION get_uptime(
    p_check_id   bigint,
    p_from       timestamptz,
    p_to         timestamptz,
    p_network_id bigint DEFAULT NULL
)
RETURNS TABLE (
    total_count    bigint,
    success_count  bigint,
    uptime_pct     numeric
)
LANGUAGE plpgsql
STABLE
AS $$
#variable_conflict use_column
DECLARE
    v_hour_from timestamptz;
    v_hour_to   timestamptz;
BEGIN
    -- First whole hour at or after p_from
    v_hour_from := date_trunc('hour', p_from, 'UTC');
    IF v_hour_from < p_from THEN
        v_hour_from := v_hour_from + interval '1 hour';
    END IF;

    -- End of the last whole hour that has been rolled up
    v_hour_to := LEAST(date_trunc('hour', p_to, 'UTC'), rollup_watermark());
    IF v_hour_to < v_hour_from THEN
        v_hour_to := v_hour_from;   -- no rollup span; raw covers everything
    END IF;

    RETURN QUERY
    WITH parts AS (
        -- Whole hours from rollups: [v_hour_from, v_hour_to)
        SELECT sum(r.total_count)   AS total,
               sum(r.success_count) AS ok
          FROM check_rollups_hourly r
         WHERE r.check_id = p_check_id
           AND (p_network_id IS NULL OR r.network_id = p_network_id)
           AND r.bucket_start >= v_hour_from
           AND r.bucket_start <  v_hour_to

        UNION ALL

        -- Everything else from raw results
        SELECT count(*),
               count(*) FILTER (WHERE cr.success)
          FROM check_results cr
         WHERE cr.check_id = p_check_id
           AND (p_network_id IS NULL OR cr.network_id = p_network_id)
           AND cr.checked_at >= p_from
           AND cr.checked_at <  p_to
           AND (cr.checked_at < v_hour_from OR cr.checked_at >= v_hour_to)
    )
    SELECT COALESCE(sum(p.total), 0)::bigint,
           COALESCE(sum(p.ok), 0)::bigint,
           round(100.0 * sum(p.ok) / NULLIF(sum(p.total), 0), 3)
      FROM parts p;
END;
$$;
