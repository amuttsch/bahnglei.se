#!/bin/zsh

TABLES=(countries osm_tiles stations platforms stop_positions platform_nodes platform_ways)

cd tmp

for table in "$TABLES"; do
  pg_dump --host localhost --port 5432 --username postgres --password postgres --file "$table.sql" --table "$table" "bahngleise"
done
