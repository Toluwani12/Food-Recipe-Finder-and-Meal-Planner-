package pkg

import (
	"Food/internal/errors"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
)

type ApiResponse struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message"`
	Code       int         `json:"code"`
	Error      string      `json:"errors,omitempty"`
}

type Response struct {
	Message interface{} `json:"message,omitempty"`
	Err     interface{} `json:"errors,omitempty"`
}

func Render(w http.ResponseWriter, r *http.Request, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	switch v := res.(type) {
	case *errors.ErrResponse:
		w.WriteHeader(v.HTTPStatusCode)
		json.NewEncoder(w).Encode(Response{Message: v.StatusText, Err: v.ErrorText})
	case error:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal Server Error", Err: v.Error()})
	default:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Message: v, Err: nil})
	}
}

// Pagination holds pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// NewPagination creates a new Pagination instance
func NewPagination(page, pageSize, totalItems int) *Pagination {
	totalPages := (totalItems + pageSize - 1) / pageSize
	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
}

// ApplyPagination applies pagination to a SQL query
func ApplyPagination(query string, page, pageSize int) string {
	offset := (page - 1) * pageSize
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", query, pageSize, offset)
}

// ParsePaginationParams parses pagination parameters from URL query
func ParsePaginationParams(query url.Values) (int, int, error) {
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	return page, pageSize, nil
}

// ParseParams parses pagination parameters from URL query
func ParseParams(query url.Values) (int, int, error) {
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	return page, pageSize, nil
}

// ApplyToQuery applies pagination to a SQL query
func ApplyToQuery(query string, page, pageSize int) string {
	offset := (page - 1) * pageSize
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", query, pageSize, offset)
}

type ExecContext func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

func Log(service, repo, userID string, otherFields ...log.Fields) *log.Entry {
	fields := log.Fields{
		"service": service,
		"repo":    repo,
		"user_id": userID,
	}
	for _, f := range otherFields {
		for k, v := range f {
			fields[k] = v
		}
	}
	return log.WithFields(fields)
}
