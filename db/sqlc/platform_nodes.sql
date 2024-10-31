-- name: CreatePlatformNode :exec
insert into platform_nodes (id, platform_id) values ($1, $2) on conflict (id) do nothing;

-- name: GetPlatformNode :one
select * from platform_nodes where id = $1;

-- name: UpdatePlatformNode :exec
update platform_nodes set coordinate = $2 where id = $1;
