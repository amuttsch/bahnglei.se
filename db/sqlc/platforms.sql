-- name: CreatePlatform :one
insert into platforms (id, station_id, positions) 
  values ($1, $2, $3) 
  returning *;

-- name: DeletePlatformsForStation :exec
delete from platforms where station_id = $1;

-- name: GetPlatformsForStation :many
select * from platforms where station_id = $1;
