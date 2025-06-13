DROP VIEW IF EXISTS user_view;
CREATE VIEW IF NOT EXISTS user_view AS
SELECT
    LOWER(customer_name) as user_name,
    LOWER(customer_email) as user_email,
    customer_phone as user_phone,
    MAX(created_at) AS last_booking_created_date
FROM bookings
WHERE 
    LOWER(customer_name) NOT IN ('test test') AND
    customer_email NOT IN ('-', 'test@test.com') AND
    customer_phone IS NOT NULL
GROUP BY LOWER(customer_name), LOWER(customer_email), customer_phone
ORDER BY LOWER(customer_name);
