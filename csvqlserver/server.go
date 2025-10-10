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

type Server struct {
	db        *sql.DB
	schemas   []*csvql.Schema
	maxMemory int64
}

func New(db *sql.DB, schemas []*csvql.Schema) *Server {
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
	s.schemas = append(s.schemas, schema)

	Send(w, http.StatusOK, H{"status": "success", "table": tablename})

}

func (s *Server) Query(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req SQLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Send(w, http.StatusBadRequest, SQLResponse{
			Error: fmt.Sprintf("invalid json request: %v", err),
			Meta:  &QueryMeta{Duration: time.Since(start).String()},
		})
		return
	}

	// Execute query on the single database
	columns, rows, err := csvql.Query(s.db, req.Query, req.Params...)
	if err != nil {
		Send(w, http.StatusBadRequest, SQLResponse{
			Error: fmt.Sprintf("Query execution error: %v", err),
			Meta:  &QueryMeta{Query: req.Query, Params: req.Params, Duration: time.Since(start).String()},
		})
		return
	}

	Send(w, http.StatusOK, SQLResponse{
		Data: &QueryData{
			Columns: columns,
			Rows:    rows,
		},
		Meta: &QueryMeta{
			RowCount: len(rows),
			Query:    req.Query,
			Params:   req.Params,
			Duration: time.Since(start).String(),
		},
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
