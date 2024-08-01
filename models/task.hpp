#include <ctime>

typedef struct {
    int             id;              // Unique task identifier
    int             projectId;       // ID of the project this task belongs to
    char*           name;            // Name of the task
    char*           description;     // Optional description of the task
    bool            taskCompleted;   // True if the task is completed, false otherwise
    time_t          dueDate;         // Due date for the task
    time_t          creationDate;    // Timestamp of task creation
    time_t          completionDate;  // Completion date for the task
} Task;
