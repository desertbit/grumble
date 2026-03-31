package grumble

import (
	"reflect"
	"testing"
	"time"
)

// helper: create a fresh FlagMap.
func newFlagMap() FlagMap {
	return make(FlagMap)
}

// helper: must not return an error.
func mustParse(t *testing.T, f *Flags, args []string) ([]string, FlagMap) {
	t.Helper()
	res := newFlagMap()
	left, err := f.parse(args, res)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return left, res
}

// helper: must return an error.
func mustFailParse(t *testing.T, f *Flags, args []string) error {
	t.Helper()
	res := newFlagMap()
	_, err := f.parse(args, res)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	return err
}

// helper: assert a function panics.
func assertPanics(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: expected panic, but did not panic", name)
		}
	}()
	fn()
}

// ---------------------------------------------------------------------------
// TestFlagRegistration – register every type and confirm no panic.
// ---------------------------------------------------------------------------

func TestFlagRegistration(t *testing.T) {
	f := &Flags{}

	// Each call must not panic.
	f.String("a", "str", "default", "help str")
	f.StringL("strl", "default", "help strl")
	f.StringList("b", "slist", []string{"x"}, "help slist")
	f.StringListL("slistl", []string{"y"}, "help slistl")
	f.Bool("c", "bflag", false, "help bool")
	f.BoolL("bflagl", true, "help booll")
	f.Int("d", "iflag", 0, "help int")
	f.IntL("iflagl", 1, "help intl")
	f.Int8("e", "i8", 0, "help int8")
	f.Int8L("i8l", 1, "help int8l")
	f.Int16("f", "i16", 0, "help int16")
	f.Int16L("i16l", 1, "help int16l")
	f.Int32("g", "i32", 0, "help int32")
	f.Int32L("i32l", 1, "help int32l")
	f.Int64("h", "i64", 0, "help int64")
	f.Int64L("i64l", 1, "help int64l")
	f.Uint("i", "uflag", 0, "help uint")
	f.UintL("uflagl", 1, "help uintl")
	f.Uint8("j", "u8", 0, "help uint8")
	f.Uint8L("u8l", 1, "help uint8l")
	f.Uint16("k", "u16", 0, "help uint16")
	f.Uint16L("u16l", 1, "help uint16l")
	f.Uint32("l", "u32", 0, "help uint32")
	f.Uint32L("u32l", 1, "help uint32l")
	f.Uint64("m", "u64", 0, "help uint64")
	f.Uint64L("u64l", 1, "help uint64l")
	f.Float32("n", "f32", 0.0, "help float32")
	f.Float32L("f32l", 1.0, "help float32l")
	f.Float64("o", "f64", 0.0, "help float64")
	f.Float64L("f64l", 1.0, "help float64l")
	f.Duration("p", "dur", time.Second, "help dur")
	f.DurationL("durl", time.Minute, "help durl")

	if f.empty() {
		t.Fatal("flags should not be empty after registration")
	}
}

// ---------------------------------------------------------------------------
// String flag tests
// ---------------------------------------------------------------------------

func TestFlagStringParse(t *testing.T) {
	t.Run("long flag", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		_, res := mustParse(t, f, []string{"--name", "hello"})
		if v := res.String("name"); v != "hello" {
			t.Fatalf("expected 'hello', got '%s'", v)
		}
		if res["name"].IsDefault {
			t.Fatal("IsDefault should be false when flag is provided")
		}
	})

	t.Run("default value", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		_, res := mustParse(t, f, []string{})
		if v := res.String("name"); v != "default" {
			t.Fatalf("expected 'default', got '%s'", v)
		}
		if !res["name"].IsDefault {
			t.Fatal("IsDefault should be true when flag is not provided")
		}
	})

	t.Run("short flag", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		_, res := mustParse(t, f, []string{"-s", "hello"})
		if v := res.String("name"); v != "hello" {
			t.Fatalf("expected 'hello', got '%s'", v)
		}
	})

	t.Run("equals form", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		_, res := mustParse(t, f, []string{"--name=hello"})
		if v := res.String("name"); v != "hello" {
			t.Fatalf("expected 'hello', got '%s'", v)
		}
	})

	t.Run("quoted string trims quotes", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		_, res := mustParse(t, f, []string{"--name", "\"hello world\""})
		if v := res.String("name"); v != "hello world" {
			t.Fatalf("expected 'hello world', got '%s'", v)
		}
	})

	t.Run("missing value", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		mustFailParse(t, f, []string{"--name"})
	})
}

func TestFlagStringLParse(t *testing.T) {
	f := &Flags{}
	f.StringL("name", "def", "help")
	_, res := mustParse(t, f, []string{"--name", "val"})
	if v := res.String("name"); v != "val" {
		t.Fatalf("expected 'val', got '%s'", v)
	}
}

