#ifndef CLIDO_H
#define CLIDO_H

#include <stdint.h>
#include <stdbool.h>
#include <time.h>
#include <stdlib.h>

// DATA STRUCTURES

// Define a type for representing the possibility of NULL integers in C
typedef struct {
    int value;      // Integer value
    int is_null;    // Flag to check if the value is NULL
} NullableInt;

// Define a type for representing the possibility of NULL datetime in C
typedef struct {
    time_t  value;    // Time value represented as time_t
    int     is_null;  // Flag to check if the value is NULL
} NullableTime;

// Enum for representing the types of projects and tasks
typedef enum {
    TOP_LEVEL,   // Top-level item
    SUB_LEVEL    // Sub-level item
} ItemType;

typedef struct Project Project;
typedef struct Task Task;

struct Project {
    int         id;                 // Unique project identifier
    NullableInt parent;             // Nullable parent project ID (NULL if top-level)
    char*       name;               // Name of the project
    char*       description;        // Optional description of the project
    time_t      creationDate;       // Timestamp of project creation
    ItemType    type;               // Enum type of the project
    Project**   subProjects;        // Pointer to an array of subproject pointers
    size_t      subProjectCount;    // Number of subprojects
    Task**      tasks;              // Pointer to an array of task pointers related to this project
    size_t      taskCount;          // Number of tasks related to this project
};

struct Task {
    int             id;              // Unique task identifier
    NullableInt     parent;          // Nullable parent task ID (NULL if top-level)
    char*           name;            // Name of the task
    char*           description;     // Optional description of the task
    NullableTime    dueDate;         // Optional due date for the task
    bool            taskCompleted;   // True if the task is completed, false otherwise
    time_t          creationDate;    // Timestamp of task creation
    NullableTime    completionDate;  // Nullable completion date of the task
    ItemType        type;            // Enum type of the task
    int             projectId;       // ID of the project this task belongs to
    Task**          subTasks;        // Pointer to an array of subtask pointers
    size_t          subTaskCount;    // Number of subtasks
};

// FUNCTION PROTOTYPES

// Project Management
Project *create_project(int projectId, const char *name, const char *desc, ItemType type);
void add_subproject(Project *parent, Project *subproject);
void add_task_to_project(Project *project, Task *task);
void free_project(Project *project);
Project *find_project_by_id(Project *root, int projectId);
void traverse_projects(Project *root, void (*visit)(Project *));
bool update_project(Project *project, const char *newName, const char *newDesc);
bool remove_project(Project **root, int projectId);
size_t count_subprojects(Project *project);
bool check_project_completion(Project *project);
void print_project_tree(Project *project, int level);

// Task Management
Task *create_task(int taskId, const char *name, const char *desc, ItemType type, int projectId);
void add_subtask(Task *parent, Task *subtask);
void free_task(Task *task);
Task *find_task_by_id(Task *root, int taskId);
void traverse_tasks(Task *root, void (*visit)(Task *));
bool update_task(Task *task, const char *newName, const char *newDesc, bool completed);
bool remove_task(Task **root, int taskId);
size_t count_subtasks(Task *task);
void print_task_tree(Task *task, int level);

#endif // CLIDO_H
