package requestx

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type ParseError struct {
	Field string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("invalid %s: %v", e.Field, e.Err)
}

func NormalizeText(v string) string {
	return strings.TrimSpace(v)
}

func NormalizeLower(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func QueryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func ParsePathID(r *http.Request, field string) (int64, error) {
	raw := httprouter.ParamsFromContext(r.Context()).ByName(field)
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || n <= 0 {
		return 0, &ParseError{Field: field, Err: errors.New("must be a positive integer")}
	}
	return n, nil
}

func ParseOptionalBool(v string) (bool, error) {
	if v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}

func ParseOptionalInt64(v string) (int64, error) {
	if v == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n <= 0 {
		return 0, errors.New("invalid int64")
	}
	return n, nil
}

func ParseOptionalSource(v string) (string, error) {
	s := NormalizeLower(v)
	if s == "" {
		return "", nil
	}
	if len(s) > 64 {
		return "", errors.New("source too long")
	}
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' {
			continue
		}
		return "", errors.New("invalid source")
	}
	return s, nil
}

func ParseOptionalQueryBool(r *http.Request, key string) (bool, error) {
	v, err := ParseOptionalBool(QueryString(r, key))
	if err != nil {
		return false, &ParseError{Field: key, Err: err}
	}
	return v, nil
}

func ParseOptionalQueryInt64(r *http.Request, key string) (int64, error) {
	v, err := ParseOptionalInt64(QueryString(r, key))
	if err != nil {
		return 0, &ParseError{Field: key, Err: err}
	}
	return v, nil
}

func ParseOptionalQuerySource(r *http.Request, key string) (string, error) {
	v, err := ParseOptionalSource(QueryString(r, key))
	if err != nil {
		return "", &ParseError{Field: key, Err: err}
	}
	return v, nil
}

func IsAdminSource(source string) bool {
	return strings.HasPrefix(source, "admin-")
}
