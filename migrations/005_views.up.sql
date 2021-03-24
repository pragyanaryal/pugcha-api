-- Materialized views of the businesses
CREATE MATERIALIZED VIEW businesses_view
AS
SELECT b.user_id,
       b.approved_by,
       b.blocked,
       b.approved,
       b.created_on,
       b.updated_on,
       b.id       AS business_id,
       c.id       AS category_id,
       c.name     AS category_name,
       bp.name    AS business_name,
       bp.picture as picture,
       bp.pan_number,
       bp.email,
       bp.website,
       bp.vat_number,
       bp.contact,
       bp.location,
       bp.established_date,
       bp.description,
       bp.geom,
       jsonb_agg(DISTINCT a2) as address,
       jsonb_agg(DISTINCT oi) as owners,
       jsonb_agg(DISTINCT ot) as opening
FROM businesses b
         INNER JOIN business_profiles bp on b.id = bp.business_id
         INNER JOIN categories c on c.id = bp.category_id
         LEFT JOIN owners_info oi on b.id = oi.business_id
         LEFT JOIN opening_times ot on b.id = ot.business_id
         LEFT JOIN addresses a2 on b.id = a2.business_id
GROUP BY b.id, b.user_id, b.approved_by, b.blocked, b.approved, b.created_on, b.updated_on, c.id, c.name, bp.name,
         bp.picture, bp.pan_number, bp.email, bp.website, bp.vat_number, bp.contact, bp.location, bp.established_date,
         bp.description, bp.geom
WITH DATA;

-- Materialized view of the users
CREATE MATERIALIZED VIEW users_view
AS
SELECT m.id,
       m.email,
       m.status,
       m.password,
       m.type,
       m.created_on,
       m.updated_on,
       up.first_name,
       up.last_name,
       up.gender,
       up.dob,
       up.middle_name,
       up.contact,
       json_strip_nulls(
               json_build_object('first_name', up.first_name, 'last_name', up.last_name, 'middle_name', up.middle_name,
                                 'dob', up.dob, 'gender', up.gender, 'contact', up.contact, 'picture',
                                 up.profile_pics)) AS profile
FROM users m
         LEFT JOIN users_profile up on m.id = up.user_id
GROUP BY m.id, m.email, m.status, m.password, m.type, m.created_on, m.updated_on, up.first_name, up.last_name,
         up.middle_name, up.dob, up.gender, up.contact, up.profile_pics
WITH DATA;

CREATE UNIQUE INDEX idx_users_view_id ON users_view (id);
CREATE UNIQUE INDEX idx_users_view_email ON users_view (email);
CREATE INDEX idx_users_view_status ON users_view (status);
CREATE INDEX idx_users_view_type ON users_view (type);
CREATE INDEX idx_users_view_gender ON users_view (gender);

CREATE UNIQUE INDEX idx_business_id ON businesses_view (business_id);
CREATE INDEX idx_business_user_id ON businesses_view (user_id);
CREATE INDEX idx_business_approved_by_2 ON businesses_view (approved_by);
CREATE INDEX idx_business_category ON businesses_view (category_id);
CREATE INDEX idx_business_geom ON businesses_view USING GIST (geom);

