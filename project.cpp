#include "project.hpp"

// Callback function for sqlite3_exec
static int callback(void *data, int argc, char **argv, char **azColName) {
    std::vector<Project>* projects = static_cast<std::vector<Project>*>(data);

    // Convert the retrieved data to appropriate types
    int id = std::stoi(argv[0]);
    std::string name(argv[1]);
    std::string description(argv[2]);
    time_t creationDate = std::stoi(argv[3]); // Assuming creationDate is stored as an integer timestamp
    size_t taskCount = std::stoull(argv[4]);

    Project p(id, name, description, creationDate, taskCount);
    projects->push_back(p);

    return 0;
}

// Create a new project
int Project::createProject(sqlite3* db) {
    std::string sql = "INSERT INTO Project (name, description, creationDate, taskCount) VALUES (?, ?, datetime('now'), ?);";
    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(db, sql.c_str(), -1, &stmt, nullptr);

    if (rc != SQLITE_OK) {
        std::cerr << "Failed to prepare statement: " << sqlite3_errmsg(db) << std::endl;
        return rc;
    }

    sqlite3_bind_text(stmt, 1, name.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, description.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_int(stmt, 3, taskCount);

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) {
        std::cerr << "Execution failed: " << sqlite3_errmsg(db) << std::endl;
        sqlite3_finalize(stmt);
        return rc;
    }

    id = sqlite3_last_insert_rowid(db); // Get the ID of the newly inserted project
    sqlite3_finalize(stmt);
    return SQLITE_OK;
}

// Update an existing project
int Project::updateProject(int id, std::string newName, std::string newDescription, sqlite3* db) {
    std::string sql = "UPDATE Project SET name = ?, description = ? WHERE id = ?";
    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(db, sql.c_str(), -1, &stmt, nullptr);

    if (rc != SQLITE_OK) {
        std::cerr << "Failed to prepare statement: " << sqlite3_errmsg(db) << std::endl;
        return rc;
    }

    sqlite3_bind_text(stmt, 1, newName.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, newDescription.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_int(stmt, 3, id);

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) {
        std::cerr << "Execution failed: " << sqlite3_errmsg(db) << std::endl;
        sqlite3_finalize(stmt);
        return rc;
    }

    sqlite3_finalize(stmt);
    return SQLITE_OK;
}

// Delete an existing project
int Project::deleteProject(int id, sqlite3* db) {
    std::string sql = "DELETE FROM Project WHERE id = ?";
    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(db, sql.c_str(), -1, &stmt, nullptr);

    if (rc != SQLITE_OK) {
        std::cerr << "Failed to prepare statement: " << sqlite3_errmsg(db) << std::endl;
        return rc;
    }

    sqlite3_bind_int(stmt, 1, id);

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) {
        std::cerr << "Execution failed: " << sqlite3_errmsg(db) << std::endl;
        sqlite3_finalize(stmt);
        return rc;
    }

    sqlite3_finalize(stmt);
    return SQLITE_OK;
}

// List all projects
std::vector<Project> Project::listProjects(sqlite3* db) {
    std::vector<Project> projects;
    const char* sql = "SELECT id, name, description, strftime('%s', creationDate), taskCount FROM Project"; // '%s' for timestamp
    char* errMsg = nullptr;

    int rc = sqlite3_exec(db, sql, callback, &projects, &errMsg);
    if (rc != SQLITE_OK) {
        std::cerr << "SQL error: " << errMsg << std::endl;
        sqlite3_free(errMsg);
    }

    return projects;
}
