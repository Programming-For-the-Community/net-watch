-- Returns the end of the last hour that rollup_hourly() has fully rolled
-- up. Everything before it is in check_rollups_hourly; everything from it
-- onward must be read from check_results.
--
-- Stored in settings under 'rollup.watermark' (written by rollup_hourly).
-- Returns -infinity before the first rollup, meaning "use raw data only".

CREATE OR REPLACE FUNCTION rollup_watermark()
RETURNS timestamptz
LANGUAGE plpgsql
STABLE
AS $$
DECLARE
    v_watermark timestamptz;
BEGIN
    SELECT (s.value #>> '{}')::timestamptz
      INTO v_watermark
      FROM settings s
     WHERE s.key = 'rollup.watermark';

    RETURN COALESCE(v_watermark, '-infinity'::timestamptz);
END;
$$;
