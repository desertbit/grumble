package grumble

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// mustNotPanic fails the test if the function panics.
func mustNotPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	fn()
}

// mustPanic fails the test if the function does NOT panic.
func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic but none occurred")
		}
	}()
	fn()
}

// ---------------------------------------------------------------------------
// TestArgRegistration
// ---------------------------------------------------------------------------

func TestArgRegistration(t *testing.T) {
	// Each scalar type followed by a list type would panic (list must be last).
	// So we register each type in its own Args instance and just confirm no
	// panic and that empty() returns false afterwards.

	types := []struct {
		label    string
		register func(a *Args)
	}{
		{"String", func(a *Args) { a.String("s", "help") }},
		{"StringList", func(a *Args) { a.StringList("sl", "help") }},
		{"Bool", func(a *Args) { a.Bool("b", "help") }},
		{"BoolList", func(a *Args) { a.BoolList("bl", "help") }},
		{"Int", func(a *Args) { a.Int("i", "help") }},
		{"IntList", func(a *Args) { a.IntList("il", "help") }},
		{"Int64", func(a *Args) { a.Int64("i64", "help") }},
		{"Int64List", func(a *Args) { a.Int64List("i64l", "help") }},
		{"Uint", func(a *Args) { a.Uint("u", "help") }},
		{"UintList", func(a *Args) { a.UintList("ul", "help") }},
		{"Uint64", func(a *Args) { a.Uint64("u64", "help") }},
		{"Uint64List", func(a *Args) { a.Uint64List("u64l", "help") }},
		{"Float64", func(a *Args) { a.Float64("f", "help") }},
		{"Float64List", func(a *Args) { a.Float64List("fl", "help") }},
		{"Duration", func(a *Args) { a.Duration("d", "help") }},
		{"DurationList", func(a *Args) { a.DurationList("dl", "help") }},
	}

	for _, tc := range types {
		t.Run(tc.label, func(t *testing.T) {
			var a Args
			if !a.empty() {
				t.Fatal("expected empty before registration")
			}
			mustNotPanic(t, func() { tc.register(&a) })
			if a.empty() {
				t.Fatal("expected non-empty after registration")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestArgStringParse
// ---------------------------------------------------------------------------

func TestArgStringParse(t *testing.T) {
	var a Args
	a.String("name", "a name")

	res := make(ArgMap)
	rest, err := a.parse([]string{"hello", "extra"}, res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.String("name") != "hello" {
		t.Fatalf("expected 'hello', got '%s'", res.String("name"))
	}
	if len(rest) != 1 || rest[0] != "extra" {
		t.Fatalf("expected [extra], got %v", rest)
	}
}

// ---------------------------------------------------------------------------
// TestArgStringListParse
// ---------------------------------------------------------------------------

func TestArgStringListParse(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var a Args
		a.StringList("names", "list of names")

		res := make(ArgMap)
		rest, err := a.parse([]string{"a", "b", "c"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.StringList("names")
		if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
			t.Fatalf("expected [a b c], got %v", got)
		}
		if len(rest) != 0 {
			t.Fatalf("expected no remaining args, got %v", rest)
		}
	})

	t.Run("default", func(t *testing.T) {
		var a Args
		a.StringList("names", "list of names", Default([]string{"x"}))

		res := make(ArgMap)
		rest, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.StringList("names")
		if len(got) != 1 || got[0] != "x" {
			t.Fatalf("expected [x], got %v", got)
		}
		if !res["names"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
		if len(rest) != 0 {
			t.Fatalf("expected no remaining args, got %v", rest)
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgBoolParse
// ---------------------------------------------------------------------------

func TestArgBoolParse(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var a Args
		a.Bool("flag", "a flag")

		res := make(ArgMap)
		_, err := a.parse([]string{"true"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Bool("flag") != true {
			t.Fatal("expected true")
		}
	})

	t.Run("false", func(t *testing.T) {
		var a Args
		a.Bool("flag", "a flag")

		res := make(ArgMap)
		_, err := a.parse([]string{"false"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Bool("flag") != false {
			t.Fatal("expected false")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var a Args
		a.Bool("flag", "a flag")

		res := make(ArgMap)
		_, err := a.parse([]string{"abc"}, res)
		if err == nil {
			t.Fatal("expected error for invalid bool")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgBoolListParse
// ---------------------------------------------------------------------------

func TestArgBoolListParse(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var a Args
		a.BoolList("flags", "some flags")

		res := make(ArgMap)
		_, err := a.parse([]string{"true", "false", "true"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.BoolList("flags")
		if len(got) != 3 || got[0] != true || got[1] != false || got[2] != true {
			t.Fatalf("expected [true false true], got %v", got)
		}
	})

	t.Run("default", func(t *testing.T) {
		var a Args
		a.BoolList("flags", "some flags", Default([]bool{false, true}))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.BoolList("flags")
		if len(got) != 2 || got[0] != false || got[1] != true {
			t.Fatalf("expected [false true], got %v", got)
		}
		if !res["flags"].IsDefault {
			t.Fatal("expected IsDefault")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgIntParse
// ---------------------------------------------------------------------------

func TestArgIntParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Int("count", "a count")

		res := make(ArgMap)
		_, err := a.parse([]string{"42"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Int("count") != 42 {
			t.Fatalf("expected 42, got %d", res.Int("count"))
		}
	})

	t.Run("negative", func(t *testing.T) {
		var a Args
		a.Int("count", "a count")

		res := make(ArgMap)
		_, err := a.parse([]string{"-7"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Int("count") != -7 {
			t.Fatalf("expected -7, got %d", res.Int("count"))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var a Args
		a.Int("count", "a count")

		res := make(ArgMap)
		_, err := a.parse([]string{"notanumber"}, res)
		if err == nil {
			t.Fatal("expected error for invalid int")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgIntListParse
// ---------------------------------------------------------------------------

func TestArgIntListParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.IntList("nums", "numbers")

		res := make(ArgMap)
		_, err := a.parse([]string{"1", "2", "3"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.IntList("nums")
		if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
			t.Fatalf("expected [1 2 3], got %v", got)
		}
	})

	t.Run("invalid element", func(t *testing.T) {
		var a Args
		a.IntList("nums", "numbers")

		res := make(ArgMap)
		_, err := a.parse([]string{"1", "abc"}, res)
		if err == nil {
			t.Fatal("expected error for invalid int in list")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgInt64Parse
// ---------------------------------------------------------------------------

func TestArgInt64Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Int64("big", "a big number")

		res := make(ArgMap)
		_, err := a.parse([]string{"9223372036854775807"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Int64("big") != 9223372036854775807 {
			t.Fatalf("expected max int64, got %d", res.Int64("big"))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var a Args
		a.Int64("big", "a big number")

		res := make(ArgMap)
		_, err := a.parse([]string{"xyz"}, res)
		if err == nil {
			t.Fatal("expected error for invalid int64")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgInt64ListParse
// ---------------------------------------------------------------------------

func TestArgInt64ListParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Int64List("bigs", "big numbers")

		res := make(ArgMap)
		_, err := a.parse([]string{"100", "-200", "300"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.Int64List("bigs")
		if len(got) != 3 || got[0] != 100 || got[1] != -200 || got[2] != 300 {
			t.Fatalf("expected [100 -200 300], got %v", got)
		}
	})

	t.Run("invalid element", func(t *testing.T) {
		var a Args
		a.Int64List("bigs", "big numbers")

		res := make(ArgMap)
		_, err := a.parse([]string{"100", "nope"}, res)
		if err == nil {
			t.Fatal("expected error for invalid int64 in list")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgUintParse
// ---------------------------------------------------------------------------

func TestArgUintParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Uint("port", "a port number")

		res := make(ArgMap)
		_, err := a.parse([]string{"8080"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Uint("port") != 8080 {
			t.Fatalf("expected 8080, got %d", res.Uint("port"))
		}
	})

	t.Run("invalid negative", func(t *testing.T) {
		var a Args
		a.Uint("port", "a port number")

		res := make(ArgMap)
		_, err := a.parse([]string{"-1"}, res)
		if err == nil {
			t.Fatal("expected error for negative uint")
		}
	})

	t.Run("invalid string", func(t *testing.T) {
		var a Args
		a.Uint("port", "a port number")

		res := make(ArgMap)
		_, err := a.parse([]string{"abc"}, res)
		if err == nil {
			t.Fatal("expected error for invalid uint")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgUintListParse
// ---------------------------------------------------------------------------

func TestArgUintListParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.UintList("ports", "port numbers")

		res := make(ArgMap)
		_, err := a.parse([]string{"80", "443", "8080"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.UintList("ports")
		if len(got) != 3 || got[0] != 80 || got[1] != 443 || got[2] != 8080 {
			t.Fatalf("expected [80 443 8080], got %v", got)
		}
	})

	t.Run("invalid element", func(t *testing.T) {
		var a Args
		a.UintList("ports", "port numbers")

		res := make(ArgMap)
		_, err := a.parse([]string{"80", "-1"}, res)
		if err == nil {
			t.Fatal("expected error for invalid uint in list")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgUint64Parse
// ---------------------------------------------------------------------------

func TestArgUint64Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Uint64("bigu", "big unsigned")

		res := make(ArgMap)
		_, err := a.parse([]string{"18446744073709551615"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var expected uint64 = 18446744073709551615
		if res.Uint64("bigu") != expected {
			t.Fatalf("expected max uint64, got %d", res.Uint64("bigu"))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var a Args
		a.Uint64("bigu", "big unsigned")

		res := make(ArgMap)
		_, err := a.parse([]string{"-1"}, res)
		if err == nil {
			t.Fatal("expected error for negative uint64")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgUint64ListParse
// ---------------------------------------------------------------------------

func TestArgUint64ListParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Uint64List("vals", "values")

		res := make(ArgMap)
		_, err := a.parse([]string{"1000", "2000"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.Uint64List("vals")
		if len(got) != 2 || got[0] != 1000 || got[1] != 2000 {
			t.Fatalf("expected [1000 2000], got %v", got)
		}
	})

	t.Run("invalid element", func(t *testing.T) {
		var a Args
		a.Uint64List("vals", "values")

		res := make(ArgMap)
		_, err := a.parse([]string{"1000", "nope"}, res)
		if err == nil {
			t.Fatal("expected error for invalid uint64 in list")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgFloat64Parse
// ---------------------------------------------------------------------------

func TestArgFloat64Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Float64("ratio", "a ratio")

		res := make(ArgMap)
		_, err := a.parse([]string{"3.14"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.Float64("ratio")
		if got < 3.139 || got > 3.141 {
			t.Fatalf("expected ~3.14, got %f", got)
		}
	})

	t.Run("integer input", func(t *testing.T) {
		var a Args
		a.Float64("ratio", "a ratio")

		res := make(ArgMap)
		_, err := a.parse([]string{"42"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Float64("ratio") != 42.0 {
			t.Fatalf("expected 42.0, got %f", res.Float64("ratio"))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var a Args
		a.Float64("ratio", "a ratio")

		res := make(ArgMap)
		_, err := a.parse([]string{"notafloat"}, res)
		if err == nil {
			t.Fatal("expected error for invalid float64")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgFloat64ListParse
// ---------------------------------------------------------------------------

func TestArgFloat64ListParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Float64List("coords", "coordinates")

		res := make(ArgMap)
		_, err := a.parse([]string{"1.1", "2.2", "3.3"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.Float64List("coords")
		if len(got) != 3 {
			t.Fatalf("expected 3 elements, got %d", len(got))
		}
		if got[0] < 1.09 || got[0] > 1.11 {
			t.Fatalf("expected ~1.1, got %f", got[0])
		}
		if got[1] < 2.19 || got[1] > 2.21 {
			t.Fatalf("expected ~2.2, got %f", got[1])
		}
		if got[2] < 3.29 || got[2] > 3.31 {
			t.Fatalf("expected ~3.3, got %f", got[2])
		}
	})

	t.Run("invalid element", func(t *testing.T) {
		var a Args
		a.Float64List("coords", "coordinates")

		res := make(ArgMap)
		_, err := a.parse([]string{"1.1", "bad"}, res)
		if err == nil {
			t.Fatal("expected error for invalid float64 in list")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgDurationParse
// ---------------------------------------------------------------------------

func TestArgDurationParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.Duration("timeout", "a timeout")

		res := make(ArgMap)
		_, err := a.parse([]string{"5s"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Duration("timeout") != 5*time.Second {
			t.Fatalf("expected 5s, got %v", res.Duration("timeout"))
		}
	})

	t.Run("complex duration", func(t *testing.T) {
		var a Args
		a.Duration("timeout", "a timeout")

		res := make(ArgMap)
		_, err := a.parse([]string{"1h30m"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := 1*time.Hour + 30*time.Minute
		if res.Duration("timeout") != expected {
			t.Fatalf("expected %v, got %v", expected, res.Duration("timeout"))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var a Args
		a.Duration("timeout", "a timeout")

		res := make(ArgMap)
		_, err := a.parse([]string{"notaduration"}, res)
		if err == nil {
			t.Fatal("expected error for invalid duration")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgDurationListParse
// ---------------------------------------------------------------------------

func TestArgDurationListParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var a Args
		a.DurationList("timeouts", "timeouts")

		res := make(ArgMap)
		_, err := a.parse([]string{"1s", "2m", "3h"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.DurationList("timeouts")
		if len(got) != 3 {
			t.Fatalf("expected 3 elements, got %d", len(got))
		}
		if got[0] != 1*time.Second {
			t.Fatalf("expected 1s, got %v", got[0])
		}
		if got[1] != 2*time.Minute {
			t.Fatalf("expected 2m, got %v", got[1])
		}
		if got[2] != 3*time.Hour {
			t.Fatalf("expected 3h, got %v", got[2])
		}
	})

	t.Run("invalid element", func(t *testing.T) {
		var a Args
		a.DurationList("timeouts", "timeouts")

		res := make(ArgMap)
		_, err := a.parse([]string{"1s", "bad"}, res)
		if err == nil {
			t.Fatal("expected error for invalid duration in list")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgDefaultValues
// ---------------------------------------------------------------------------

func TestArgDefaultValues(t *testing.T) {
	t.Run("string default", func(t *testing.T) {
		var a Args
		a.String("name", "a name", Default("world"))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.String("name") != "world" {
			t.Fatalf("expected 'world', got '%s'", res.String("name"))
		}
		if !res["name"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("bool default", func(t *testing.T) {
		var a Args
		a.Bool("flag", "a flag", Default(true))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Bool("flag") != true {
			t.Fatal("expected true default")
		}
		if !res["flag"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("int default", func(t *testing.T) {
		var a Args
		a.Int("count", "a count", Default(10))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !res["count"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("int64 default", func(t *testing.T) {
		var a Args
		a.Int64("big", "big number", Default(int64(99)))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !res["big"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("uint default", func(t *testing.T) {
		var a Args
		a.Uint("port", "a port", Default(uint(8080)))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !res["port"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("uint64 default", func(t *testing.T) {
		var a Args
		a.Uint64("val", "a value", Default(uint64(999)))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !res["val"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("float64 default", func(t *testing.T) {
		var a Args
		a.Float64("ratio", "a ratio", Default(2.718))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !res["ratio"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("duration default", func(t *testing.T) {
		var a Args
		a.Duration("timeout", "a timeout", Default(5*time.Second))

		res := make(ArgMap)
		_, err := a.parse([]string{}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Duration("timeout") != 5*time.Second {
			t.Fatalf("expected 5s, got %v", res.Duration("timeout"))
		}
		if !res["timeout"].IsDefault {
			t.Fatal("expected IsDefault to be true")
		}
	})

	t.Run("provided value is not default", func(t *testing.T) {
		var a Args
		a.String("name", "a name", Default("world"))

		res := make(ArgMap)
		_, err := a.parse([]string{"alice"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.String("name") != "alice" {
			t.Fatalf("expected 'alice', got '%s'", res.String("name"))
		}
		if res["name"].IsDefault {
			t.Fatal("expected IsDefault to be false when value is provided")
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgMissingMandatory
// ---------------------------------------------------------------------------

func TestArgMissingMandatory(t *testing.T) {
	var a Args
	a.String("name", "a name")

	res := make(ArgMap)
	_, err := a.parse([]string{}, res)
	if err == nil {
		t.Fatal("expected error for missing mandatory argument")
	}
}

// ---------------------------------------------------------------------------
// TestArgMinMax
// ---------------------------------------------------------------------------

func TestArgMinMax(t *testing.T) {
	t.Run("min violation", func(t *testing.T) {
		var a Args
		a.StringList("names", "list of names", Min(2))

		res := make(ArgMap)
		_, err := a.parse([]string{"only_one"}, res)
		if err == nil {
			t.Fatal("expected error for min violation")
		}
	})

	t.Run("max violation", func(t *testing.T) {
		var a Args
		a.StringList("names", "list of names", Max(2))

		res := make(ArgMap)
		_, err := a.parse([]string{"a", "b", "c"}, res)
		if err == nil {
			t.Fatal("expected error for max violation")
		}
	})

	t.Run("within range", func(t *testing.T) {
		var a Args
		a.StringList("names", "list of names", Min(1), Max(3))

		res := make(ArgMap)
		_, err := a.parse([]string{"a", "b"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.StringList("names")
		if len(got) != 2 || got[0] != "a" || got[1] != "b" {
			t.Fatalf("expected [a b], got %v", got)
		}
	})

	t.Run("exactly at min", func(t *testing.T) {
		var a Args
		a.StringList("names", "list of names", Min(2), Max(4))

		res := make(ArgMap)
		_, err := a.parse([]string{"a", "b"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.StringList("names")
		if len(got) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(got))
		}
	})

	t.Run("exactly at max", func(t *testing.T) {
		var a Args
		a.StringList("names", "list of names", Min(1), Max(3))

		res := make(ArgMap)
		_, err := a.parse([]string{"a", "b", "c"}, res)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := res.StringList("names")
		if len(got) != 3 {
			t.Fatalf("expected 3 elements, got %d", len(got))
		}
	})
}

// ---------------------------------------------------------------------------
// TestArgRegistrationPanics
// ---------------------------------------------------------------------------

func TestArgRegistrationPanics(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		mustPanic(t, func() {
			var a Args
			a.String("", "some help")
		})
	})

	t.Run("empty help", func(t *testing.T) {
		mustPanic(t, func() {
			var a Args
			a.String("name", "")
		})
	})

	t.Run("duplicate name", func(t *testing.T) {
		mustPanic(t, func() {
			var a Args
			a.String("name", "first")
			a.String("name", "second")
		})
	})

	t.Run("arg after list", func(t *testing.T) {
		mustPanic(t, func() {
			var a Args
			a.StringList("names", "list of names")
			a.String("extra", "extra arg")
		})
	})

	t.Run("mandatory after optional", func(t *testing.T) {
		mustPanic(t, func() {
			var a Args
			a.String("opt", "optional", Default("x"))
			a.String("mand", "mandatory")
		})
	})
}

// ---------------------------------------------------------------------------
// TestArgMapAccessorPanics
// ---------------------------------------------------------------------------

func TestArgMapAccessorPanics(t *testing.T) {
	res := make(ArgMap)

	t.Run("String", func(t *testing.T) {
		mustPanic(t, func() { res.String("missing") })
	})
	t.Run("StringList", func(t *testing.T) {
		mustPanic(t, func() { res.StringList("missing") })
	})
	t.Run("Bool", func(t *testing.T) {
		mustPanic(t, func() { res.Bool("missing") })
	})
	t.Run("BoolList", func(t *testing.T) {
		mustPanic(t, func() { res.BoolList("missing") })
	})
	t.Run("Int", func(t *testing.T) {
		mustPanic(t, func() { res.Int("missing") })
	})
	t.Run("IntList", func(t *testing.T) {
		mustPanic(t, func() { res.IntList("missing") })
	})
	t.Run("Int64", func(t *testing.T) {
		mustPanic(t, func() { res.Int64("missing") })
	})
	t.Run("Int64List", func(t *testing.T) {
		mustPanic(t, func() { res.Int64List("missing") })
	})
	t.Run("Uint", func(t *testing.T) {
		mustPanic(t, func() { res.Uint("missing") })
	})
	t.Run("UintList", func(t *testing.T) {
		mustPanic(t, func() { res.UintList("missing") })
	})
	t.Run("Uint64", func(t *testing.T) {
		mustPanic(t, func() { res.Uint64("missing") })
	})
	t.Run("Uint64List", func(t *testing.T) {
		mustPanic(t, func() { res.Uint64List("missing") })
	})
	t.Run("Float64", func(t *testing.T) {
		mustPanic(t, func() { res.Float64("missing") })
	})
	t.Run("Float64List", func(t *testing.T) {
		mustPanic(t, func() { res.Float64List("missing") })
	})
	t.Run("Duration", func(t *testing.T) {
		mustPanic(t, func() { res.Duration("missing") })
	})
	t.Run("DurationList", func(t *testing.T) {
		mustPanic(t, func() { res.DurationList("missing") })
	})
}

// ---------------------------------------------------------------------------
// TestArgMapTypeMismatchPanics
// ---------------------------------------------------------------------------

func TestArgMapTypeMismatchPanics(t *testing.T) {
	// Store an int, then try to read it with every non-int accessor.
	res := ArgMap{
		"val": &ArgMapItem{Value: 42},
	}

	t.Run("String on int", func(t *testing.T) {
		mustPanic(t, func() { res.String("val") })
	})
	t.Run("StringList on int", func(t *testing.T) {
		mustPanic(t, func() { res.StringList("val") })
	})
	t.Run("Bool on int", func(t *testing.T) {
		mustPanic(t, func() { res.Bool("val") })
	})
	t.Run("BoolList on int", func(t *testing.T) {
		mustPanic(t, func() { res.BoolList("val") })
	})
	t.Run("IntList on int", func(t *testing.T) {
		mustPanic(t, func() { res.IntList("val") })
	})
	t.Run("Int64 on int", func(t *testing.T) {
		mustPanic(t, func() { res.Int64("val") })
	})
	t.Run("Int64List on int", func(t *testing.T) {
		mustPanic(t, func() { res.Int64List("val") })
	})
	t.Run("Uint on int", func(t *testing.T) {
		mustPanic(t, func() { res.Uint("val") })
	})
	t.Run("UintList on int", func(t *testing.T) {
		mustPanic(t, func() { res.UintList("val") })
	})
	t.Run("Uint64 on int", func(t *testing.T) {
		mustPanic(t, func() { res.Uint64("val") })
	})
	t.Run("Uint64List on int", func(t *testing.T) {
		mustPanic(t, func() { res.Uint64List("val") })
	})
	t.Run("Float64 on int", func(t *testing.T) {
		mustPanic(t, func() { res.Float64("val") })
	})
	t.Run("Float64List on int", func(t *testing.T) {
		mustPanic(t, func() { res.Float64List("val") })
	})
	t.Run("Duration on int", func(t *testing.T) {
		mustPanic(t, func() { res.Duration("val") })
	})
	t.Run("DurationList on int", func(t *testing.T) {
		mustPanic(t, func() { res.DurationList("val") })
	})
}

// ---------------------------------------------------------------------------
// TestArgMultipleParse
// ---------------------------------------------------------------------------

func TestArgMultipleParse(t *testing.T) {
	var a Args
	a.String("greeting", "a greeting")
	a.Int("count", "a count")
	a.StringList("names", "names")

	res := make(ArgMap)
	rest, err := a.parse([]string{"hello", "42", "a", "b"}, res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rest) != 0 {
		t.Fatalf("expected no remaining args, got %v", rest)
	}

	if res.String("greeting") != "hello" {
		t.Fatalf("expected 'hello', got '%s'", res.String("greeting"))
	}
	if res.Int("count") != 42 {
		t.Fatalf("expected 42, got %d", res.Int("count"))
	}
	names := res.StringList("names")
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("expected [a b], got %v", names)
	}
}
