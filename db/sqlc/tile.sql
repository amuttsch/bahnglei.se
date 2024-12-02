-- name: GetTile :one
select * from osm_tiles where x = $1 and y = $2 and z = $3;

-- name: CreateTile :one
insert into osm_tiles (x, y, z, data) values ($1, $2, $3, $4) returning *;
