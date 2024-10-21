-- name: CreateStopPosition :one
insert into stop_positions (id, station_id, platform, lat, lng, neighbors) 
  values ($1, $2, $3, $4, $5, $6) 
  returning *;

-- name: DeleteStopPositionsForStation :exec
delete from stop_positions where station_id = $1;

-- name: GetStopPositionsForStation :many
select * from stop_positions where station_id = $1;

-- name: GetStopPositionsForStationAndPlatform :one
select * from stop_positions where station_id = $1 and platform = $2;

