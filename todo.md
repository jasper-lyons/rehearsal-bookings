# TODO

1. [x] - Write migration tool to manage database migrations
1.1. [ ] - Handle migrations for local sqlite db and production postgresql db

2. [x] - Data Layer Abstraction
2.1. [x] - Lets use the inbuilt database/sql library
2.2. [x] - SQLite Driver
2.2.1. [x] - Query
2.2.2. [x] - Insert
2.2.2. [x] - Update
2.2.2. [x] - Delete
2.3. [ ] - Postgres Driver

3. [x] - Data Access Abstraction
3.1. [x] - Find
3.2. [x] - All
3.3. [x] - Where
3.4. [x] - Create
3.5. [x] - Update
3.6. [x] - Delete
3.7. [ ] - Handle converting between ints in sqlite to bools in go

4.1. [x] - Identify solution for handlers
4.1.1. [x] - We'll have a big routes file and have each handler exist in it's own file. Each handler can then be tested.
4.1.2. [x] - Figure out middleware so that all handlers can have auth checks etc.

5. [ ] - Booking Use Case: As a Guest I can book a rehearsal slot.
5.1. [x] - As a Guest I can select a rehearsal type so that I pay an appropriate ammount for my rehearsal slot.
5.2. [x] - As a Guest I can select a rehearsal room so that I know where I'm rehearsing.
5.3. [x] - As a Guest I can select a date so that the the rehearsal space stays available for me.
5.4. [ ] - As a Guest I can select a time so that the the rehearsal space stays available for me on the selected date.
5.4.1. [ ] - Ensure no one can book overlapping timeslots in the same room.
5.5. [ ] - As a Guest I can enter my contact details (name, email, phone number) so that I can confirm that my booking was successful.
5.6. [ ] - As a Guest I can pay for my booking to ensure that no one else books my slot!
5.6.1. [x] - We need a cart page to allow guests to check they are booking the correct slot.
5.6.2. [ ] - We need a payment confirmation page so that Guests can trust they've paid for their rehearsal slot.

6. [ ] - Launch new Bookings
6.3. [x] - Test SumUp Details - No test details needed with test accounts 
6.1. [x] - Email on successful booking
6.1.1 [x] - Soft Delete held booking on payment fail
6.2. [x] - Booking Success Page
6.4. [x] - Password protected Admin Interface
6.5. [x] - Additional Items on bookings 
6.5.1 [x] - Cymbals
6.5.2 [x] - Notes
6.6. [ ] - Discount Codes
6.7. [x] - Speed up Sumup Widget Loading
6.8. [ ] - Store transaction id from stripe in database as stripe-<id>
6.9. [ ] - Store booking id as transaction reference in stripe
6.10. [ ] - Implement Backup Strategy
6.10.1. [ ] - daily backups, to 7 days
6.10.2. [ ] - weekly backups, to 4 weeks
6.10.3. [ ] - monthly backups, to 12 months
