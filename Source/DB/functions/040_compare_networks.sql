-- Compares one check across every network it has run on over
-- [p_from, p_to). Mainly for global checks (internet / DNS baselines),
-- e.g. "latency to 1.1.1.1 at home vs. the office".
--
-- Hour-granular (see hourly_check_stats). latency_p95 is a weighted
-- average of hourly p95s, an approximation.

CREATE OR REPLACE FUNCTION compare_networks(
    p_check_id bigint,
    p_from     timestamptz,
    p_to       timestamptz
)
RETURNS TABLE (
    network_id      bigint,
    network_name    text,
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
    WITH stats AS (
        SELECT h.network_id                          AS nid,
               sum(h.total_count)                    AS total,
               sum(h.success_count)                  AS ok,
               sum(h.latency_avg * h.success_count)  AS lat_sum,
               sum(h.latency_p95 * h.success_count)  AS p95_sum
          FROM hourly_check_stats(ARRAY[p_check_id], p_from, p_to) h
         GROUP BY h.network_id
    ),
    inc AS (
        SELECT i.network_id AS nid,
               count(*)     AS cnt
          FROM incidents i
         WHERE i.check_id = p_check_id
           AND i.started_at >= p_from
           AND i.started_at <  p_to
         GROUP BY i.network_id
    )
    SELECT net.id,
           net.name,
           s.total::bigint,
           s.ok::bigint,
           round(100.0 * s.ok / NULLIF(s.total, 0), 3),
           (s.lat_sum / NULLIF(s.ok, 0)::double precision)::real,
           (s.p95_sum / NULLIF(s.ok, 0)::double precision)::real,
           COALESCE(inc.cnt, 0)
      FROM stats s
      JOIN networks net ON net.id = s.nid
      LEFT JOIN inc     ON inc.nid = s.nid
     ORDER BY net.name;
END;
$$;
