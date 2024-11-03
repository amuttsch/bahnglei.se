create extension if not exists pg_trgm;

CREATE TABLE IF NOT EXISTS countries (
  iso_code text NOT NULL PRIMARY KEY,
  name text not null,
  osm_url text not null,
  created_at timestamp
  with
    time zone not null default current_timestamp,
    updated_at timestamp
  with
    time zone
);

CREATE TABLE IF NOT EXISTS stations (
  id bigint PRIMARY KEY,
  country_iso_code text not null,
  name text not null,
  coordinate point not null,
  operator text not null,
  wikidata text not null,
  wikipedia text not null,
  tracks bigint not null,
  created_at timestamp
  with
    time zone not null default current_timestamp,
    updated_at timestamp
  with
    time zone not null default current_timestamp,
    CONSTRAINT fk_stations_country FOREIGN KEY (country_iso_code) REFERENCES countries (iso_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS platforms (
  id bigint primary key,
  country_iso_code text not null,
  positions text not null,
  station_id bigint,
  coordinate point,
  created_at timestamp
  with
    time zone not null default current_timestamp,
    updated_at timestamp
  with
    time zone not null default current_timestamp,
    CONSTRAINT fk_stations_platforms FOREIGN KEY (station_id) REFERENCES stations (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT fk_country FOREIGN KEY (country_iso_code) REFERENCES countries (iso_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE platform_nodes (
  id bigint primary key,
  platform_id bigint not null,
  country_iso_code text not null,
  coordinate point,
  CONSTRAINT fk_platforms FOREIGN KEY (platform_id) REFERENCES platforms (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT fk_country FOREIGN KEY (country_iso_code) REFERENCES countries (iso_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS stop_positions (
  id bigint PRIMARY KEY,
  station_id bigint,
  country_iso_code text not null,
  platform text not null default '',
  coordinate point not null,
  neighbors text not null default '',
  created_at timestamp
  with
    time zone default current_timestamp,
    updated_at timestamp
  with
    time zone not null default current_timestamp,
    CONSTRAINT fk_country FOREIGN KEY (country_iso_code) REFERENCES countries (iso_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS osm_tiles (
  id BIGSERIAL PRIMARY KEY,
  x bigint not null,
  y bigint not null,
  z bigint not null,
  data bytea not null,
  created_at timestamp
  with
    time zone default current_timestamp,
    updated_at timestamp
  with
    time zone not null default current_timestamp
);

CREATE TABLE IF NOT EXISTS import_state (
  id SERIAL PRIMARY KEY,
  country_iso_code text NOT NULL,
  number_stations int not null,
  number_platforms int not null,
  state text not null,
  created_at timestamp
  with
    time zone default current_timestamp,
    updated_at timestamp
  with
    time zone not null default current_timestamp,
    CONSTRAINT fk_import_state_country FOREIGN KEY (country_iso_code) REFERENCES countries (iso_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TEMPORARY TABLE tmp_nodes (id bigint primary key, coordinate point) ON COMMIT DELETE ROWS;
