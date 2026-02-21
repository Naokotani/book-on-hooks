package requestx

import (
	"errors"
	"strconv"
	"strings"
)

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
	s := strings.ToLower(strings.TrimSpace(v))
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
