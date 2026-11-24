-- Rolls raw check_results up into check_rollups_hourly, one hour at a time,
-- committing after each hour, then advances the rollup watermark.
--
-- Only complete hours are rolled up. Re-running over the same range is
-- safe (upsert).
--
-- Defaults:
--   p_from  the current watermark, or the oldest raw result on first run
--   p_to    the start of the current hour
--
-- Call outside an explicit transaction (COMMIT inside a procedure is only
-- allowed then):  CALL rollup_hourly();

CREATE OR REPLACE PROCEDURE rollup_hourly(
    p_from timestamptz DEFAULT NULL,
    p_to   timestamptz DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_from  timestamptz;
    v_to    timestamptz;
    v_hour  timestamptz;
BEGIN
    v_to := date_trunc('hour', COALESCE(p_to, now()), 'UTC');

    v_from := p_from;
    IF v_from IS NULL THEN
        v_from := rollup_watermark();
        IF v_from = '-infinity'::timestamptz THEN
            SELECT min(cr.checked_at) INTO v_from FROM check_results cr;
        END IF;
    END IF;

    IF v_from IS NULL THEN
        RETURN;   -- no data yet
    END IF;

    v_hour := date_trunc('hour', v_from, 'UTC');

    WHILE v_hour < v_to LOOP
        INSERT INTO check_rollups_hourly AS r (
            check_id, network_id, bucket_start,
            total_count, success_count,
            latency_min, latency_avg, latency_max, latency_p95
        )
        SELECT cr.check_id,
               cr.network_id,
               v_hour,
               count(*),
               count(*) FILTER (WHERE cr.success),
               min(cr.latency_ms) FILTER (WHERE cr.success),
               (avg(cr.latency_ms) FILTER (WHERE cr.success))::real,
               max(cr.latency_ms) FILTER (WHERE cr.success),
               (percentile_cont(0.95) WITHIN GROUP (ORDER BY cr.latency_ms)
                    FILTER (WHERE cr.success))::real
          FROM check_results cr
         WHERE cr.checked_at >= v_hour
           AND cr.checked_at <  v_hour + interval '1 hour'
         GROUP BY cr.check_id, cr.network_id
        ON CONFLICT (check_id, network_id, bucket_start) DO UPDATE
           SET total_count   = EXCLUDED.total_count,
               success_count = EXCLUDED.success_count,
               latency_min   = EXCLUDED.latency_min,
               latency_avg   = EXCLUDED.latency_avg,
               latency_max   = EXCLUDED.latency_max,
               latency_p95   = EXCLUDED.latency_p95;

        v_hour := v_hour + interval '1 hour';

        -- Advance the watermark (never move it backwards on a re-run)
        INSERT INTO settings (key, value)
        VALUES ('rollup.watermark', to_jsonb(v_hour))
        ON CONFLICT (key) DO UPDATE
           SET value = EXCLUDED.value
         WHERE (settings.value #>> '{}')::timestamptz < v_hour;

        COMMIT;
    END LOOP;
END;
$$;
