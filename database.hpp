#ifndef DATABASE_HPP
#define DATABASE_HPP

#include <sqlite3.h>
#include <string>

class Database {
private:
    static std::string databasePath;  // Path to the database file
    static sqlite3* db;               // Pointer to the SQLite database

    // Create the database file if it doesn't exist
    static bool createDatabase();

    // Execute SQL commands to create necessary tables
    static bool initializeSchema();

    // Check if the database exists at the specified path
    static bool isDatabaseCreated();

public:
    // Initialize the database and get the connection
    static bool initialize();

    // Get the database connection
    static sqlite3* getConnection();

    // Clean up the database connection
    static void closeConnection();
};

#endif // DATABASE_HPP
