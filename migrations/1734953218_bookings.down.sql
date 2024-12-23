-- Drop the trigger first
DROP TRIGGER IF EXISTS update_bookings_updated_at;

-- Drop the index
DROP INDEX IF EXISTS idx_bookings_room_times;

-- Drop the bookings table
DROP TABLE IF EXISTS bookings;
