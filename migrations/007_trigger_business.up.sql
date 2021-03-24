CREATE OR REPLACE FUNCTION send_notify_business() RETURNS trigger
    LANGUAGE PLPGSQL AS
$$
BEGIN
    PERFORM pg_notify('changes', 'business'::text);
    RETURN NULL;
END;
$$;

CREATE TRIGGER business_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON businesses
    FOR EACH STATEMENT
EXECUTE PROCEDURE send_notify_business();

CREATE TRIGGER business_profile_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON business_profiles
    FOR EACH STATEMENT
EXECUTE PROCEDURE send_notify_business();

CREATE TRIGGER address_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON addresses
    FOR EACH STATEMENT
EXECUTE PROCEDURE send_notify_business();

CREATE TRIGGER owners_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON owners_info
    FOR EACH STATEMENT
EXECUTE PROCEDURE send_notify_business();

CREATE TRIGGER opening_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON opening_times
    FOR EACH STATEMENT
EXECUTE PROCEDURE send_notify_business();