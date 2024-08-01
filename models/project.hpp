#include <ctime>

typedef struct {
    int             id;                 // Unique project identifier
    char*           name;               // Name of the project
    char*           description;        // Optional description of the project
    time_t          creationDate;       // Timestamp of project creation
    size_t          taskCount;          // Number of tasks related to this project
} Project;
