package validator

import "testing"

func TestValidatorValid(t *testing.T) {
	v := &Validator{}
	if !v.Valid() {
		t.Fatalf("expected empty validator to be valid")
	}

	v.AddFieldError("name", "required")
	if v.Valid() {
		t.Fatalf("expected validator with errors to be invalid")
	}
}

func TestAddFieldError(t *testing.T) {
	v := &Validator{}
	v.AddFieldError("name", "required")
	if got := v.FieldErrors["name"]; got != "required" {
		t.Fatalf("expected first error message to be kept, got %q", got)
	}

	v.AddFieldError("name", "different")
	if got := v.FieldErrors["name"]; got != "required" {
		t.Fatalf("expected existing field error not to be overwritten, got %q", got)
	}
}

func TestCheckField(t *testing.T) {
	v := &Validator{}
	v.CheckField(true, "name", "required")
	if !v.Valid() {
		t.Fatalf("expected no error when condition is true")
	}

	v.CheckField(false, "name", "required")
	if got := v.FieldErrors["name"]; got != "required" {
		t.Fatalf("expected error to be added when condition is false, got %q", got)
	}
}

func TestNotBlank(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		valid bool
	}{
		{name: "plain text", in: "abc", valid: true},
		{name: "surrounded by spaces", in: "  abc  ", valid: true},
		{name: "spaces only", in: "   ", valid: false},
		{name: "empty", in: "", valid: false},
		{name: "tabs and newlines", in: "\t\n", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotBlank(tt.in); got != tt.valid {
				t.Fatalf("NotBlank(%q) = %v, want %v", tt.in, got, tt.valid)
			}
		})
	}
}

func TestMaxChars(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		n     int
		valid bool
	}{
		{name: "under limit", in: "abcd", n: 5, valid: true},
		{name: "at limit", in: "abcde", n: 5, valid: true},
		{name: "over limit", in: "abcdef", n: 5, valid: false},
		{name: "unicode counts runes", in: "日本語", n: 3, valid: true},
		{name: "unicode over limit", in: "日本語", n: 2, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxChars(tt.in, tt.n); got != tt.valid {
				t.Fatalf("MaxChars(%q, %d) = %v, want %v", tt.in, tt.n, got, tt.valid)
			}
		})
	}
}

func TestPermittedInt(t *testing.T) {
	tests := []struct {
		name      string
		in        int
		permitted []int
		valid     bool
	}{
		{name: "present", in: 2, permitted: []int{1, 2, 3}, valid: true},
		{name: "missing", in: 4, permitted: []int{1, 2, 3}, valid: false},
		{name: "empty permitted", in: 1, permitted: []int{}, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermittedInt(tt.in, tt.permitted...); got != tt.valid {
				t.Fatalf("PermittedInt(%d, %v) = %v, want %v", tt.in, tt.permitted, got, tt.valid)
			}
		})
	}
}

func TestPositiveInt(t *testing.T) {
	tests := []struct {
		name  string
		in    int
		valid bool
	}{
		{name: "positive", in: 1, valid: true},
		{name: "zero", in: 0, valid: false},
		{name: "negative", in: -1, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PositiveInt(tt.in); got != tt.valid {
				t.Fatalf("PositiveInt(%d) = %v, want %v", tt.in, got, tt.valid)
			}
		})
	}
}

func TestValidDollarValue(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		valid bool
	}{
		{name: "integer no symbol", in: "10", valid: true},
		{name: "integer with symbol", in: "$10", valid: true},
		{name: "two decimals", in: "10.99", valid: true},
		{name: "two decimals with symbol", in: "$10.99", valid: true},
		{name: "one decimal invalid", in: "10.9", valid: false},
		{name: "three decimals invalid", in: "10.999", valid: false},
		{name: "alpha invalid", in: "abc", valid: false},
		{name: "negative invalid", in: "-10.00", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidDollarValue(tt.in); got != tt.valid {
				t.Fatalf("ValidDollarValue(%q) = %v, want %v", tt.in, got, tt.valid)
			}
		})
	}
}
