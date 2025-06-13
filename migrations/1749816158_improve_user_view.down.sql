DROP VIEW IF EXISTS user_view;
CREATE VIEW IF NOT EXISTS user_view AS
SELECT
    customer_name as user_name,
    customer_email as user_email,
    customer_phone as user_phone,
    MAX(created_at) AS last_booking_created_date
FROM bookings
GROUP BY customer_name, customer_email, customer_phone
ORDER BY customer_name;
