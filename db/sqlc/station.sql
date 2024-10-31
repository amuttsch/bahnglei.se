-- name: CreateStation :one
insert into stations (id, country_iso_code, name, coordinate, operator, wikidata, wikipedia, tracks) 
  values ($1, $2, $3, $4, $5, $6, $7, $8) 
  on conflict (id) do update set name = $3, coordinate = $4, operator = $5, wikidata = $6, wikipedia = $7, tracks = $8
  returning *;
 
-- name: UpdateStationNumberOfTracks :exec
update stations 
  set tracks = $2
  where id = $1;

-- name: GetStation :one
select * from stations where id = $1;

-- name: FindStations :many
select * from stations where id IN (sqlc.slice('ids'));

-- name: DeleteStation :exec
delete from stations where id = $1;

-- name: CountStations :one
select count(*) from stations;

-- name: SearchStations :many
select * from stations where name ILIKE $1 order by tracks desc limit 20;

-- name: SetStationNumberOfTracks :exec
with cte as(
	select station_id, count(*) num_tracks from stop_positions group by station_id
)
update stations set tracks = cte.num_tracks
from cte
where stations.id = cte.station_id;
