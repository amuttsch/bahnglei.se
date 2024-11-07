-- name: CreatePlatformWay :exec
insert into platform_ways (id, platform_id, country_iso_code) values ($1, $2, $3) on conflict (id) do nothing;

-- name: GetPlatformWay :one
select * from platform_ways where id = $1;

