#ifndef TUIDO_H
#define TUIDO_H

#include <stdint.h>
#include <stdbool.h>
#include <time.h>
#include <stdlib.h>

// Define a type for representing the possibility of NULL integers in C
typedef struct {
    int value;
    int is_null;
} NullableInt;

// Define a type for representing the possibility of NULL datetime in C
typedef struct {
    time_t value;    // Time value represented as time_t
    int is_null;     // Flag to check if the value is NULL
} NullableTime;

// Enum for representing the types of projects and tasks
typedef enum {
    TOP_LEVEL,   // Top-level item
    SUB_LEVEL    // Sub-level item
} ItemType;

typedef struct Project Project;
typedef struct Task Task;

struct Project {
    int projectId;                   // Unique project identifier
    NullableInt projectParent;       // Nullable parent project ID (NULL if top-level)
    char *projectName;               // Name of the project
    char *projectDesc;               // Optional description of the project
    time_t projectCreationDate;      // Timestamp of project creation
    ItemType projectType;            // Enum type of the project
    Project **subProjects;           // Pointer to an array of subproject pointers
    size_t subProjectCount;          // Number of subprojects
    Task **tasks;                    // Pointer to an array of task pointers related to this project
    size_t taskCount;                // Number of tasks related to this project
};

struct Task {
    int taskId;                      // Unique task identifier
    NullableInt taskParent;          // Nullable parent task ID (NULL if top-level)
    char *taskName;                  // Name of the task
    char *taskDesc;                  // Optional description of the task
    NullableTime taskDueDate;        // Optional due date for the task
    bool taskCompleted;              // True if the task is completed, false otherwise
    time_t taskCreationDate;         // Timestamp of task creation
    NullableTime taskCompletionDate; // Nullable completion date of the task
    ItemType taskType;               // Enum type of the task
    int projectId;                   // ID of the project this task belongs to
    Task **subTasks;                 // Pointer to an array of subtask pointers
    size_t subTaskCount;             // Number of subtasks
};

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

#endif // TUIDO_H
