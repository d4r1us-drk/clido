#include <iostream>
#include <cstdlib>
#include <sys/stat.h>
#include "database.hpp"

// Initialize static members
std::string Database::databasePath = "";
sqlite3* Database::db = nullptr;

// Helper function to get the default database path
static std::string getDefaultDatabasePath() {
    const char* xdgDataHome = std::getenv("XDG_DATA_HOME");
    if (xdgDataHome) {
        return std::string(xdgDataHome) + "/clido/data.db";
    } else {
        return std::string(getenv("HOME")) + "/.local/share/clido/data.db";
    }
}

// SQL commands to create the necessary tables
static const char* createTablesSQL = R"(
CREATE TABLE IF NOT EXISTS Project(
    id           INTEGER PRIMARY KEY,
    name         TEXT NOT NULL,
    description  TEXT NULL,
    creationDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    taskCount    INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Task(
    id              INTEGER PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT NULL,
    dueDate         DATETIME NULL,
    completed       UNSIGNED INTEGER DEFAULT 0,
    creationDate    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completionDate  DATETIME NULL,
    projectId       INTEGER NOT NULL,
    FOREIGN KEY (projectId) REFERENCES Project(id)
);
)";

// Create the database file if it doesn't exist
bool Database::createDatabase() {
    if (isDatabaseCreated()) {
        return true;
    }

    // Ensure the directory exists
    size_t pos = databasePath.find_last_of("/\\");
    std::string directory = databasePath.substr(0, pos);
    mkdir(directory.c_str(), 0755);

    // Create or open the database file
    int rc = sqlite3_open(databasePath.c_str(), &db);
    if (rc != SQLITE_OK) {
        std::cerr << "Failed to create/open database: " << sqlite3_errmsg(db) << std::endl;
        return false;
    }

    // Initialize the database schema
    if (!initializeSchema()) {
        sqlite3_close(db);
        db = nullptr;
        return false;
    }

    return true;
}

// Check if the database exists at the specified path
bool Database::isDatabaseCreated() {
    struct stat buffer;
    return (stat(databasePath.c_str(), &buffer) == 0);
}

// Execute SQL commands to create necessary tables
bool Database::initializeSchema() {
    char* errMsg = nullptr;
    int rc = sqlite3_exec(db, createTablesSQL, nullptr, nullptr, &errMsg);
    if (rc != SQLITE_OK) {
        std::cerr << "Failed to create tables: " << errMsg << std::endl;
        sqlite3_free(errMsg);
        return false;
    }
    return true;
}

// Initialize the database and get the connection
bool Database::initialize() {
    databasePath = getDefaultDatabasePath();
    return createDatabase();
}

// Get the database connection
sqlite3* Database::getConnection() {
    return db;
}

// Clean up the database connection
void Database::closeConnection() {
    if (db) {
        sqlite3_close(db);
        db = nullptr;
    }
}
