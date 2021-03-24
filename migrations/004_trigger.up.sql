CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_geom() RETURNS trigger
    LANGUAGE PLPGSQL AS
$$
BEGIN
    UPDATE business_profiles
    SET geom = ST_SetSRID(ST_MakePoint((NEW.location).longitude, (NEW.location).latitude), 6207)
    WHERE business_id = NEW.business_id;
    RETURN NEW;
END;
$$;

CREATE TRIGGER update_geom_trigger
    AFTER INSERT OR UPDATE of location
    ON "business_profiles"
    FOR EACH ROW
EXECUTE PROCEDURE update_geom();