// ---------------------------------------------------------------------------
// StringList flag tests
// ---------------------------------------------------------------------------

func TestFlagStringListParse(t *testing.T) {
	t.Run("default value", func(t *testing.T) {
		f := &Flags{}
		f.StringList("t", "tags", []string{"a", "b"}, "help")
		_, res := mustParse(t, f, []string{})
		got := res.StringList("tags")
		if !reflect.DeepEqual(got, []string{"a", "b"}) {
			t.Fatalf("expected [a b], got %v", got)
		}
	})

	t.Run("single value", func(t *testing.T) {
		f := &Flags{}
		f.StringList("t", "tags", []string{"a", "b"}, "help")
		_, res := mustParse(t, f, []string{"--tags", "foo"})
		got := res.StringList("tags")
		if !reflect.DeepEqual(got, []string{"foo"}) {
			t.Fatalf("expected [foo], got %v", got)
		}
	})

	t.Run("short flag", func(t *testing.T) {
		f := &Flags{}
		f.StringList("t", "tags", []string{"a"}, "help")
		_, res := mustParse(t, f, []string{"-t", "bar"})
		got := res.StringList("tags")
		if !reflect.DeepEqual(got, []string{"bar"}) {
			t.Fatalf("expected [bar], got %v", got)
		}
	})

	t.Run("equals form", func(t *testing.T) {
		f := &Flags{}
		f.StringList("t", "tags", []string{}, "help")
		_, res := mustParse(t, f, []string{"--tags=baz"})
		got := res.StringList("tags")
		if !reflect.DeepEqual(got, []string{"baz"}) {
			t.Fatalf("expected [baz], got %v", got)
		}
	})
}

func TestFlagStringListLParse(t *testing.T) {
	f := &Flags{}
	f.StringListL("items", []string{"x"}, "help")
	_, res := mustParse(t, f, []string{"--items", "y"})
	got := res.StringList("items")
	if !reflect.DeepEqual(got, []string{"y"}) {
		t.Fatalf("expected [y], got %v", got)
	}
}

// ---------------------------------------------------------------------------
// Bool flag tests
// ---------------------------------------------------------------------------

func TestFlagBoolParse(t *testing.T) {
	t.Run("long flag no value", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", false, "help")
		_, res := mustParse(t, f, []string{"--verbose"})
		if v := res.Bool("verbose"); !v {
			t.Fatal("expected true")
		}
	})

	t.Run("equals true", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", false, "help")
		_, res := mustParse(t, f, []string{"--verbose=true"})
		if v := res.Bool("verbose"); !v {
			t.Fatal("expected true")
		}
	})

	t.Run("equals false", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", false, "help")
		_, res := mustParse(t, f, []string{"--verbose=false"})
		if v := res.Bool("verbose"); v {
			t.Fatal("expected false")
		}
	})

	t.Run("value false as next arg", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", false, "help")
		_, res := mustParse(t, f, []string{"--verbose", "false"})
		if v := res.Bool("verbose"); v {
			t.Fatal("expected false")
		}
	})

	t.Run("default false", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", false, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Bool("verbose"); v {
			t.Fatal("expected false (default)")
		}
	})

	t.Run("default true", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", true, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Bool("verbose"); !v {
			t.Fatal("expected true (default)")
		}
	})

	t.Run("short flag", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", false, "help")
		_, res := mustParse(t, f, []string{"-b"})
		if v := res.Bool("verbose"); !v {
			t.Fatal("expected true via short flag")
		}
	})
}

func TestFlagBoolLParse(t *testing.T) {
	f := &Flags{}
	f.BoolL("debug", false, "help")
	_, res := mustParse(t, f, []string{"--debug"})
	if !res.Bool("debug") {
		t.Fatal("expected true")
	}
}

// ---------------------------------------------------------------------------
// Int flag tests
// ---------------------------------------------------------------------------

