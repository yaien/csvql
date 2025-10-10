# CSVQL

A high-performance CSV query tool written in Go, powered by DuckDB. Query CSV files using SQL syntax with automatic schema detection and type inference.

## Features

- 🚀 **Fast CSV Processing**: Built on DuckDB for optimal performance
- 🔍 **Automatic Schema Detection**: Intelligent type inference for columns
- 📊 **SQL Query Support**: Full SQL capabilities on CSV data
- 🕒 **Smart Date/Time Parsing**: Multiple date and time format support
- 📋 **Tabular Output**: Clean tabwriter-formatted table display
- 🎯 **Type Detection**: Automatic detection of integers, floats, booleans, dates, and strings

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

```bash
csvql [flags] <file>

Flags:
  -q, --query string   SQL query to execute on the CSV data (default "SELECT * FROM data LIMIT 20;")
  -t, --table string   Name of the table for SQL queries (default "data")
  -h, --help          Help for csvql
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
├── internal/
│   └── printer/        # Table formatting and output
├── testdata/           # Sample CSV files
├── parser.go           # CSV parsing and schema detection
├── parser_test.go      # Unit tests
└── go.mod             # Go module definition
```

### Core Components

1. **Parser**: Handles CSV reading and type inference
2. **Schema Detection**: Automatically determines column types
3. **DuckDB Integration**: Provides SQL query engine
4. **Table Printer**: Formats output for terminal display

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
