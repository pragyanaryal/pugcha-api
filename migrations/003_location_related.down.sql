DROP INDEX idx_geom;

ALTER TABLE business_profiles DROP COLUMN geom;

DROP EXTENSION IF EXISTS postgis;