func TestFlagIntParse(t *testing.T) {
	t.Run("long flag", func(t *testing.T) {
		f := &Flags{}
		f.Int("i", "count", 42, "help")
		_, res := mustParse(t, f, []string{"--count", "10"})
		if v := res.Int("count"); v != 10 {
			t.Fatalf("expected 10, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Int("i", "count", 42, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Int("count"); v != 42 {
			t.Fatalf("expected 42, got %d", v)
		}
	})

	t.Run("short flag", func(t *testing.T) {
		f := &Flags{}
		f.Int("i", "count", 42, "help")
		_, res := mustParse(t, f, []string{"-i", "7"})
		if v := res.Int("count"); v != 7 {
			t.Fatalf("expected 7, got %d", v)
		}
	})

	t.Run("equals form", func(t *testing.T) {
		f := &Flags{}
		f.Int("i", "count", 42, "help")
		_, res := mustParse(t, f, []string{"--count=99"})
		if v := res.Int("count"); v != 99 {
			t.Fatalf("expected 99, got %d", v)
		}
	})

	t.Run("negative value", func(t *testing.T) {
		f := &Flags{}
		f.Int("i", "count", 0, "help")
		_, res := mustParse(t, f, []string{"--count=-5"})
		if v := res.Int("count"); v != -5 {
			t.Fatalf("expected -5, got %d", v)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		f := &Flags{}
		f.Int("i", "count", 42, "help")
		mustFailParse(t, f, []string{"--count", "abc"})
	})
}

func TestFlagIntLParse(t *testing.T) {
	f := &Flags{}
	f.IntL("level", 5, "help")
	_, res := mustParse(t, f, []string{"--level", "3"})
	if v := res.Int("level"); v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Int8 flag tests
// ---------------------------------------------------------------------------

func TestFlagInt8Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Int8("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "127"})
		if v := res.Int8("val"); v != 127 {
			t.Fatalf("expected 127, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Int8("", "val", 10, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Int8("val"); v != 10 {
			t.Fatalf("expected 10, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Int8("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "abc"})
	})

	t.Run("overflow", func(t *testing.T) {
		f := &Flags{}
		f.Int8("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "200"})
	})
}

func TestFlagInt8LParse(t *testing.T) {
	f := &Flags{}
	f.Int8L("small", 1, "help")
	_, res := mustParse(t, f, []string{"--small", "2"})
	if v := res.Int8("small"); v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Int16 flag tests
// ---------------------------------------------------------------------------

func TestFlagInt16Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Int16("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "1000"})
		if v := res.Int16("val"); v != 1000 {
			t.Fatalf("expected 1000, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Int16("", "val", 500, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Int16("val"); v != 500 {
			t.Fatalf("expected 500, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Int16("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "xyz"})
	})

	t.Run("overflow", func(t *testing.T) {
		f := &Flags{}
		f.Int16("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "40000"})
	})
}

func TestFlagInt16LParse(t *testing.T) {
	f := &Flags{}
	f.Int16L("mid", 100, "help")
	_, res := mustParse(t, f, []string{"--mid", "200"})
	if v := res.Int16("mid"); v != 200 {
		t.Fatalf("expected 200, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Int32 flag tests
// ---------------------------------------------------------------------------

func TestFlagInt32Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Int32("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "100000"})
		if v := res.Int32("val"); v != 100000 {
			t.Fatalf("expected 100000, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Int32("", "val", 77, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Int32("val"); v != 77 {
			t.Fatalf("expected 77, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Int32("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "notanum"})
	})

	t.Run("overflow", func(t *testing.T) {
		f := &Flags{}
		f.Int32("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "3000000000"})
	})
}

func TestFlagInt32LParse(t *testing.T) {
	f := &Flags{}
	f.Int32L("big", 11, "help")
	_, res := mustParse(t, f, []string{"--big", "22"})
	if v := res.Int32("big"); v != 22 {
		t.Fatalf("expected 22, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Int64 flag tests
// ---------------------------------------------------------------------------

func TestFlagInt64Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Int64("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "9999999999"})
		if v := res.Int64("val"); v != 9999999999 {
			t.Fatalf("expected 9999999999, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Int64("", "val", 64, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Int64("val"); v != 64 {
			t.Fatalf("expected 64, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Int64("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "bad"})
	})
}

func TestFlagInt64LParse(t *testing.T) {
	f := &Flags{}
	f.Int64L("huge", 1, "help")
	_, res := mustParse(t, f, []string{"--huge", "2"})
	if v := res.Int64("huge"); v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Uint flag tests
// ---------------------------------------------------------------------------

func TestFlagUintParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Uint("u", "num", 0, "help")
		_, res := mustParse(t, f, []string{"--num", "55"})
		if v := res.Uint("num"); v != 55 {
			t.Fatalf("expected 55, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Uint("u", "num", 12, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Uint("num"); v != 12 {
			t.Fatalf("expected 12, got %d", v)
		}
	})

	t.Run("short flag", func(t *testing.T) {
		f := &Flags{}
		f.Uint("u", "num", 0, "help")
		_, res := mustParse(t, f, []string{"-u", "8"})
		if v := res.Uint("num"); v != 8 {
			t.Fatalf("expected 8, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Uint("u", "num", 0, "help")
		mustFailParse(t, f, []string{"--num", "abc"})
	})

	t.Run("negative rejected", func(t *testing.T) {
		f := &Flags{}
		f.Uint("u", "num", 0, "help")
		mustFailParse(t, f, []string{"--num", "-1"})
	})
}

func TestFlagUintLParse(t *testing.T) {
	f := &Flags{}
	f.UintL("unum", 3, "help")
	_, res := mustParse(t, f, []string{"--unum", "4"})
	if v := res.Uint("unum"); v != 4 {
		t.Fatalf("expected 4, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Uint8 flag tests
// ---------------------------------------------------------------------------

func TestFlagUint8Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Uint8("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "255"})
		if v := res.Uint8("val"); v != 255 {
			t.Fatalf("expected 255, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Uint8("", "val", 8, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Uint8("val"); v != 8 {
			t.Fatalf("expected 8, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Uint8("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "nope"})
	})

	t.Run("overflow", func(t *testing.T) {
		f := &Flags{}
		f.Uint8("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "256"})
	})
}

func TestFlagUint8LParse(t *testing.T) {
	f := &Flags{}
	f.Uint8L("byte", 1, "help")
	_, res := mustParse(t, f, []string{"--byte", "2"})
	if v := res.Uint8("byte"); v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Uint16 flag tests
// ---------------------------------------------------------------------------

func TestFlagUint16Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Uint16("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "60000"})
		if v := res.Uint16("val"); v != 60000 {
			t.Fatalf("expected 60000, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Uint16("", "val", 16, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Uint16("val"); v != 16 {
			t.Fatalf("expected 16, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Uint16("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "bad"})
	})

	t.Run("overflow", func(t *testing.T) {
		f := &Flags{}
		f.Uint16("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "70000"})
	})
}

func TestFlagUint16LParse(t *testing.T) {
	f := &Flags{}
	f.Uint16L("port", 8080, "help")
	_, res := mustParse(t, f, []string{"--port", "9090"})
	if v := res.Uint16("port"); v != 9090 {
		t.Fatalf("expected 9090, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Uint32 flag tests
// ---------------------------------------------------------------------------

func TestFlagUint32Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Uint32("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "100000"})
		if v := res.Uint32("val"); v != 100000 {
			t.Fatalf("expected 100000, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Uint32("", "val", 32, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Uint32("val"); v != 32 {
			t.Fatalf("expected 32, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Uint32("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "nope"})
	})

	t.Run("overflow", func(t *testing.T) {
		f := &Flags{}
		f.Uint32("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "5000000000"})
	})
}

func TestFlagUint32LParse(t *testing.T) {
	f := &Flags{}
	f.Uint32L("id", 1, "help")
	_, res := mustParse(t, f, []string{"--id", "2"})
	if v := res.Uint32("id"); v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Uint64 flag tests
// ---------------------------------------------------------------------------

func TestFlagUint64Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Uint64("", "val", 0, "help")
		_, res := mustParse(t, f, []string{"--val", "18446744073709551615"})
		if v := res.Uint64("val"); v != 18446744073709551615 {
			t.Fatalf("expected max uint64, got %d", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Uint64("", "val", 64, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Uint64("val"); v != 64 {
			t.Fatalf("expected 64, got %d", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Uint64("", "val", 0, "help")
		mustFailParse(t, f, []string{"--val", "notanumber"})
	})
}

func TestFlagUint64LParse(t *testing.T) {
	f := &Flags{}
	f.Uint64L("bigid", 10, "help")
	_, res := mustParse(t, f, []string{"--bigid", "20"})
	if v := res.Uint64("bigid"); v != 20 {
		t.Fatalf("expected 20, got %d", v)
	}
}

// ---------------------------------------------------------------------------
// Float32 flag tests
// ---------------------------------------------------------------------------

func TestFlagFloat32Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Float32("", "ratio", 0.0, "help")
		_, res := mustParse(t, f, []string{"--ratio", "3.14"})
		if v := res.Float32("ratio"); v < 3.13 || v > 3.15 {
			t.Fatalf("expected ~3.14, got %f", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Float32("", "ratio", 1.5, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Float32("ratio"); v != 1.5 {
			t.Fatalf("expected 1.5, got %f", v)
		}
	})

	t.Run("equals form", func(t *testing.T) {
		f := &Flags{}
		f.Float32("", "ratio", 0, "help")
		_, res := mustParse(t, f, []string{"--ratio=2.5"})
		if v := res.Float32("ratio"); v != 2.5 {
			t.Fatalf("expected 2.5, got %f", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Float32("", "ratio", 0.0, "help")
		mustFailParse(t, f, []string{"--ratio", "abc"})
	})
}

func TestFlagFloat32LParse(t *testing.T) {
	f := &Flags{}
	f.Float32L("rate", 0.5, "help")
	_, res := mustParse(t, f, []string{"--rate", "0.75"})
	if v := res.Float32("rate"); v < 0.74 || v > 0.76 {
		t.Fatalf("expected ~0.75, got %f", v)
	}
}

// ---------------------------------------------------------------------------
// Float64 flag tests
// ---------------------------------------------------------------------------

func TestFlagFloat64Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Float64("", "precise", 0.0, "help")
		_, res := mustParse(t, f, []string{"--precise", "3.141592653589793"})
		if v := res.Float64("precise"); v != 3.141592653589793 {
			t.Fatalf("expected pi, got %f", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Float64("", "precise", 2.718, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Float64("precise"); v != 2.718 {
			t.Fatalf("expected 2.718, got %f", v)
		}
	})

	t.Run("negative", func(t *testing.T) {
		f := &Flags{}
		f.Float64("", "precise", 0, "help")
		_, res := mustParse(t, f, []string{"--precise=-1.5"})
		if v := res.Float64("precise"); v != -1.5 {
			t.Fatalf("expected -1.5, got %f", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Float64("", "precise", 0.0, "help")
		mustFailParse(t, f, []string{"--precise", "xyz"})
	})
}

func TestFlagFloat64LParse(t *testing.T) {
	f := &Flags{}
	f.Float64L("score", 0.0, "help")
	_, res := mustParse(t, f, []string{"--score", "99.9"})
	if v := res.Float64("score"); v != 99.9 {
		t.Fatalf("expected 99.9, got %f", v)
	}
}

// ---------------------------------------------------------------------------
// Duration flag tests
// ---------------------------------------------------------------------------

func TestFlagDurationParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f := &Flags{}
		f.Duration("d", "timeout", time.Second, "help")
		_, res := mustParse(t, f, []string{"--timeout", "5s"})
		if v := res.Duration("timeout"); v != 5*time.Second {
			t.Fatalf("expected 5s, got %v", v)
		}
	})

	t.Run("default", func(t *testing.T) {
		f := &Flags{}
		f.Duration("d", "timeout", time.Second, "help")
		_, res := mustParse(t, f, []string{})
		if v := res.Duration("timeout"); v != time.Second {
			t.Fatalf("expected 1s, got %v", v)
		}
	})

	t.Run("short flag", func(t *testing.T) {
		f := &Flags{}
		f.Duration("d", "timeout", time.Second, "help")
		_, res := mustParse(t, f, []string{"-d", "10m"})
		if v := res.Duration("timeout"); v != 10*time.Minute {
			t.Fatalf("expected 10m, got %v", v)
		}
	})

	t.Run("equals form", func(t *testing.T) {
		f := &Flags{}
		f.Duration("d", "timeout", time.Second, "help")
		_, res := mustParse(t, f, []string{"--timeout=2h"})
		if v := res.Duration("timeout"); v != 2*time.Hour {
			t.Fatalf("expected 2h, got %v", v)
		}
	})

	t.Run("complex duration", func(t *testing.T) {
		f := &Flags{}
		f.Duration("d", "timeout", 0, "help")
		_, res := mustParse(t, f, []string{"--timeout", "1h30m"})
		if v := res.Duration("timeout"); v != time.Hour+30*time.Minute {
			t.Fatalf("expected 1h30m, got %v", v)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		f := &Flags{}
		f.Duration("d", "timeout", time.Second, "help")
		mustFailParse(t, f, []string{"--timeout", "notaduration"})
	})
}

func TestFlagDurationLParse(t *testing.T) {
	f := &Flags{}
	f.DurationL("wait", time.Minute, "help")
	_, res := mustParse(t, f, []string{"--wait", "30s"})
	if v := res.Duration("wait"); v != 30*time.Second {
		t.Fatalf("expected 30s, got %v", v)
	}
}

// ---------------------------------------------------------------------------
// Leftover (positional) arguments
// ---------------------------------------------------------------------------

func TestFlagParseLeftover(t *testing.T) {
	t.Run("flags then positional", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		left, res := mustParse(t, f, []string{"--name", "hello", "arg1", "arg2"})
		if v := res.String("name"); v != "hello" {
			t.Fatalf("expected 'hello', got '%s'", v)
		}
		if !reflect.DeepEqual(left, []string{"arg1", "arg2"}) {
			t.Fatalf("expected [arg1 arg2], got %v", left)
		}
	})

	t.Run("positional only", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		left, res := mustParse(t, f, []string{"arg1", "arg2"})
		if v := res.String("name"); v != "default" {
			t.Fatalf("expected 'default', got '%s'", v)
		}
		if !reflect.DeepEqual(left, []string{"arg1", "arg2"}) {
			t.Fatalf("expected [arg1 arg2], got %v", left)
		}
	})

	t.Run("no args at all", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		left, _ := mustParse(t, f, []string{})
		if len(left) != 0 {
			t.Fatalf("expected no leftover, got %v", left)
		}
	})
}

// ---------------------------------------------------------------------------
// Double-dash stops flag parsing
// ---------------------------------------------------------------------------

func TestFlagParseDoubleDash(t *testing.T) {
	t.Run("double dash stops parsing", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		left, res := mustParse(t, f, []string{"--name", "hello", "--", "--other", "stuff"})
		if v := res.String("name"); v != "hello" {
			t.Fatalf("expected 'hello', got '%s'", v)
		}
		if !reflect.DeepEqual(left, []string{"--other", "stuff"}) {
			t.Fatalf("expected [--other stuff], got %v", left)
		}
	})

	t.Run("double dash as first arg", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		left, res := mustParse(t, f, []string{"--", "--name", "hello"})
		if v := res.String("name"); v != "default" {
			t.Fatalf("expected 'default', got '%s'", v)
		}
		if !reflect.DeepEqual(left, []string{"--name", "hello"}) {
			t.Fatalf("expected [--name hello], got %v", left)
		}
	})

	t.Run("double dash only", func(t *testing.T) {
		f := &Flags{}
		f.Bool("b", "verbose", false, "help")
		left, _ := mustParse(t, f, []string{"--"})
		if len(left) != 0 {
			t.Fatalf("expected no leftover, got %v", left)
		}
	})
}

// ---------------------------------------------------------------------------
// Invalid / unknown flags
// ---------------------------------------------------------------------------

func TestFlagParseInvalidFlag(t *testing.T) {
	t.Run("unknown long flag", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		mustFailParse(t, f, []string{"--unknown", "val"})
	})

	t.Run("unknown short flag", func(t *testing.T) {
		f := &Flags{}
		f.String("s", "name", "default", "help")
		mustFailParse(t, f, []string{"-x", "val"})
	})

	t.Run("no flags registered", func(t *testing.T) {
		f := &Flags{}
		mustFailParse(t, f, []string{"--anything"})
	})
}

// ---------------------------------------------------------------------------
// Registration panics
// ---------------------------------------------------------------------------

func TestFlagParseRegistrationPanics(t *testing.T) {
	t.Run("empty long", func(t *testing.T) {
		assertPanics(t, "empty long", func() {
			f := &Flags{}
			f.String("s", "", "default", "help")
		})
	})

	t.Run("long starting with dash", func(t *testing.T) {
		assertPanics(t, "long starts with -", func() {
			f := &Flags{}
			f.String("s", "-name", "default", "help")
		})
	})

	t.Run("long starting with double dash", func(t *testing.T) {
		assertPanics(t, "long starts with --", func() {
			f := &Flags{}
			f.String("s", "--name", "default", "help")
		})
	})

	t.Run("duplicate short", func(t *testing.T) {
		assertPanics(t, "duplicate short", func() {
			f := &Flags{}
			f.String("s", "name1", "default", "help")
			f.String("s", "name2", "default", "help")
		})
	})

	t.Run("duplicate long", func(t *testing.T) {
		assertPanics(t, "duplicate long", func() {
			f := &Flags{}
			f.String("a", "name", "default", "help")
			f.String("b", "name", "default", "help")
		})
	})

	t.Run("empty help", func(t *testing.T) {
		assertPanics(t, "empty help", func() {
			f := &Flags{}
			f.String("s", "name", "default", "")
		})
	})

	t.Run("short longer than 1 char", func(t *testing.T) {
		assertPanics(t, "short too long", func() {
			f := &Flags{}
			f.String("ab", "name", "default", "help")
		})
	})

	t.Run("short equals dash", func(t *testing.T) {
		assertPanics(t, "short is dash", func() {
			f := &Flags{}
			f.String("-", "name", "default", "help")
		})
	})
}

// ---------------------------------------------------------------------------
// Multiple flags combined
// ---------------------------------------------------------------------------

func TestFlagParseMultipleFlags(t *testing.T) {
	f := &Flags{}
	f.String("s", "name", "default", "help")
	f.Int("n", "count", 0, "help")
	f.Bool("v", "verbose", false, "help")
	f.Duration("d", "timeout", time.Second, "help")

	left, res := mustParse(t, f, []string{
		"--name", "test",
		"-n", "5",
		"--verbose",
		"--timeout=3s",
		"positional",
	})

	if v := res.String("name"); v != "test" {
		t.Fatalf("expected 'test', got '%s'", v)
	}
	if v := res.Int("count"); v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}
	if v := res.Bool("verbose"); !v {
		t.Fatal("expected true")
	}
	if v := res.Duration("timeout"); v != 3*time.Second {
		t.Fatalf("expected 3s, got %v", v)
	}
	if !reflect.DeepEqual(left, []string{"positional"}) {
		t.Fatalf("expected [positional], got %v", left)
	}
}

// ---------------------------------------------------------------------------
// Empty Flags struct
// ---------------------------------------------------------------------------

func TestFlagEmptyFlags(t *testing.T) {
	f := &Flags{}
	if !f.empty() {
		t.Fatal("expected empty")
	}
	f.String("s", "name", "default", "help")
	if f.empty() {
		t.Fatal("expected not empty after registration")
	}
}

// ---------------------------------------------------------------------------
// Sort
// ---------------------------------------------------------------------------

func TestFlagSort(t *testing.T) {
	f := &Flags{}
	f.StringL("zebra", "z", "help")
	f.StringL("apple", "a", "help")
	f.StringL("mango", "m", "help")

	f.sort()

	if f.list[0].Long != "apple" {
		t.Fatalf("expected 'apple' first, got '%s'", f.list[0].Long)
	}
	if f.list[1].Long != "mango" {
		t.Fatalf("expected 'mango' second, got '%s'", f.list[1].Long)
	}
	if f.list[2].Long != "zebra" {
		t.Fatalf("expected 'zebra' third, got '%s'", f.list[2].Long)
	}
}

// ---------------------------------------------------------------------------
// FlagMap accessor panics for missing flags
// ---------------------------------------------------------------------------

func TestFlagMapAccessorPanics(t *testing.T) {
	res := newFlagMap()

	t.Run("String panics", func(t *testing.T) {
		assertPanics(t, "String", func() { res.String("missing") })
	})
	t.Run("Bool panics", func(t *testing.T) {
		assertPanics(t, "Bool", func() { res.Bool("missing") })
	})
	t.Run("Int panics", func(t *testing.T) {
		assertPanics(t, "Int", func() { res.Int("missing") })
	})
	t.Run("Int8 panics", func(t *testing.T) {
		assertPanics(t, "Int8", func() { res.Int8("missing") })
	})
	t.Run("Int16 panics", func(t *testing.T) {
		assertPanics(t, "Int16", func() { res.Int16("missing") })
	})
	t.Run("Int32 panics", func(t *testing.T) {
		assertPanics(t, "Int32", func() { res.Int32("missing") })
	})
	t.Run("Int64 panics", func(t *testing.T) {
		assertPanics(t, "Int64", func() { res.Int64("missing") })
	})
	t.Run("Uint panics", func(t *testing.T) {
		assertPanics(t, "Uint", func() { res.Uint("missing") })
	})
	t.Run("Uint8 panics", func(t *testing.T) {
		assertPanics(t, "Uint8", func() { res.Uint8("missing") })
	})
	t.Run("Uint16 panics", func(t *testing.T) {
		assertPanics(t, "Uint16", func() { res.Uint16("missing") })
	})
	t.Run("Uint32 panics", func(t *testing.T) {
		assertPanics(t, "Uint32", func() { res.Uint32("missing") })
	})
	t.Run("Uint64 panics", func(t *testing.T) {
		assertPanics(t, "Uint64", func() { res.Uint64("missing") })
	})
	t.Run("Float32 panics", func(t *testing.T) {
		assertPanics(t, "Float32", func() { res.Float32("missing") })
	})
	t.Run("Float64 panics", func(t *testing.T) {
		assertPanics(t, "Float64", func() { res.Float64("missing") })
	})
	t.Run("Duration panics", func(t *testing.T) {
		assertPanics(t, "Duration", func() { res.Duration("missing") })
	})
	t.Run("StringList panics", func(t *testing.T) {
		assertPanics(t, "StringList", func() { res.StringList("missing") })
	})
}

// ---------------------------------------------------------------------------
// FlagMap type assertion panics (wrong type stored)
// ---------------------------------------------------------------------------

func TestFlagMapTypeMismatchPanics(t *testing.T) {
	res := newFlagMap()
	res["wrong"] = &FlagMapItem{Value: 42} // int stored

	t.Run("String on int panics", func(t *testing.T) {
		assertPanics(t, "String on int", func() { res.String("wrong") })
	})
	t.Run("Bool on int panics", func(t *testing.T) {
		assertPanics(t, "Bool on int", func() { res.Bool("wrong") })
	})
}

// ---------------------------------------------------------------------------
// IsDefault field check
// ---------------------------------------------------------------------------

func TestFlagIsDefault(t *testing.T) {
	f := &Flags{}
	f.String("s", "given", "def", "help")
	f.String("x", "notgiven", "def", "help")

	_, res := mustParse(t, f, []string{"--given", "val"})

	if res["given"].IsDefault {
		t.Fatal("expected IsDefault=false for provided flag")
	}
	if !res["notgiven"].IsDefault {
		t.Fatal("expected IsDefault=true for unprovided flag")
	}
}

// ---------------------------------------------------------------------------
// Duplicate empty short flags should NOT panic (both empty)
// ---------------------------------------------------------------------------

func TestFlagDuplicateEmptyShortNoPanic(t *testing.T) {
	f := &Flags{}
	f.StringL("first", "a", "help")
	f.StringL("second", "b", "help") // both have short="" — no conflict
	// If we got here without panic, the test passes.
}

// ---------------------------------------------------------------------------
// match helper tests
// ---------------------------------------------------------------------------

func TestFlagMatch(t *testing.T) {
	f := &Flags{}

	if !f.match("-s", "s", "name") {
		t.Fatal("expected match for short")
	}
	if !f.match("--name", "s", "name") {
		t.Fatal("expected match for long")
	}
	if f.match("--other", "s", "name") {
		t.Fatal("expected no match")
	}
	if f.match("-s", "", "name") {
		t.Fatal("expected no match for empty short")
	}
	if f.match("--name", "s", "") {
		t.Fatal("expected no match for empty long")
	}
}

// ---------------------------------------------------------------------------
// showDefault on flagItem
// ---------------------------------------------------------------------------

func TestFlagShowDefault(t *testing.T) {
	t.Run("bool does not show default", func(t *testing.T) {
		fi := &flagItem{Default: false}
		if fi.showDefault() {
			t.Fatal("bool flags should not show default")
		}
	})

	t.Run("string shows default", func(t *testing.T) {
		fi := &flagItem{Default: "hello"}
		if !fi.showDefault() {
			t.Fatal("string flags should show default")
		}
	})

	t.Run("int shows default", func(t *testing.T) {
		fi := &flagItem{Default: 42}
		if !fi.showDefault() {
			t.Fatal("int flags should show default")
		}
	})
}

// ---------------------------------------------------------------------------
// copyMissingValues on FlagMap
// ---------------------------------------------------------------------------

func TestFlagMapCopyMissingValues(t *testing.T) {
	t.Run("copies non-default values", func(t *testing.T) {
		dst := newFlagMap()
		src := newFlagMap()
		src["key"] = &FlagMapItem{Value: "val", IsDefault: false}

		dst.copyMissingValues(src, false)
		if dst["key"] == nil || dst["key"].Value != "val" {
			t.Fatal("expected key to be copied")
		}
	})

	t.Run("skips default when copyDefault is false", func(t *testing.T) {
		dst := newFlagMap()
		src := newFlagMap()
		src["key"] = &FlagMapItem{Value: "val", IsDefault: true}

		dst.copyMissingValues(src, false)
		if dst["key"] != nil {
			t.Fatal("expected key NOT to be copied (IsDefault + copyDefault=false)")
		}
	})

	t.Run("copies default when copyDefault is true", func(t *testing.T) {
		dst := newFlagMap()
		src := newFlagMap()
		src["key"] = &FlagMapItem{Value: "val", IsDefault: true}

		dst.copyMissingValues(src, true)
		if dst["key"] == nil || dst["key"].Value != "val" {
			t.Fatal("expected key to be copied (copyDefault=true)")
		}
	})

	t.Run("does not overwrite existing", func(t *testing.T) {
		dst := newFlagMap()
		dst["key"] = &FlagMapItem{Value: "existing"}
		src := newFlagMap()
		src["key"] = &FlagMapItem{Value: "new"}

		dst.copyMissingValues(src, true)
		if dst["key"].Value != "existing" {
			t.Fatal("expected existing value to be preserved")
		}
	})
}

// ---------------------------------------------------------------------------
// Bool flag: value that is NOT a bool literal after the flag is treated
// as next positional arg (because bool allows empty value).
// ---------------------------------------------------------------------------

func TestFlagBoolPositionalAfterFlag(t *testing.T) {
	f := &Flags{}
	f.Bool("v", "verbose", false, "help")

	// "--verbose" followed by something starting with "-" — should
	// NOT consume the next flag-like arg as a value for verbose.
	f2 := &Flags{}
	f2.Bool("v", "verbose", false, "help")
	f2.String("n", "name", "def", "help")

	_, res := mustParse(t, f2, []string{"--verbose", "--name", "hello"})
	if !res.Bool("verbose") {
		t.Fatal("expected verbose=true")
	}
	if v := res.String("name"); v != "hello" {
		t.Fatalf("expected 'hello', got '%s'", v)
	}
}

// ---------------------------------------------------------------------------
// Missing value for non-bool flag at end of args
// ---------------------------------------------------------------------------

func TestFlagParseMissingValueAtEnd(t *testing.T) {
	f := &Flags{}
	f.String("s", "name", "default", "help")
	mustFailParse(t, f, []string{"--name"})
}

// ---------------------------------------------------------------------------
// trimQuotes helper
// ---------------------------------------------------------------------------

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`hello`, "hello"},
		{`"`, `"`},
		{`""`, ""},
		{`"a`, `"a`},
		{`a"`, `a"`},
		{"", ""},
	}
	for _, tc := range tests {
		got := trimQuotes(tc.input)
		if got != tc.expected {
			t.Errorf("trimQuotes(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
