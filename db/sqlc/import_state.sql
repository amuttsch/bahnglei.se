-- name: CreateImportState :one
insert into import_state (country_iso_code, state, number_platforms, number_stations) 
  values ($1, 'starting', 0, 0) 
  returning *;

-- name: UpdateImportState :exec
update import_state  
  set state = $2, number_stations = $3, number_platforms = $4
  where id = $1;

