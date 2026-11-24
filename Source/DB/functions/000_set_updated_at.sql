-- Trigger function: keeps updated_at current on UPDATE.
-- Attached by the CREATE TRIGGER statements in the table files, so this
-- file must run before the tables folder.
--
-- Skips no-op updates (row unchanged) so updated_at only moves when data
-- actually changes.

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW IS DISTINCT FROM OLD THEN
        NEW.updated_at := now();
    END IF;
    RETURN NEW;
END;
$$;
