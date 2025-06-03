ALTER TABLE
    bookings
ADD
    COLUMN updated_at TIMESTAMP;

ALTER TABLE
    bookings DROP COLUMN seat_id,
    DROP COLUMN confirmed_at,
    DROP COLUMN cancelled_at;