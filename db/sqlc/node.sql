-- name: CreateTemporaryNodeTable :exec
CREATE TEMPORARY TABLE tmp_nodes(
  id bigint primary key,
  coordinate point
) ON COMMIT DELETE ROWS;

-- name: InsertTemporaryNodes :copyfrom
insert into tmp_nodes (id, coordinate) values ($1, $2);

-- name: MergeNodesIntoPlatformNodes :many
update platform_nodes u set coordinate = tn.coordinate
from tmp_nodes tn
where u.id = tn.id
returning u.id;
