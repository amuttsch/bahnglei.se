-- name: CreateRoute :one
insert into routes (route_id, stop_position_id, from_station, to_station, via, ref, name, operator, network, service, country_iso_code) 
  values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
  on conflict (route_id, stop_position_id) do update set from_station = $3, to_station = $4, via = $5, ref = $6, name = $7, operator = $8, network = $9, service = $10, updated_at = CURRENT_TIMESTAMP
  returning *;
 
-- name: FindRoutesForStopPosition :many
select * from routes where stop_position_id = $1 order by name;

-- name: DeleteRoutesUpdatedBefore :exec
delete from routes where country_iso_code = $1 and updated_at < $2;
