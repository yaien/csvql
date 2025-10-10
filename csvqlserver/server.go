package csvqlserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/yaien/csvql"
)

type SQLRequest struct {
	Query  string `json:"query"`
	Params []any  `json:"params,omitempty"`
}

type SQLResponse struct {
	Data          *QueryData `json:"data,omitempty"`
	Meta          *QueryMeta `json:"meta"`
	Error         string     `json:"error,omitempty"`
	ExecutionTime string     `json:"executionTime"`
}

type QueryData struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
}

type QueryMeta struct {
	RowCount int    `json:"rowCount"`
	Query    string `json:"query"`
	Params   []any  `json:"params"`
}

type Server struct {
	db        *sql.DB
	schemas   map[string]*csvql.Schema
	maxMemory int64
}

func New(db *sql.DB, schemas map[string]*csvql.Schema) *Server {
	if schemas == nil {
		schemas = make(map[string]*csvql.Schema)
	}

	return &Server{
		db:        db,
		schemas:   schemas,
		maxMemory: 1024 * 1024 * 1024, // 1GB
	}
}

func (s *Server) SetMaxMemory(max int64) { s.maxMemory = max }

func (s *Server) Submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(s.maxMemory)
	if err != nil {
		Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}

	defer func() { _ = file.Close() }()

	tablename := r.FormValue("table")
	if tablename == "" {
		tablename = strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	}

	// Parse CSV file
	schema, rows, err := csvql.Parse(file, tablename)
	if err != nil {
		Error(w, fmt.Sprintf("failed to parse csv file: %v", err), http.StatusBadRequest)
		return
	}

	err = csvql.ImportOnDB(s.db, schema, rows)
	if err != nil {
		Error(w, fmt.Sprintf("failed to import csv file into database: %v", err), http.StatusInternalServerError)
		return
	}

	// Store schema for API metadata
	s.schemas[tablename] = schema

	Send(w, http.StatusOK, H{"status": "success", "table": tablename})

}

func (s *Server) Query(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req SQLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Send(w, http.StatusBadRequest, SQLResponse{
			Error:         fmt.Sprintf("invalid json request: %v", err),
			ExecutionTime: time.Since(start).String(),
			Meta:          &QueryMeta{},
		})
		return
	}

	// Execute query on the single database
	rows, err := s.db.Query(req.Query, req.Params...)
	if err != nil {
		Send(w, http.StatusBadRequest, SQLResponse{
			Error:         fmt.Sprintf("Query execution error: %v", err),
			ExecutionTime: time.Since(start).String(),
			Meta:          &QueryMeta{Query: req.Query, Params: req.Params},
		})
		return
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		Send(w, http.StatusInternalServerError, SQLResponse{
			Error:         fmt.Sprintf("failed to get columns: %v", err),
			ExecutionTime: time.Since(start).String(),
			Meta:          &QueryMeta{Query: req.Query, Params: req.Params},
		})
		return
	}

	// Collect results
	var results [][]any
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			Send(w, http.StatusInternalServerError, SQLResponse{
				Error:         fmt.Sprintf("failed to scan row: %v", err),
				ExecutionTime: time.Since(start).String(),
				Meta:          &QueryMeta{Query: req.Query, Params: req.Params},
			})
			return
		}

		results = append(results, values)
	}

	if err = rows.Err(); err != nil {
		Send(w, http.StatusInternalServerError, SQLResponse{
			Error:         fmt.Sprintf("row iteration error: %v", err),
			ExecutionTime: time.Since(start).String(),
			Meta:          &QueryMeta{Query: req.Query, Params: req.Params},
		})
		return
	}

	Send(w, http.StatusOK, SQLResponse{
		Data: &QueryData{
			Columns: columns,
			Rows:    results,
		},
		Meta: &QueryMeta{
			RowCount: len(results),
			Query:    req.Query,
			Params:   req.Params,
		},
		ExecutionTime: time.Since(start).String(),
	})

}

func (s *Server) Schemas(w http.ResponseWriter, r *http.Request) {
	Send(w, http.StatusOK, s.schemas)
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	Send(w, http.StatusOK, H{
		"status":    "healthy",
		"tables":    len(s.schemas),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
