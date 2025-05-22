ALTER TABLE
    bookings
ADD
    COLUMN created_at TIMESTAMP,
    COLUMN updated_at TIMESTAMP;

ALTER TABLE
    bookings DROP COLUMN seat_id,
    DROP COLUMN booked_at,
    DROP COLUMN cancelled_at;