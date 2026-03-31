package grumble

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// TestStrToInt
// ---------------------------------------------------------------------------

func TestStrToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "decimal positive", input: "42", want: 42},
		{name: "decimal negative", input: "-1", want: -1},
		{name: "hex", input: "0x1F", want: 31},
		{name: "octal", input: "0o17", want: 15},
		{name: "binary", input: "0b1010", want: 10},
		{name: "zero", input: "0", want: 0},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "abc", wantErr: true},
		{name: "float-like", input: "3.14", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToInt(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToInt(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToInt(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToInt(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToInt64
// ---------------------------------------------------------------------------

func TestStrToInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{name: "decimal positive", input: "42", want: 42},
		{name: "decimal negative", input: "-99999", want: -99999},
		{name: "large positive", input: "9223372036854775807", want: math.MaxInt64},
		{name: "large negative", input: "-9223372036854775808", want: math.MinInt64},
		{name: "hex", input: "0xFF", want: 255},
		{name: "binary", input: "0b11111111", want: 255},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "xyz", wantErr: true},
		{name: "overflow", input: "9223372036854775808", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToInt64(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToInt64(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToInt64(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToInt64(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToInt32
// ---------------------------------------------------------------------------

func TestStrToInt32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int32
		wantErr bool
	}{
		{name: "decimal positive", input: "1000", want: 1000},
		{name: "decimal negative", input: "-1000", want: -1000},
		{name: "max int32", input: "2147483647", want: math.MaxInt32},
		{name: "min int32", input: "-2147483648", want: math.MinInt32},
		{name: "hex", input: "0x7F", want: 127},
		{name: "overflow", input: "2147483648", wantErr: true},
		{name: "underflow", input: "-2147483649", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "nope", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToInt32(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToInt32(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToInt32(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToInt32(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToInt16
// ---------------------------------------------------------------------------

func TestStrToInt16(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int16
		wantErr bool
	}{
		{name: "decimal positive", input: "300", want: 300},
		{name: "decimal negative", input: "-300", want: -300},
		{name: "max int16", input: "32767", want: math.MaxInt16},
		{name: "min int16", input: "-32768", want: math.MinInt16},
		{name: "hex", input: "0xFF", want: 255},
		{name: "overflow", input: "32768", wantErr: true},
		{name: "underflow", input: "-32769", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "bad", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToInt16(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToInt16(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToInt16(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToInt16(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToInt8
// ---------------------------------------------------------------------------

func TestStrToInt8(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int8
		wantErr bool
	}{
		{name: "decimal positive", input: "42", want: 42},
		{name: "decimal negative", input: "-42", want: -42},
		{name: "max int8", input: "127", want: 127},
		{name: "min int8", input: "-128", want: -128},
		{name: "hex", input: "0x7F", want: 127},
		{name: "binary", input: "0b1010", want: 10},
		{name: "overflow", input: "128", wantErr: true},
		{name: "underflow", input: "-129", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToInt8(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToInt8(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToInt8(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToInt8(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToUint
// ---------------------------------------------------------------------------

func TestStrToUint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint
		wantErr bool
	}{
		{name: "decimal", input: "42", want: 42},
		{name: "zero", input: "0", want: 0},
		{name: "hex", input: "0xFF", want: 255},
		{name: "octal", input: "0o77", want: 63},
		{name: "binary", input: "0b1100", want: 12},
		{name: "negative", input: "-1", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToUint(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToUint(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToUint(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToUint(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToUint64
// ---------------------------------------------------------------------------

func TestStrToUint64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{name: "decimal", input: "42", want: 42},
		{name: "large value", input: "18446744073709551615", want: math.MaxUint64},
		{name: "hex", input: "0xFFFF", want: 65535},
		{name: "binary", input: "0b1111", want: 15},
		{name: "negative", input: "-1", wantErr: true},
		{name: "overflow", input: "18446744073709551616", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "nope", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToUint64(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToUint64(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToUint64(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToUint64(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToUint32
// ---------------------------------------------------------------------------

func TestStrToUint32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint32
		wantErr bool
	}{
		{name: "decimal", input: "1000", want: 1000},
		{name: "max uint32", input: "4294967295", want: math.MaxUint32},
		{name: "hex", input: "0xFFFF", want: 65535},
		{name: "zero", input: "0", want: 0},
		{name: "overflow", input: "4294967296", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "bad", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToUint32(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToUint32(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToUint32(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToUint32(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToUint16
// ---------------------------------------------------------------------------

func TestStrToUint16(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint16
		wantErr bool
	}{
		{name: "decimal", input: "500", want: 500},
		{name: "max uint16", input: "65535", want: math.MaxUint16},
		{name: "hex", input: "0xFF", want: 255},
		{name: "zero", input: "0", want: 0},
		{name: "overflow", input: "65536", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "xyz", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToUint16(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToUint16(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToUint16(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToUint16(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestStrToUint8
// ---------------------------------------------------------------------------

func TestStrToUint8(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint8
		wantErr bool
	}{
		{name: "decimal", input: "42", want: 42},
		{name: "max uint8", input: "255", want: math.MaxUint8},
		{name: "hex", input: "0xFF", want: 255},
		{name: "binary", input: "0b11111111", want: 255},
		{name: "zero", input: "0", want: 0},
		{name: "overflow", input: "256", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "non-numeric", input: "nah", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToUint8(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("strToUint8(%q): expected error, got %d", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("strToUint8(%q): unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("strToUint8(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
