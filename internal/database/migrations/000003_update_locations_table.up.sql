ALTER TABLE locations 
ADD COLUMN description TEXT,
ADD COLUMN capacity INT NOT NULL,
ADD COLUMN owner_id INT NOT NULL;

ALTER TABLE events
DROP COLUMN event_date,
DROP COLUMN total_tickets;