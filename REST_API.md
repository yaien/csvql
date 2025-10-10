# CSVQL REST API Documentation

This document provides comprehensive examples of how to use the CSVQL REST API server.

The server uses a single DuckDB database instance where you can load multiple CSV files as tables and query them using SQL.

## Starting the Server

```bash
# Start server with a CSV file (default port 8047)
csvql serve employees.csv

# Start server with custom settings
csvql serve data.csv --port 9000 --table mydata --dbpath ./database.db

# Start with increased memory limit for large files
csvql serve large_file.csv --max-memory 2147483648  # 2GB
```

## Server Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8047` | Port to listen on |
| `--table` | `data` | Table name for the initial CSV file |
| `--dbpath` | (empty) | Path to persistent DuckDB file (in-memory if empty) |
| `--max-memory` | `1073741824` | Max memory for file uploads (1GB) |

## API Endpoints

### 1. Upload CSV Files

**Endpoint:** `POST /csvql/submit/`

Upload and import CSV files into the database.

**Example:**
```bash
# Upload a CSV file with automatic table name (filename without extension)
curl -X POST http://localhost:8047/csvql/submit/ \
  -F "file=@users.csv"

# Upload with custom table name
curl -X POST http://localhost:8047/csvql/submit/ \
  -F "file=@data.csv" \
  -F "table=employees"
```

**Success Response:**
```json
{
  "status": "success",
  "table": "users"
}
```

**Error Response:**
```json
{
  "error": "failed to parse csv file: invalid format"
}
```

### 2. Execute SQL Queries

**Endpoint:** `POST /csvql/query/`

Execute SQL queries against the loaded tables.

**Request Body:**
```json
{
  "query": "SELECT * FROM data WHERE age > ?",
  "params": [25]
}
```

**Example:**
```bash
# Simple query
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM data LIMIT 10"}'

# Parameterized query
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT name, age FROM users WHERE department = ? AND age > ?",
    "params": ["Engineering", 25]
  }'

# Join multiple tables
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT u.name, e.salary FROM users u JOIN employees e ON u.id = e.user_id"
  }'

# Aggregation query
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT department, COUNT(*) as count, AVG(salary) as avg_salary FROM data GROUP BY department"
  }'
```

**Success Response:**
```json
{
  "data": {
    "columns": ["name", "age", "department"],
    "rows": [
      ["John Doe", 30, "Engineering"],
      ["Jane Smith", 28, "Engineering"],
      ["Bob Wilson", 35, "Engineering"]
    ]
  },
  "meta": {
    "rowCount": 3,
    "query": "SELECT name, age, department FROM users WHERE department = ?",
    "params": ["Engineering"]
  },
  "executionTime": "2.5ms"
}
```

**Error Response:**
```json
{
  "error": "Query execution error: table 'nonexistent' does not exist",
  "meta": {
    "query": "SELECT * FROM nonexistent"
  },
  "executionTime": "1.2ms"
}
```

### 3. Get Table Schemas

**Endpoint:** `GET /csvql/schemas/`

Retrieve schema information for all loaded tables.

**Example:**
```bash
curl http://localhost:8047/csvql/schemas/
```

**Response:**
```json
{
  "users": {
    "name": "users",
    "columns": [
      {"Name": "id", "Type": "INTEGER"},
      {"Name": "name", "Type": "VARCHAR"},
      {"Name": "email", "Type": "VARCHAR"},
      {"Name": "age", "Type": "INTEGER"},
      {"Name": "department", "Type": "VARCHAR"}
    ]
  },
  "employees": {
    "name": "employees", 
    "columns": [
      {"Name": "user_id", "Type": "INTEGER"},
      {"Name": "salary", "Type": "DOUBLE"},
      {"Name": "hire_date", "Type": "DATE"},
      {"Name": "is_active", "Type": "BOOLEAN"}
    ]
  }
}
```

## Error Handling

