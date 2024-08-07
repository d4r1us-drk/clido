# clido

Clido is an awesome CLI to-do list management application that helps you keep track of your projects and tasks efficiently.

## Table of Contents
1. [About The Project](#about-the-project)
    - [Built With](#built-with)
2. [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Building](#building)
3. [Usage](#usage)
4. [Roadmap](#roadmap)
5. [License](#license)
6. [Acknowledgments](#acknowledgments)

## About The Project

Clido is a simple yet powerful CLI tool designed to help you manage your projects and tasks effectively from the terminal. Whether you are a developer, a project manager, or just someone who loves to keep things organized, Clido is the perfect tool for you.

### Built With

* [Go](https://golang.org/)
* [SQLite](https://www.sqlite.org/index.html)
* [Color](https://github.com/fatih/color) - For colored terminal output
* [Tablewriter](https://github.com/olekukonko/tablewriter) - For table formatting in terminal

## Getting Started

To get a local copy up and running follow these simple steps.

### Installation

1. Grab the official binary from the [releases page](https://github.com/d4r1us-drk/clido/releases) for your operating system and computer architecture. Currently supported operating systems are Windows, Mac and Linux, each both on x86 and ARM.
2. Move the binary anywhere in your PATH.
3. Enjoy.

### Building

Make sure to have Go and SQLite installed.

1. Clone the project
    ```sh
    git clone https://github.com/d4r1us-drk/clido.git
    cd clido
    ```

2. Compile and run
    ```sh
    go build
    ./clido help
    ```

## Usage

Clido allows you to manage projects and tasks with various commands. Below are some usage examples.

### Commands

- Create a new project
  ```sh
  clido new project -n "New Project" -d "Project Description"
  ```

- Create a new task
  ```sh
  clido new task -n "New Task" -d "Task Description" -D "2024-08-15 23:00" -p "Existing Project"
  ```

- Edit an existing project
  ```sh
  clido edit project 1 -n "Updated Project Name" -d "Updated Description"
  ```

- List all projects
  ```sh
  clido list projects
  ```

- List tasks by project number
  ```sh
  clido list tasks -P 1
  ```

- Remove a project
  ```sh
  clido remove project 1
  ```

- Toggle task completion
  ```sh
  clido toggle 1
  ```

For detailed help, use the help command:
```sh
clido help
```

## Roadmap

- [x] Add task and project management
- [ ] Add priority levels for tasks
- [ ] Add reminders and notifications
- [ ] Add sub-tasks and sub-projects
- [ ] Add a TUI interface

See the [open issues](https://github.com/d4r1us-drk/clido/issues) for a full list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the GPLv3 License. See `LICENSE.txt` for more information.