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
  lat float not null,
  lng float not null,
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
  id serial primary key,
  positions text not null,
  station_id bigint not null,
  created_at timestamp
  with
    time zone not null default current_timestamp,
    updated_at timestamp
  with
    time zone not null default current_timestamp,
    CONSTRAINT fk_stations_platforms FOREIGN KEY (station_id) REFERENCES stations (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS stop_positions (
  id SERIAL PRIMARY KEY,
  station_id bigint NOT NULL,
  platform text not null,
  lat float not null,
  lng float not null,
  neighbors text not null,
  created_at timestamp
  with
    time zone default current_timestamp,
    updated_at timestamp
  with
    time zone not null default current_timestamp,
    CONSTRAINT fk_stations_stop_position FOREIGN KEY (station_id) REFERENCES public.stations (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS osm_tiles (
  id SERIAL PRIMARY KEY,
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
)