### Invalid SQL Query
```bash
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM nonexistent_table"}'
```

**Response:**
```json
{
  "error": "Query execution error: no such table: nonexistent_table",
  "meta": {
    "query": "SELECT * FROM nonexistent_table",
    "params": [],
  },
  "executionTime": "1.2ms"
}
```

### Invalid JSON Request
```bash
curl -X POST http://localhost:8047/csvql/query/ \
  -H "Content-Type: application/json" \
  -d '{"invalid": json}'
```

**Response:**
```json
{
  "error": "invalid json request: invalid character 'j' looking for beginning of value",
  "meta": {},
  "executionTime": "0.1ms"
}
```

## cURL Examples

```bash
#!/bin/bash

# Set base URL
BASE_URL="http://localhost:8047"

# Upload CSV file
echo "Uploading CSV file..."
curl -X POST "$BASE_URL/csvql/submit/" \
  -F "file=@employees.csv" \
  -F "table=staff"

echo -e "\n\nQuerying data..."
# Simple query
curl -X POST "$BASE_URL/csvql/query/" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT COUNT(*) as total_employees FROM staff"
  }'

echo -e "\n\nParameterized query..."
# Query with parameters
curl -X POST "$BASE_URL/csvql/query/" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT name, salary FROM staff WHERE department = ? AND salary > ?",
    "params": ["Engineering", 60000]
  }'

echo -e "\n\nGetting schemas..."
# Get schemas
curl "$BASE_URL/csvql/schemas/"
```

## Data Type Support

The API automatically detects and handles the following data types:

| Type | Examples | SQL Type |
|------|----------|----------|
| **INTEGER** | `1`, `42`, `-10` | `INTEGER` |
| **FLOAT** | `3.14`, `-2.5`, `1.0` | `DOUBLE` |
| **BOOLEAN** | `true`, `false` (case-insensitive) | `BOOLEAN` |
| **DATE** | `2024-01-15`, `01/15/2024` | `DATE` |
| **TIME** | `14:30:00`, `2:30PM` | `TIME` |
| **TIMESTAMP** | `2024-01-15T14:30:00Z` | `TIMESTAMP` |
| **VARCHAR** | Any text string | `VARCHAR` |

## Performance Considerations

### Large File Uploads
- Use `--max-memory` flag to increase memory limit for large CSV files
- Consider splitting very large files into smaller chunks
- Use persistent database (`--dbpath`) for better performance with multiple sessions

### Query Optimization
- Use parameterized queries for better performance and security
- Create indexes on frequently queried columns (via SQL DDL)
- Use `LIMIT` clauses for large result sets

### Memory Usage
- In-memory database (default) is faster but limited by RAM
- Persistent database (`--dbpath`) can handle larger datasets
- Monitor memory usage with large datasets

## Security Considerations

### Input Validation
- Use parameterized queries to prevent SQL injection
- Validate file types and sizes before upload
- Consider implementing authentication for production use

### Network Security
- Run server behind reverse proxy (nginx, Apache) in production
- Use HTTPS in production environments
- Implement rate limiting and request size limits

## CORS Support

The server includes CORS headers for web browser compatibility.

## Troubleshooting

### Common Issues

1. **"Connection refused"**
   - Ensure server is running: `csvql serve data.csv`
   - Check port availability: `netstat -tlnp | grep 8047`

2. **"Table does not exist"**
   - Verify table was uploaded successfully
   - Check available tables: `curl http://localhost:8047/csvql/schemas/`

3. **"File too large"**
   - Increase memory limit: `--max-memory 2147483648`
   - Split large files into smaller chunks

4. **"Invalid CSV format"**
   - Ensure CSV file has proper headers
   - Check for encoding issues (use UTF-8)
   - Verify delimiter and quote characters

### Debug Mode
```bash
# Start server with verbose logging
csvql serve data.csv --port 8047 2>&1 | tee server.log
```