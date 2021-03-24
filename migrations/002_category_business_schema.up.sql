-- Create Category Table
CREATE TABLE categories
(
    id         uuid PRIMARY KEY,
    name       varchar(30) UNIQUE,
    picture    varchar(150),
    created_by uuid REFERENCES users (id),
    created_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create Business Table
CREATE TABLE businesses
(
    id          uuid PRIMARY KEY,
    approved    bool,
    user_id     uuid REFERENCES users (id),
    approved_by uuid REFERENCES users (id),
    blocked     bool,
    created_on  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_business_approved ON businesses (approved);
CREATE INDEX idx_business_approved_by ON businesses (approved_by);
CREATE INDEX idx_business_blocked ON businesses (blocked);

CREATE TYPE geo_loc AS
(
    latitude  float8,
    longitude float8
);

CREATE TABLE IF NOT EXISTS owners_info
(
    id           uuid PRIMARY KEY,
    business_id  uuid REFERENCES businesses (id),
    user_id      uuid REFERENCES users (id),
    profile_pics varchar(150),
    email        varchar(50),
    contact      char(10)[]
);

CREATE TABLE IF NOT EXISTS addresses
(
    id           uuid PRIMARY KEY,
    business_id  uuid REFERENCES businesses (id),
    street       varchar(50),
    ward         int,
    municipality varchar(30),
    district     varchar(20),
    state        varchar(20),
    country      varchar(30),
    contact      char(10)[]
);

CREATE TABLE IF NOT EXISTS opening_times
(
    business_id  uuid REFERENCES businesses (id),
    week_day     int,
    opening_time time with time zone,
    closing_time time with time zone,
    opened       bool,
    PRIMARY KEY (business_id, week_day)
);

CREATE TABLE business_profiles
(
    business_id      uuid PRIMARY KEY REFERENCES businesses (id),
    category_id      uuid REFERENCES categories (id),
    user_id          uuid REFERENCES users (id),
    name             varchar(80),
    pan_number       varchar(15),
    email            varchar(50),
    website          varchar(70),
    vat_number       varchar(15),
    picture          varchar(150)[],
    contact          char(10)[],
    location         geo_loc,
    established_date date,
    description      text,
    created_on       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_business_profile_category ON business_profiles (category_id);
CREATE INDEX idx_business_profile_user ON business_profiles (user_id);
CREATE INDEX idx_business_profile_contact ON business_profiles (contact);
CREATE INDEX idx_business_profile_location ON business_profiles (location);
CREATE INDEX idx_business_profile_website ON business_profiles (website);
