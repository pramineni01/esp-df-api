# Enterprise SaaS Platform Demand Forecast APIs

## Development & Testing

Running the server needs a database access currently provided by environment variables.
To setup database start by setting up the following environment variables.

    DB_USER
    DB_NAME
    DB_PASSWORD
    DB_HOST

------
*NOTE*

We currently support only password base authentication for database, which can easily be changed in the 
future.

----------

### test_utils
One of the most cumbersome part of writing integration tests is managing the sanity of the data involved.
To address that issue `test_utils` provides a couple of public functions responsible for creating and managing 
test databases.

### Writing Tests
To write tests that depends on data we start by adding that data to `seed.sql` so that the we can load and unload it whenever we need.
After that we ask `test_utils` for a database connection(`test_db`) by calling `GetTestDBConn` which is responsible for taking the schema dump of the provided database, create two database named 
`seed` and `test` from the schema and load the data from `seed.sql` to `seed` database  and return the database connection to the test database.
After that we can just call `LoadTables` and `UnLoadTables` to load and unload tables as and when needed.

### Running the tests
To run tests we need to set an environment variable(apart from the environment variables already set for database connection) named `SERVER_FILE_PATH`, responsible for providing the file system
path of the server.go file so that `test_utils** can start the server if needed.

-----
*NOTE*

Please make sure the database user provided for testing must be super user(because it needs to create and drop database etc.).

------

