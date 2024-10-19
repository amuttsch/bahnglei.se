-- name: CreateStation :one
insert into stations (id, country_iso_code, name, lat, lng, operator, wikidata, wikipedia, tracks) 
  values ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
  returning *;
 
-- name: UpdateStationNumberOfTracks :exec
update stations 
  set tracks = $2
  where id = $1;

-- name: GetStation :one
select * from stations where id = $1;

-- name: DeleteStation :exec
delete from stations where id = $1;

-- name: CountStations :one
select count(*) from stations;

-- name: SearchStations :many
select * from stations where name ILIKE $1 order by tracks desc limit 20;

