/*
 * clido.cpp
 *
 * Copyright 2024 Darius Drake
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * SPDX-License-Identifier: GPL-3.0-only
 */

#include <iostream>
#include <getopt.h>

#define NAME "clido"
#define VERSION 0.1

void displayHelp();
void displayVersion();

int main(int argc, char** argv) {
    int option;
    static const char* short_opts = "hv";
    static struct option long_opts[] = {
        {"help",    no_argument, 0, 'h'},
        {"version", no_argument, 0, 'v'},
        {NULL, 0, NULL, 0}
    };

    while ((option = getopt_long(argc, argv, short_opts, long_opts, NULL)) != -1) {
        switch (option) {
            case 'h':
                displayHelp();
                break;
            case 'v':
                displayVersion();
                break;
            case '?':
            default:
                fprintf(stderr, "Use '-h, --help' for help.\n");
        }
    }

    return 0;
}

void displayHelp() {
    std::cout << "Usage: " << NAME << " [OPTIONS]" << std::endl;
    std::cout << "Manage your tasks and projects with ease." << std::endl << std::endl;
    std::cout << "Options:" << std::endl;
    std::cout << "\t-h, --help           Display this help message and exit." << std::endl;
    std::cout << "\t-v, --version        Display version and exit." << std::endl;
}

void displayVersion() {
    printf("%s v%.1f\n", NAME, VERSION);
}
