-- Deletes data past its retention period, in batches, committing after
-- each batch so locks stay short and the WAL doesn't balloon.
--
-- Retention comes from settings (defaults if missing):
--   retention.raw_days     check_results            (default 30)
--   retention.rollup_days  check_rollups_hourly     (default 365)
--
-- Raw results that haven't been rolled up yet (at or after the rollup
-- watermark) are never deleted, so a stalled rollup can't lose data.
--
-- Run rollup_hourly() first. Call outside an explicit transaction:
--   CALL purge_old_data();

CREATE OR REPLACE PROCEDURE purge_old_data(
    p_batch_size int DEFAULT 10000
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_raw_days     int;
    v_rollup_days  int;
    v_raw_cutoff   timestamptz;
    v_roll_cutoff  timestamptz;
    v_deleted      int;
BEGIN
    SELECT COALESCE(max((s.value #>> '{}')::int) FILTER (WHERE s.key = 'retention.raw_days'), 30),
           COALESCE(max((s.value #>> '{}')::int) FILTER (WHERE s.key = 'retention.rollup_days'), 365)
      INTO v_raw_days, v_rollup_days
      FROM settings s
     WHERE s.key IN ('retention.raw_days', 'retention.rollup_days');

    v_raw_cutoff  := LEAST(now() - make_interval(days => v_raw_days), rollup_watermark());
    v_roll_cutoff := now() - make_interval(days => v_rollup_days);

    -- Raw results
    LOOP
        DELETE FROM check_results
         WHERE id IN (SELECT cr.id
                        FROM check_results cr
                       WHERE cr.checked_at < v_raw_cutoff
                       LIMIT p_batch_size);
        GET DIAGNOSTICS v_deleted = ROW_COUNT;
        COMMIT;
        EXIT WHEN v_deleted < p_batch_size;
    END LOOP;

    -- Hourly rollups (no id column, so batch by ctid)
    LOOP
        DELETE FROM check_rollups_hourly
         WHERE ctid = ANY (ARRAY(SELECT r.ctid
                                   FROM check_rollups_hourly r
                                  WHERE r.bucket_start < v_roll_cutoff
                                  LIMIT p_batch_size));
        GET DIAGNOSTICS v_deleted = ROW_COUNT;
        COMMIT;
        EXIT WHEN v_deleted < p_batch_size;
    END LOOP;
END;
$$;
