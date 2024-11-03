create index idx_coordinates on osm_tiles (
  x ASC NULLS LAST,
  y ASC NULLS LAST,
  z ASC NULLS LAST
);

CREATE INDEX stations_name_gin_trgm_idx ON stations USING gin (name gin_trgm_ops);
