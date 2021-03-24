CREATE EXTENSION IF NOT EXISTS postgis;

SELECT ADDGEOMETRYCOLUMN('business_profiles', 'geom', 6207, 'POINT', 2);

CREATE INDEX idx_geom ON business_profiles (geom);