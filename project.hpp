#ifndef PROJECT_HPP
#define PROJECT_HPP

#include <iostream>
#include <sqlite3.h>
#include <vector>
#include <ctime>

class Project {
    private:
        int id;
        std::string name;
        std::string description;
        time_t creationDate;
        size_t taskCount;

    public:
        Project(int id, std::string name, std::string description, time_t creationDate, size_t taskCount)
            : id(id), name(name), description(description), creationDate(creationDate), taskCount(taskCount) {}

        int createProject(sqlite3* db);
        int updateProject(int id, std::string newName, std::string newDescription, sqlite3* db);
        int deleteProject(int id, sqlite3* db);
        static std::vector<Project> listProjects(sqlite3* db);
};

#endif // PROJECT_HPP
