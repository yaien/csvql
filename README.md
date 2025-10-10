# CSVQL

A high-performance CSV query tool written in Go, powered by DuckDB. Query CSV files using SQL syntax with automatic schema detection and type inference.

## Features

- 🚀 **Fast CSV Processing**: Built on DuckDB for optimal performance
- 🔍 **Automatic Schema Detection**: Intelligent type inference for columns
- 📊 **SQL Query Support**: Full SQL capabilities on CSV data
- 🕒 **Smart Date/Time Parsing**: Multiple date and time format support
- 📋 **Tabular Output**: Clean tabwriter-formatted table display
- 🎯 **Type Detection**: Automatic detection of integers, floats, booleans, dates, and strings
- 🌐 **REST API Server**: HTTP API for web applications and remote access

## Installation

### From Source

```bash
git clone https://github.com/yaien/csvql.git
cd csvql
go build -o csvql ./cmd/csvql
```

### Using Go Install

```bash
go install github.com/yaien/csvql/cmd/csvql@latest
```

## Quick Start

### Command Line Interface

```bash
# Query a CSV file (default: shows first 20 rows)
csvql data.csv

# Custom SQL query
csvql users.csv -q "SELECT * FROM data WHERE age > 25"

# Custom table name
csvql employees.csv -t employees -q "SELECT department, AVG(salary) as avg_salary FROM employees GROUP BY department"

# Aggregate data
csvql products.csv -q "SELECT COUNT(*), AVG(price) FROM data WHERE in_stock = true"
```

### REST API Server

```bash
# Start REST API server with a CSV file
csvql serve employees.csv --port 8047

# Upload additional CSV files and query via HTTP
curl -X POST http://localhost:8047/csvql/submit/ \
  -F "file=@users.csv" \
  -F "table=users"

# Execute SQL queries
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM data WHERE age > 25"}'
```

## Usage

### Basic Query (Default)
```bash
# Shows first 20 rows with inferred schema
csvql data.csv
```

### Custom SQL Queries
```bash
# Filter data
csvql users.csv -q "SELECT name, age FROM data WHERE age > 30"

# Aggregations
csvql products.csv -q "SELECT COUNT(*), AVG(price) FROM data WHERE in_stock = true"

# Date/Time queries
csvql events.csv -q "SELECT * FROM data WHERE event_date > '2024-01-01'"

# Custom table name
csvql employees.csv -t emp -q "SELECT * FROM emp WHERE department = 'Engineering'"
```

## Supported Data Types

The tool automatically detects and converts the following data types:

### Boolean Values
- `true`, `false` (case-insensitive only)

### Date Formats
- ISO format: `2006-01-02`
- US format: `01/02/2006`, `02/01/2006`
- European format: `02-01-2006`, `01-02-2006`

### Time Formats
- 24-hour: `15:04:05`, `15:04`
- 12-hour: `3:04PM`, `3:04:05PM`

### Timestamp Formats
- ISO 8601: `2006-01-02T15:04:05Z07:00`
- RFC 3339: `2006-01-02T15:04:05Z`
- Custom: `2006-01-02 15:04:05`

### Numeric Types
- **INTEGER**: Whole numbers
- **FLOAT**: Decimal numbers
- **VARCHAR**: Text strings

## Example Data

The repository includes sample CSV files in the `testdata/` directory:

- `testdata.csv` - Employee data with mixed types
- `numeric.csv` - Product catalog with prices
- `datetime.csv` - Event scheduling data
- `simple.csv` - Student information

### Sample Query Results

```bash
$ csvql testdata/testdata.csv -q "SELECT name, age, salary FROM data WHERE department = 'Engineering'"
```

```
name            age     salary
John Doe        30      75000.5
Alice Brown     28      72000
Frank Miller    29      71500.5
Ivy Chen        33      77000
```

## Command Line Options

### Query Mode
```bash
csvql [flags] <file>

Flags:
  -q, --query string   SQL query to execute on the CSV data (default "SELECT * FROM data LIMIT 20;")
  -t, --table string   Name of the table for SQL queries (default "data")
  -h, --help          Help for csvql
```

### Server Mode
```bash
csvql serve [flags] <file.csv>

Flags:
  -p, --port string       Port to listen on (default "8047")
      --max-memory int    Maximum memory for multipart form parsing in bytes (default 1073741824)
      --dbpath string     Path to DuckDB database file (empty for in-memory)
      --table string      Table name to use for the CSV data (default "data")
  -h, --help             Help for csvql serve
```

### Examples

