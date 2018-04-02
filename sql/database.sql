/*
 * Copyright jean-francois PHILIPPE 2014-2018
 */


-- Table des donnees brutes
CREATE TABLE IF NOT EXISTS raw_datas (
    ts  timestamp NOT NULL,
    node_id smallint NOT NULL, -- int2
    sensor_id smallint NOT NULL,
    value integer NOT NULL );

CREATE INDEX idx_raw_datas ON raw_datas ( node_id, sensor_id);


