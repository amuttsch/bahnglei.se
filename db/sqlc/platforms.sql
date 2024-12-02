-- name: CreatePlatform :one
insert into platforms (id, station_id, positions, coordinate, country_iso_code) 
  values ($1, $2, $3, $4, $5) 
  on conflict (id) do update set positions = $3, coordinate = $4, updated_at = CURRENT_TIMESTAMP
  returning *;

-- name: DeletePlatformsForStation :exec
delete from platforms where station_id = $1;

-- name: GetPlatformsForStation :many
select * from platforms where station_id = $1;

-- name: FindPlatforms :many
select * from platforms where id IN (sqlc.slice('ids'));

-- name: UpdatePlatformSetStationId :exec
update platforms set station_id = $2, updated_at = CURRENT_TIMESTAMP where id = $1;

-- name: SetPlatformToNearestStation :exec
update platforms u set station_id = s.id, updated_at = CURRENT_TIMESTAMP
from platforms p
left join lateral (
  select * from stations s where s.country_iso_code = $1 order by s.coordinate <-> p.coordinate limit 1
) s on TRUE
WHERE u.id = p.id and p.id IN (SELECT id
             FROM platforms inner_platform
             WHERE inner_platform.country_iso_code = $1 and inner_platform.station_id is NULL);

-- name: DeletePlatformsUpdatedBefore :exec
delete from platforms where country_iso_code = $1 and updated_at < $2;
