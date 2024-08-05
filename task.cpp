#include <sqlite3.h>
#include <iostream>
#include <vector>
#include <ctime>
#include "task.hpp"

// Callback function for sqlite3_exec to populate a vector of Task objects
static int taskCallback(void *data, int argc, char **argv, char **azColName) {
    std::vector<Task>* tasks = static_cast<std::vector<Task>*>(data);

    // Extract fields from the row
    int id = std::stoi(argv[0]);
    int parentProjectId = std::stoi(argv[1]);
    std::string name(argv[2]);
    std::string description(argv[3]);
    bool taskCompleted = (std::stoi(argv[4]) != 0);  // Convert from int to bool
    time_t dueDate = std::stoll(argv[5]);  // Assuming time_t is stored as a large integer
    time_t creationDate = std::stoll(argv[6]);
    time_t completionDate = std::stoll(argv[7]);

    // Create Task object and add it to the vector
    Project* parentProject = nullptr;  // You may need to fetch the Project details if required
    Task t(id, parentProject, name, description, taskCompleted, dueDate, creationDate, completionDate);
    tasks->push_back(t);

    return 0;
}

// Create a new task in the database
int Task::createTask(sqlite3 *db) {
    std::string sql = "INSERT INTO Task (name, description, dueDate, completed, creationDate, completionDate, projectId) VALUES (?, ?, ?, ?, datetime('now'), ?, ?);";
    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(db, sql.c_str(), -1, &stmt, nullptr);

    if (rc != SQLITE_OK) {
        std::cerr << "Failed to prepare statement: " << sqlite3_errmsg(db) << std::endl;
        return rc;
    }

    sqlite3_bind_text(stmt, 1, name.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, description.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_int(stmt, 3, static_cast<int>(dueDate));  // Assuming time_t can be cast to int
    sqlite3_bind_int(stmt, 4, taskCompleted ? 1 : 0);
    sqlite3_bind_int(stmt, 5, static_cast<int>(completionDate));
    sqlite3_bind_int(stmt, 6, parentProject ? parentProject->getId() : -1);  // Bind project ID

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) {
        std::cerr << "Execution failed: " << sqlite3_errmsg(db) << std::endl;
        sqlite3_finalize(stmt);
        return rc;
    }

    id = sqlite3_last_insert_rowid(db); // Retrieve the new task's ID
    sqlite3_finalize(stmt);
    return SQLITE_OK;
}

// Update an existing task
int Task::updateTask(int id, std::string newName, std::string newDescription, time_t newDueDate, sqlite3 *db) {
    std::string sql = "UPDATE Task SET name = ?, description = ?, dueDate = ? WHERE id = ?";
    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(db, sql.c_str(), -1, &stmt, nullptr);

    if (rc != SQLITE_OK) {
        std::cerr << "Failed to prepare statement: " << sqlite3_errmsg(db) << std::endl;
        return rc;
    }

    sqlite3_bind_text(stmt, 1, newName.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, newDescription.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_int(stmt, 3, static_cast<int>(newDueDate));  // Assuming time_t can be cast to int
    sqlite3_bind_int(stmt, 4, id);

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) {
        std::cerr << "Execution failed: " << sqlite3_errmsg(db) << std::endl;
        sqlite3_finalize(stmt);
        return rc;
    }

    sqlite3_finalize(stmt);
    return SQLITE_OK;
}

// Mark the task as completed
int Task::setAsCompleted(int id, sqlite3 *db) {
    std::string sql = "UPDATE Task SET completed = 1, completionDate = datetime('now') WHERE id = ?";
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

// Delete a task from the database
int Task::deleteTask(int id, sqlite3 *db) {
    std::string sql = "DELETE FROM Task WHERE id = ?";
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

// List all tasks
std::vector<Task> Task::listTasks(sqlite3 *db) {
    std::vector<Task> tasks;
    const char* sql = "SELECT id, projectId, name, description, completed, dueDate, creationDate, completionDate FROM Task";
    char* errMsg = nullptr;

    int rc = sqlite3_exec(db, sql, taskCallback, &tasks, &errMsg);
    if (rc != SQLITE_OK) {
        std::cerr << "SQL error: " << errMsg << std::endl;
        sqlite3_free(errMsg);
    }

    return tasks;
}

// List tasks for a specific project
std::vector<Task> listTasksInProject(int projectId, sqlite3 *db) {
    std::vector<Task> tasks;
    std::string sql = "SELECT id, name, description, completed, dueDate, creationDate, completionDate FROM Task WHERE projectId = ?";
    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(db, sql.c_str(), -1, &stmt, nullptr);

    if (rc != SQLITE_OK) {
        std::cerr << "Failed to prepare statement: " << sqlite3_errmsg(db) << std::endl;
        return tasks;
    }

    sqlite3_bind_int(stmt, 1, projectId);

    while ((rc = sqlite3_step(stmt)) == SQLITE_ROW) {
        int id = sqlite3_column_int(stmt, 0);
        std::string name(reinterpret_cast<const char*>(sqlite3_column_text(stmt, 1)));
        std::string description(reinterpret_cast<const char*>(sqlite3_column_text(stmt, 2)));
        bool taskCompleted = (sqlite3_column_int(stmt, 3) != 0);
        time_t dueDate = static_cast<time_t>(sqlite3_column_int(stmt, 4));
        time_t creationDate = static_cast<time_t>(sqlite3_column_int(stmt, 5));
        time_t completionDate = static_cast<time_t>(sqlite3_column_int(stmt, 6));

        Task t(id, nullptr, name, description, taskCompleted, dueDate, creationDate, completionDate);
        tasks.push_back(t);
    }

    if (rc != SQLITE_DONE) {
        std::cerr << "Execution failed: " << sqlite3_errmsg(db) << std::endl;
    }

    sqlite3_finalize(stmt);
    return tasks;
}
