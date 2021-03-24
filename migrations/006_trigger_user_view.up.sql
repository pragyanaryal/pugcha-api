CREATE OR REPLACE FUNCTION send_notify_user() RETURNS trigger
    LANGUAGE PLPGSQL AS
$$
BEGIN
    PERFORM pg_notify('changes', 'users'::text);
    RETURN NULL;
END;
$$;

CREATE TRIGGER user_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON users
    FOR EACH STATEMENT
EXECUTE PROCEDURE send_notify_user();

CREATE TRIGGER user_profile_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON users_profile
    FOR EACH STATEMENT
EXECUTE PROCEDURE send_notify_user();
