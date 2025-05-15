-- create locations table
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- update events table
ALTER TABLE events
ADD COLUMN start_time TIMESTAMP NOT NULL,
ADD COLUMN end_time TIMESTAMP NOT NULL,
ADD COLUMN location_id INT REFERENCES locations(id);

-- create ticket_types table
CREATE TABLE ticket_types (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);