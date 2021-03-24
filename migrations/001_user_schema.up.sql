/** Create User **/
CREATE TYPE user_status AS ENUM ('active', 'blocked', 'approval_needed', 'verification_needed');
CREATE TYPE user_type AS ENUM ('user', 'admin', 'super');

CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY,
    email      varchar(50) unique,
    status     user_status              NOT NULL DEFAULT 'verification_needed',
    password   varchar,
    type       user_type                NOT NULL DEFAULT 'user',
    created_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON users (id);
CREATE INDEX idx_user_email ON users (email);


/** Create User Profile **/
CREATE TYPE gender AS ENUM ('male', 'female', 'rest', '');

CREATE TABLE IF NOT EXISTS users_profile
(
    user_id      uuid PRIMARY KEY REFERENCES users (id),
    email        varchar(50) unique,
    first_name   varchar(20),
    middle_name  varchar(20),
    last_name    varchar(20),
    dob          date,
    gender       gender,
    contact      char(10),
    profile_pics varchar(150),
    created_on   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_profile_id ON users_profile (user_id);
CREATE INDEX idx_user_profile_email ON users_profile (email);


/** Create Google Profile **/
CREATE TABLE IF NOT EXISTS google_profile
(
    user_id     uuid PRIMARY KEY REFERENCES users (id),
    email       varchar(50) unique,
    google_id   varchar unique,
    access_key  varchar,
    refresh_key varchar
);

/** Create Facebook Profile **/
CREATE TABLE IF NOT EXISTS facebook_profile
(
    user_id     uuid PRIMARY KEY REFERENCES users (id),
    email       varchar(50) unique,
    facebook_id varchar unique,
    access_key  varchar,
    refresh_key varchar
);