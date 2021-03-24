DROP INDEX IF EXISTS idx_business_profile_category;
DROP INDEX IF EXISTS idx_business_profile_contact;
DROP INDEX IF EXISTS idx_business_profile_location;
DROP INDEX IF EXISTS idx_business_profile_website;
DROP INDEX IF EXISTS idx_business_profile_address;

DROP TABLE IF EXISTS business_profiles;

DROP TYPE IF EXISTS address;
DROP TYPE IF EXISTS opening_hours;
DROP TYPE IF EXISTS owner_info;
DROP TYPE IF EXISTS geo_loc;

DROP INDEX IF EXISTS idx_business_approved;
DROP INDEX IF EXISTS idx_business_approved_by;
DROP INDEX IF EXISTS idx_business_blocked;

DROP TABLE IF EXISTS businesses;

DROP TABLE IF EXISTS categories;