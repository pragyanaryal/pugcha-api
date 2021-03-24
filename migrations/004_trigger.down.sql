DROP TRIGGER IF EXISTS update_geom_trigger on "public"."business_profiles";

DROP FUNCTION IF EXISTS update_geom();

DROP EXTENSION IF EXISTS "uuid-ossp";
