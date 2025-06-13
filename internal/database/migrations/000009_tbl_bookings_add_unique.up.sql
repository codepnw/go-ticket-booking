ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_user_id_key;

ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_event_id_key;

ALTER TABLE bookings 
ADD CONSTRAINT unique_booking UNIQUE (user_id, event_id, seat_id);