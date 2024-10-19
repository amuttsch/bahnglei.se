create extension pg_trgm;

CREATE INDEX stations_name_gin_trgm_idx ON stations USING gin (name gin_trgm_ops);
