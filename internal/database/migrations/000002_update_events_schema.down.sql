ALTER TABLE events DROP COLUMN start_time;

ALTER TABLE events DROP COLUMN end_time;

ALTER TABLE events DROP COLUMN location_id;

DROP TABLE IF EXISTS locations;

DROP TABLE IF EXISTS ticket_types;