```bash
# View first 20 rows (default behavior)
csvql employees.csv

# Custom query with default table name "data"
csvql employees.csv -q "SELECT department, COUNT(*) FROM data GROUP BY department"

# Custom table name
csvql employees.csv -t emp -q "SELECT * FROM emp WHERE salary > 50000"

# Complex aggregation
csvql sales.csv -q "SELECT DATE(order_date) as day, SUM(amount) FROM data GROUP BY day ORDER BY day"

# Start REST API server with a CSV file
csvql serve employees.csv --port 8047 --table emp

# Start server with persistent database
csvql serve data.csv --dbpath ./my_database.db --port 9000
```

## REST API

CSVQL can run as a REST API server, providing HTTP access to CSV querying capabilities. See [REST_API.md](REST_API.md) for detailed documentation.

### Quick API Example

```bash
# Start server with a CSV file
csvql serve employees.csv

# Upload additional CSV files
curl -X POST http://localhost:8047/csvql/submit/ \
  -F "file=@users.csv" \
  -F "table=users"

# Get available schemas
curl http://localhost:8047/csvql/schemas/

# Execute SQL query
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT name, age FROM data WHERE department = ?",
    "params": ["Engineering"]
  }'
```

### API Endpoints

- `POST /csvql/submit/` - Upload and import CSV files
- `POST /csvql/query/` - Execute SQL queries  
- `GET /csvql/schemas/` - Get table schemas and metadata

### Response Format

```json
{
  "data": {
    "columns": ["name", "age", "department"],
    "rows": [
      ["John Doe", 30, "Engineering"],
      ["Jane Smith", 28, "Engineering"]
    ]
  },
  "meta": {
    "rowCount": 2,
    "query": "SELECT name, age, department FROM data WHERE department = 'Engineering'"
  },
  "executionTime": "2.5ms"
}
```

### CSV Parsing Options

The parser automatically handles:
- Different delimiter detection
- Header row identification
- Quote character handling
- Escape character processing

## Development

### Prerequisites

- Go 1.24.3 or later
- GCC (for CGO compilation)

### Building

```bash
# Build the application
go build -o csvql ./cmd/csvql

# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

### Testing

```bash
# Run all tests
go test -v ./...

# Test with sample data
./csvql testdata/testdata.csv -q "SELECT COUNT(*) FROM data"
```

## Architecture

```
csvql/
├── cmd/csvql/          # CLI application entry point
│   ├── main.go         # Main application entry
│   ├── root.go         # Root CLI command and query functionality
│   └── server.go       # REST API server command
├── csvqlserver/        # REST API server package
│   ├── server.go       # Server implementation and handlers
│   ├── routes.go       # HTTP routing configuration
│   └── utils.go        # Utility functions for responses
├── testdata/           # Sample CSV files
├── parser.go           # CSV parsing and schema detection
├── parser_test.go      # Unit tests
├── web.go             # Web utilities
└── go.mod             # Go module definition
```

### Core Components

1. **Parser**: Handles CSV reading and type inference
2. **Schema Detection**: Automatically determines column types
3. **DuckDB Integration**: Provides SQL query engine with single database instance
4. **CLI Interface**: Command-line tool with query and server modes
5. **REST API Server**: HTTP API for web applications and remote access
6. **Table Printer**: Formats output for terminal display

## Performance

CSVQL leverages DuckDB's columnar storage and vectorized execution for optimal performance:

- **Large Files**: Efficiently processes GB-sized CSV files
- **Complex Queries**: Supports joins, aggregations, and window functions
- **Memory Efficient**: Streaming processing for large datasets

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run tests: `go test ./...`
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## Quick Reference

### CLI Commands
```bash
# Query CSV file
csvql data.csv -q "SELECT * FROM data WHERE age > 25"

# Start REST API server
csvql serve data.csv --port 8047
```

### API Endpoints
```bash
# Upload CSV
curl -X POST http://localhost:8047/csvql/submit/ -F "file=@data.csv"

# Query data
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM data LIMIT 10"}'

# Get schemas
curl http://localhost:8047/csvql/schemas/
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [DuckDB](https://duckdb.org/) - High-performance analytical database
- [Apache Arrow](https://arrow.apache.org/) - Columnar memory format
- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go

## Related Projects

- [csvkit](https://csvkit.readthedocs.io/) - CSV processing toolkit in Python
- [q](https://github.com/harelba/q) - Run SQL directly on CSV files
- [textql](https://github.com/dinedal/textql) - Execute SQL against structured text
