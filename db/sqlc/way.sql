-- name: CreateTemporaryWaysTable :exec
CREATE TEMPORARY TABLE IF NOT EXISTS tmp_ways(
  id bigint not null,
  node bigint not null
) ON COMMIT DELETE ROWS;

-- name: InsertTemporaryWays :copyfrom
insert into tmp_ways (id, node) values ($1, $2);

-- name: InsertPlatformNodesFromPlatformWays :many
INSERT INTO platform_nodes (id, platform_id, country_iso_code)
select node, w.platform_id, w.country_iso_code from tmp_ways t
inner join platform_ways w on t.id = w.id
ON CONFLICT (id) DO NOTHING
returning *;
