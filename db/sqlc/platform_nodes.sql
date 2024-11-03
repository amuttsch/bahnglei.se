-- name: CreatePlatformNode :exec
insert into platform_nodes (id, platform_id, country_iso_code) values ($1, $2, $3) on conflict (id) do nothing;

-- name: GetPlatformNode :one
select * from platform_nodes where id = $1;

-- name: UpdatePlatformNode :exec
update platform_nodes set coordinate = $2 where id = $1;
