-- Checks the poller should run now on the given network.
--
-- A check is due when:
--   * it is enabled
--   * neither the check nor its device is paused
--   * its device is on p_network_id, or is global (network_id IS NULL)
--   * it has never run, or its interval has elapsed since the last run
--
-- Never-run checks come first, then the most overdue.

CREATE OR REPLACE FUNCTION get_due_checks(
    p_network_id bigint,
    p_limit      int DEFAULT 100
)
RETURNS TABLE (
    check_id            bigint,
    device_id           bigint,
    network_id          bigint,
    address             text,
    type                text,
    config              jsonb,
    timeout_ms          int,
    interval_seconds    int,
    failure_threshold   int,
    recovery_threshold  int,
    last_checked_at     timestamptz
)
LANGUAGE plpgsql
STABLE
AS $$
#variable_conflict use_column
BEGIN
    RETURN QUERY
    SELECT c.id,
           d.id,
           d.network_id,
           d.address,
           c.type,
           c.config,
           c.timeout_ms,
           c.interval_seconds,
           c.failure_threshold,
           c.recovery_threshold,
           s.last_checked_at
      FROM checks c
      JOIN devices d           ON d.id = c.device_id
      LEFT JOIN check_status s ON s.check_id = c.id
     WHERE c.enabled
       AND (c.paused_until IS NULL OR c.paused_until <= now())
       AND (d.paused_until IS NULL OR d.paused_until <= now())
       AND (d.network_id = p_network_id OR d.network_id IS NULL)
       AND (s.last_checked_at IS NULL
            OR s.last_checked_at + make_interval(secs => c.interval_seconds) <= now())
     ORDER BY s.last_checked_at NULLS FIRST
     LIMIT p_limit;
END;
$$;
