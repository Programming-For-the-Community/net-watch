-- Baseline health of one network over [p_from, p_to): one row per check on
-- the baseline devices that apply to it:
--   * the network's own gateway device   (system_role = 'gateway')
--   * the global internet / DNS devices  (system_role = 'internet' / 'dns')
-- Results are limited to those taken on p_network_id.
--
-- Baseline checks with no data in the range are still returned (zeros).
-- Hour-granular (see hourly_check_stats).
--
-- The system_role values are defined in Go (model.SystemRole); keep the
-- literals below in sync.

CREATE OR REPLACE FUNCTION get_network_health(
    p_network_id bigint,
    p_from       timestamptz,
    p_to         timestamptz
)
RETURNS TABLE (
    system_role     text,
    device_id       bigint,
    check_id        bigint,
    check_name      text,
    sample_count    bigint,
    success_count   bigint,
    uptime_pct      numeric,
    latency_avg     real,
    latency_p95     real,
    incident_count  bigint
)
LANGUAGE plpgsql
STABLE
AS $$
#variable_conflict use_column
BEGIN
    RETURN QUERY
    WITH baseline AS (
        SELECT d.system_role AS role,
               d.id          AS did,
               c.id          AS cid,
               c.name        AS cname
          FROM devices d
          JOIN checks c ON c.device_id = d.id
         WHERE (d.system_role = 'gateway' AND d.network_id = p_network_id)
            OR (d.system_role IN ('internet', 'dns') AND d.network_id IS NULL)
    ),
    stats AS (
        SELECT h.check_id                            AS cid,
               sum(h.total_count)                    AS total,
               sum(h.success_count)                  AS ok,
               sum(h.latency_avg * h.success_count)  AS lat_sum,
               sum(h.latency_p95 * h.success_count)  AS p95_sum
          FROM hourly_check_stats(ARRAY(SELECT b.cid FROM baseline b), p_from, p_to) h
         WHERE h.network_id = p_network_id
         GROUP BY h.check_id
    ),
    inc AS (
        SELECT i.check_id AS cid,
               count(*)   AS cnt
          FROM incidents i
         WHERE i.network_id = p_network_id
           AND i.started_at >= p_from
           AND i.started_at <  p_to
         GROUP BY i.check_id
    )
    SELECT b.role,
           b.did,
           b.cid,
           b.cname,
           COALESCE(s.total, 0)::bigint,
           COALESCE(s.ok, 0)::bigint,
           round(100.0 * s.ok / NULLIF(s.total, 0), 3),
           (s.lat_sum / NULLIF(s.ok, 0)::double precision)::real,
           (s.p95_sum / NULLIF(s.ok, 0)::double precision)::real,
           COALESCE(inc.cnt, 0)
      FROM baseline b
      LEFT JOIN stats s ON s.cid = b.cid
      LEFT JOIN inc     ON inc.cid = b.cid
     ORDER BY CASE b.role WHEN 'gateway' THEN 1 WHEN 'internet' THEN 2 ELSE 3 END,
              b.cname;
END;
$$;
