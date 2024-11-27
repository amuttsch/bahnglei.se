create index IF NOT EXISTS idx_coordinates on osm_tiles (
  x ASC NULLS LAST,
  y ASC NULLS LAST,
  z ASC NULLS LAST
);

CREATE INDEX IF NOT EXISTS stations_name_gin_trgm_idx ON stations USING gin (name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS routes_route_id ON routes (route_id ASC);

CREATE INDEX IF NOT EXISTS routes_stop_position_id ON routes (stop_position_id ASC);
