# go-licenser
CLI tool written in Go for gathering licenses of dependencies in your Go project.

## Basic usage
You can gather all licenses in your project into one big text file via `licenser gather <path to module> <output file>` command.
Tool will automatically gather all dependencies that are collected by `go vendor` and extract licenses from them to one output text file.
