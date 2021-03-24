DROP TRIGGER IF EXISTS opening_trigger on "public"."opening_times";

DROP TRIGGER IF EXISTS owners_trigger on "public"."owners_info";

DROP TRIGGER IF EXISTS address_trigger on "public"."addresses";

DROP TRIGGER IF EXISTS business_profile_trigger on "public"."business_profiles";

DROP TRIGGER IF EXISTS business_trigger on "public"."businesses";

DROP FUNCTION IF EXISTS send_notify_business();