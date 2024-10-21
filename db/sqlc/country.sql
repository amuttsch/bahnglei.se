-- name: SaveCountry :one
INSERT INTO countries (created_at, updated_at, iso_code, name, osm_url) 
  VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $1, $2, $3) 
  ON CONFLICT (iso_code) DO UPDATE SET name = $2, osm_url = $3
  RETURNING *;

-- name: CountCountries :one
select count(*) from countries;

-- name: GetCountries :many
select c.*, count(s.id) stations from countries c left join stations s on s.country_iso_code = c.iso_code group by c.iso_code;

