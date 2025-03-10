# JFrog Home Assignment

This repository contains the implementation of a downloader that reads URLs from a CSV file, downloads the content, and writes it to files.

## Prerequisites

- Go 1.23 or later
- `make` for running build and test commands
- `mockgen` for generating mocks for unit tests

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/jfrog-home-assignment.git
    cd jfrog-home-assignment
    ```

2. Install dependencies:

    ```sh
    go mod download
    ```

3. Install `mockgen`:

    ```sh
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

4. Install `make`:

    On Ubuntu/Debian:

    ```sh
    sudo apt-get update
    sudo apt-get install make
    ```

## Generating Mocks

To generate mocks for the interfaces, run the following command:

```sh
make mocks
```

## Building the Binary

To build the binary, run the following command:

```sh
make build
```

This will create an executable named `home-assignment` in the `bin` directory.

## Running the Binary

To run the binary, use the following command:

```sh
./bin/home-assignment -csv-file path/to/your/csvfile.csv -out-dir path/to/output/directory
```

Replace `path/to/your/csvfile.csv` with the path to your CSV file and `path/to/output/directory` with the path to the directory where you want to save the downloaded files.

## Running Tests

To run the tests, use the following command:

```sh
make test
```

This will run all the tests in the repository and display detailed output.

## Project Structure

```
.
├── cmd
│   └── home-assignment
│       └── home-assignment.go
├── internal
│   ├── config
│   │   ├── config.go
│   │   ├── config_test.go
│   │   └── types.go
│   ├── csv-reader
│   │   ├── csv_reader.go
│   │   └── csv_reader_test.go
│   ├── downloader
│   │   ├── downloader.go
│   │   └── downloader_test.go
│   ├── file-writer
│   │   ├── file_writer.go
│   │   └── file_writer_test.go
│   └── types
│       ├── logger_stub.go
│       └── types.go
├── go.mod
├── go.sum
├── Makefile
└── ReadMe.md
```
