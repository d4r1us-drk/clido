-- Table to store data about projects (can be recursive, ie: can have subprojects)
CREATE TABLE IF NOT EXISTS Project(
    id           INTEGER PRIMARY KEY,
    parent       INTEGER NULL,
    name         TEXT NOT NULL,
    description  TEXT NULL,
    creationDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type         UNSIGNED INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (parent) REFERENCES Project(id)
);

-- Table to store data about tasks (can be recursive, ie: can have subtasks)
CREATE TABLE IF NOT EXISTS Task(
    id              INTEGER PRIMARY KEY,
    parent          INTEGER NULL,
    name            TEXT NOT NULL,
    description     TEXT NULL,
    dueDate         DATETIME NULL,
    completed       UNSIGNED INTEGER DEFAULT 0,
    creationDate    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completionDate  DATETIME NULL,
    type            UNSIGNED INTEGER NOT NULL DEFAULT 0,
    projectId       INTEGER NOT NULL,
    FOREIGN KEY (parent)        REFERENCES Task(id),
    FOREIGN KEY (projectId)     REFERENCES Project(id)
);
