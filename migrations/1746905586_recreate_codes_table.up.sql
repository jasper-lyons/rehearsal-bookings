DROP TRIGGER IF EXISTS update_codes_updated_at;
DROP TABLE IF EXISTS codes;
-- Recreate the codes table with updated code_value type
-- to TEXT to accommodate longer codes
CREATE TABLE IF NOT EXISTS codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code_name VARCHAR(255) NOT NULL,
    code_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_codes_updated_at 
    AFTER UPDATE ON codes
    BEGIN
        UPDATE codes 
        SET updated_at = DATETIME('now') 
        WHERE id = NEW.id;
    END;

INSERT INTO codes (code_name, code_value) 
VALUES 
    ('Room 1', 1234),
    ('Room 2', 1234),
    ('Monday Front Door', 1234),
    ('Tuesday Front Door', 1234),
    ('Wednesday Front Door', 1234),
    ('Thursday Front Door', 1234),
    ('Friday Front Door', 1234),
    ('Saturday Front Door', 1234),
    ('Sunday Front Door', 1234),
    ('Room 2 Store', 1234),
    ('Rec Room Store', 1234),
    ('Rec Room Keybo', 1234);
