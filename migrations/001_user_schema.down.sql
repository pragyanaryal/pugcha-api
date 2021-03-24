/** Drop Social Account Tables **/
DROP TABLE IF EXISTS facebook_profile CASCADE;
DROP TABLE IF EXISTS google_profile CASCADE;

/** Drop User Profile Tables and Indexes **/
DROP INDEX IF EXISTS id_user_profile_email;
DROP INDEX IF EXISTS id_user_profile_id;

DROP TABLE IF EXISTS users_profile;

DROP TYPE IF EXISTS gender;

/** Drop User Table and Indexes **/
DROP INDEX IF EXISTS id_user_email;
DROP INDEX IF EXISTS id_user_id;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS user_type;
DROP TYPE IF EXISTS user_status;