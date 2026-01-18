package util

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// ParseID extracts the ID from the end of the request URL path
func ParseID(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 0 {
		return 0, errors.New("invalid path")
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return id, nil
}
