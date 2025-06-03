ALTER TABLE
    bookings DROP COLUMN updated_at;

ALTER TABLE
    bookings
ADD
    COLUMN seat_id BIGINT NOT NULL REFERENCES seats(id),
ADD
    COLUMN confirmed_at TIMESTAMP,
ADD
    COLUMN cancelled_at TIMESTAMP;