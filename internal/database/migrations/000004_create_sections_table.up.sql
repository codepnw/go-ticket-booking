CREATE TABLE sections (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    seat_count INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    section_id INT NOT NULL REFERENCES sections(id),
    row_label VARCHAR(10) NOT NULL,
    seat_number INT DEFAULT 1,
    is_available BOOLEAN DEFAULT TRUE
);