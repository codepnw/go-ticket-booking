ALTER TABLE
    bookings DROP COLUMN created_at;

ALTER TABLE
    bookings DROP COLUMN updated_at;

ALTER TABLE
    bookings
ADD
    COLUMN seat_id BIGINT NOT NULL REFERENCES seats(id),
ADD
    COLUMN booked_at TIMESTAMP NOT NULL,
ADD
    COLUMN cancelled_at TIMESTAMP;