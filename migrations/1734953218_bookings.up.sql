-- Create bookings table
CREATE TABLE bookings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Customer information
    customer_name TEXT NOT NULL,
    customer_email TEXT NOT NULL,
    customer_phone TEXT,  -- Optional phone number
    
    -- Booking details
    room_name TEXT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    
    -- Optional notes field for any special requirements
    notes TEXT,
    
    -- Ensure end_time is after start_time
    CHECK (end_time > start_time)
);

-- Create index for common queries
CREATE INDEX idx_bookings_room_times ON bookings(room_name, start_time, end_time);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_bookings_updated_at 
    AFTER UPDATE ON bookings
    BEGIN
        UPDATE bookings 
        SET updated_at = DATETIME('now') 
        WHERE id = NEW.id;
    END;
