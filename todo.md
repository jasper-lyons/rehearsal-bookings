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
4.1. [x] - Identify solution for handlers
4.1.1. [x] - We'll have a big routes file and have each handler exist in it's own file. Each handler can then be tested.
4.1.2. [ ] - Figure out middleware so that all handlers can have auth checks etc.

