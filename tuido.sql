-- Table to store data about projects (can be recursive, ie: can have subprojects)
CREATE TABLE IF NOT EXISTS Project(
    projectId           INTEGER PRIMARY KEY,
    projectParent       INTEGER NULL,
    projectName         TEXT NOT NULL,
    projectDesc         TEXT NULL,
    projectCreationDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    projectType         UNSIGNED INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (projectParent) REFERENCES Project(projectId)
);

-- Table to store data about tasks (can be recursive, ie: can have subtasks)
CREATE TABLE IF NOT EXISTS Task(
    taskId              INTEGER PRIMARY KEY,
    taskParent          INTEGER NULL,
    taskName            TEXT NOT NULL,
    taskDesc            TEXT NULL,
    taskDueDate         DATETIME NULL,
    taskCompleted       UNSIGNED INTEGER DEFAULT 0,
    taskCreationDate    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    taskCompletionDate  DATETIME NULL,
    taskType            UNSIGNED INTEGER NOT NULL DEFAULT 0,
    projectId           INTEGER NOT NULL,
    FOREIGN KEY (taskParent)        REFERENCES Task(taskId),
    FOREIGN KEY (projectId)         REFERENCES Project(projectId)
);
