# CSVQL

A Simple tool that imports all csv data into a duckdb database

## Installation

### From Source

```bash
git clone https://github.com/yaien/csvql.git
cd csvql
go install ./cmd/csvql
```

### Using Go Install

```bash
go install github.com/yaien/csvql/cmd/csvql@latest
```

## Quick Start

### Command Line Interface

```bash
# Import data.csv into a table with name "datatable" from the "data.db" database
csvql data.csv datatable data.db
```
