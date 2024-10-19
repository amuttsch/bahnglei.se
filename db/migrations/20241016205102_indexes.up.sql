create index idx_coordinates on osm_tiles (
  x ASC NULLS LAST,
  y ASC NULLS LAST,
  z ASC NULLS LAST
);
