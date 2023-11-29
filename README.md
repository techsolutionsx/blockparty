# Test Project Setup and Run Guide

## Step 1: Navigate to the Test Project

Open your terminal and navigate to the test project directory. You can choose between two different directories:

```bash
cd test1&2
```

or

```bash
cd test3
```

## Step 2: Run the Project

Execute the following command to run the main.go file:

```bash
go run main.go
```

# Troubleshooting: Fixing Connection Error

If you encounter the following error:

```plaintext
Error inserting metadata into database: pq: no pg_hba.conf entry for host "50.3.70.137", user "jlee", database "jlee", no encryption
```

Follow the steps below to resolve the issue:

## PostgreSQL Host-Based Authentication Configuration

1. Open the PostgreSQL configuration file `pg_hba.conf`. This file is usually located in the PostgreSQL data directory.

2. Add an entry that allows connections from the host "50.3.70.137" for the user "jlee" to the specified database without encryption. Ensure that the entry resembles the following:

   ```plaintext
   host    jlee    jlee    50.3.70.137/32    trust
   ```

   Note: The "trust" method is used here for simplicity. In a production environment, it's recommended to use more secure authentication methods.

3. Save the changes to `pg_hba.conf`.

## Restart PostgreSQL Server

After making changes to the `pg_hba.conf` file, restart the PostgreSQL server to apply the new configuration.

```bash
# Example using systemctl
sudo systemctl restart postgresql
```

## Retry Running the Project

Once the configuration is updated and the PostgreSQL server is restarted, rerun the project using the previously mentioned command:

```bash
go run main.go
```

This should resolve the connection error, and the application should now be able to connect to the PostgreSQL database successfully.

Make sure to replace "50.3.70.137," "jlee," and "jlee" with the actual IP address, username, and database name specified in your application configuration. Additionally, adjust the authentication method based on your security requirements.
