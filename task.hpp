#ifndef TASK_HPP
#define TASK_HPP

#include <iostream>
#include <sqlite3.h>
#include <vector>
#include <string>
#include <ctime>
#include "project.hpp"

class Task {
    private:
        int id;
        Project *parentProject;
        std::string name;
        std::string description;
        bool taskCompleted;
        time_t dueDate;
        time_t creationDate;
        time_t completionDate;

    public:
        // Constructor
        Task(int id, Project *parentProject, std::string name, std::string description,
             bool taskCompleted, time_t dueDate, time_t creationDate, time_t completionDate)
            : id(id), parentProject(parentProject), name(name), description(description),
              taskCompleted(taskCompleted), dueDate(dueDate), creationDate(creationDate),
              completionDate(completionDate) {}

        // Method declarations
        int createTask(sqlite3 *db);
        int updateTask(int id, std::string newName, std::string newDescription, time_t newDueDate, sqlite3 *db);
        int setAsCompleted(int id, sqlite3 *db);
        int deleteTask(int id, sqlite3 *db);
        static std::vector<Task> listTasks(sqlite3 *db);
        static std::vector<Task> listTasksInProject(int projectId);
};

#endif // TASK_HPP
