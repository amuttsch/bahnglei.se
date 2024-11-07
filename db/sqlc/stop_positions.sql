-- name: CreateStopPosition :one
insert into stop_positions (id, station_id, platform, coordinate, neighbors, country_iso_code) 
  values ($1, $2, $3, $4, $5, $6) 
  on conflict(id) do update set coordinate = $4
  returning *;

-- name: DeleteStopPositionsForStation :exec
delete from stop_positions where station_id = $1;

-- name: GetStopPositionsForStation :many
select * from stop_positions where station_id = $1;

-- name: GetStopPositionsForStationAndPlatform :one
select * from stop_positions where station_id = $1 and platform = $2;

-- name: FindStopPositions :many
select * from stop_positions where id IN (sqlc.slice('ids'));

-- name: SetStopPositionStationIdToNearestStation :exec
update stop_positions u set station_id = s.id
from stop_positions sp
left join lateral (
	select * from stations s order by s.coordinate <-> sp.coordinate limit 1
) s on TRUE
WHERE u.id = sp.id and sp.id IN (SELECT id
             FROM stop_positions sp
             WHERE station_id is null
             and sp.country_iso_code = $1);

-- name: SetStopPositionNeighbors :exec
with cte as (
	select sp.*, string_agg(p.positions, ';') n from stop_positions sp
	inner join platforms p on sp.station_id = p.station_id and sp.platform = ANY(STRING_TO_ARRAY(p.positions, ';'))
  where sp.country_iso_code = $1
	group by sp.id
)
update stop_positions u set neighbors = cte.n
from cte where u.id = cte.id;
