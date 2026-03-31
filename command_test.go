package grumble

import (
	"testing"
)

// ---------------------------------------------------------------------------
// TestCommandValidate
// ---------------------------------------------------------------------------

func TestCommandValidate(t *testing.T) {
	tests := []struct {
		name    string
		cmd     Command
		wantErr bool
	}{
		{
			name:    "empty name",
			cmd:     Command{Name: "", Help: "some help"},
			wantErr: true,
		},
		{
			name:    "name starts with dash",
			cmd:     Command{Name: "-bad", Help: "some help"},
			wantErr: true,
		},
		{
			name:    "empty help",
			cmd:     Command{Name: "good", Help: ""},
			wantErr: true,
		},
		{
			name:    "valid name and help",
			cmd:     Command{Name: "good", Help: "does things"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestCommandAddCommand
// ---------------------------------------------------------------------------

func TestCommandAddCommand(t *testing.T) {
	parent := &Command{Name: "parent", Help: "parent help"}

	child := &Command{Name: "child", Help: "child help"}
	parent.AddCommand(child)

	if child.Parent() != parent {
		t.Fatal("child.Parent() should equal parent")
	}

	// Add a second child and verify both are present.
	child2 := &Command{Name: "child2", Help: "child2 help"}
	parent.AddCommand(child2)

	if child2.Parent() != parent {
		t.Fatal("child2.Parent() should equal parent")
	}

	all := parent.commands.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 children, got %d", len(all))
	}
	if all[0].Name != "child" {
		t.Fatalf("expected first child named 'child', got %q", all[0].Name)
	}
	if all[1].Name != "child2" {
		t.Fatalf("expected second child named 'child2', got %q", all[1].Name)
	}
}

// ---------------------------------------------------------------------------
// TestCommandAddCommandPanics
// ---------------------------------------------------------------------------

func TestCommandAddCommandPanics(t *testing.T) {
	parent := &Command{Name: "parent", Help: "parent help"}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when adding command with empty name")
		}
	}()

	// Empty name should fail validation and cause a panic.
	parent.AddCommand(&Command{Name: "", Help: "bad"})
}

// ---------------------------------------------------------------------------
// TestCommandRegisterFlagsAndArgs
// ---------------------------------------------------------------------------

func TestCommandRegisterFlagsAndArgs(t *testing.T) {
	t.Run("with help flag", func(t *testing.T) {
		cmd := Command{
			Name: "test",
			Help: "test help",
			Flags: func(f *Flags) {
				f.String("v", "verbose", "no", "enable verbose")
			},
			Args: func(a *Args) {
				a.String("filename", "the file to process")
			},
		}
		cmd.registerFlagsAndArgs(true)

		// Parse flags to verify they were registered.
		flagRes := make(FlagMap)
		remaining, err := cmd.flags.parse([]string{"--verbose", "yes", "-h", "--", "myfile.txt"}, flagRes)
		if err != nil {
			t.Fatalf("flag parse error: %v", err)
		}

		if flagRes["verbose"] == nil {
			t.Fatal("expected 'verbose' flag to be registered")
		}
		if flagRes["help"] == nil {
			t.Fatal("expected 'help' flag to be registered when addHelpFlag=true")
		}

		// Parse args to verify they were registered.
		argRes := make(ArgMap)
		_, err = cmd.args.parse(remaining, argRes)
		if err != nil {
			t.Fatalf("arg parse error: %v", err)
		}
		if argRes["filename"] == nil {
			t.Fatal("expected 'filename' arg to be registered")
		}
		if argRes["filename"].Value != "myfile.txt" {
			t.Fatalf("expected filename='myfile.txt', got %v", argRes["filename"].Value)
		}
	})

	t.Run("without help flag", func(t *testing.T) {
		cmd := Command{
			Name: "test",
			Help: "test help",
			Flags: func(f *Flags) {
				f.Int("p", "port", 8080, "port number")
			},
		}
		cmd.registerFlagsAndArgs(false)

		// Only the user-defined flag should exist.
		flagRes := make(FlagMap)
		_, err := cmd.flags.parse([]string{}, flagRes)
		if err != nil {
			t.Fatalf("flag parse error: %v", err)
		}

		if flagRes["port"] == nil {
			t.Fatal("expected 'port' flag to be registered")
		}
		if _, ok := flagRes["help"]; ok {
			t.Fatal("did not expect 'help' flag when addHelpFlag=false")
		}
	})

	t.Run("nil flags and args callbacks", func(t *testing.T) {
		cmd := Command{
			Name: "bare",
			Help: "bare help",
		}
		cmd.registerFlagsAndArgs(true)

		// Should still have the help flag.
		flagRes := make(FlagMap)
		_, err := cmd.flags.parse([]string{}, flagRes)
		if err != nil {
			t.Fatalf("flag parse error: %v", err)
		}
		if flagRes["help"] == nil {
			t.Fatal("expected 'help' flag even when Flags callback is nil")
		}
	})
}

// ---------------------------------------------------------------------------
// TestCommandFlagsInheritance
// ---------------------------------------------------------------------------

func TestCommandFlagsInheritance(t *testing.T) {
	parent := &Command{
		Name: "parent",
		Help: "parent help",
		Flags: func(f *Flags) {
			f.Bool("d", "debug", false, "enable debug")
		},
	}
	parent.registerFlagsAndArgs(true)

	child := &Command{
		Name: "child",
		Help: "child help",
		Flags: func(f *Flags) {
			f.String("o", "output", "stdout", "output destination")
		},
	}
	parent.AddCommand(child)

	// Parse parent flags independently.
	parentFlags := make(FlagMap)
	_, err := parent.flags.parse([]string{"--debug"}, parentFlags)
	if err != nil {
		t.Fatalf("parent flag parse error: %v", err)
	}
	if parentFlags["debug"] == nil {
		t.Fatal("expected 'debug' flag on parent")
	}

	// Parse child flags independently.
	childFlags := make(FlagMap)
	_, err = child.flags.parse([]string{"--output", "/tmp/out"}, childFlags)
	if err != nil {
		t.Fatalf("child flag parse error: %v", err)
	}
	if childFlags["output"] == nil {
		t.Fatal("expected 'output' flag on child")
	}
	if childFlags["output"].Value != "/tmp/out" {
		t.Fatalf("expected output='/tmp/out', got %v", childFlags["output"].Value)
	}

	// Parent should not have child's flag.
	parentFlags2 := make(FlagMap)
	_, err = parent.flags.parse([]string{"--output", "x"}, parentFlags2)
	if err == nil {
		t.Fatal("expected error when parsing child flag on parent")
	}
}

// ---------------------------------------------------------------------------
// TestNewContext
// ---------------------------------------------------------------------------

func TestNewContext(t *testing.T) {
	cmd := &Command{Name: "ctx-cmd", Help: "ctx help"}
	flagMap := FlagMap{
		"verbose": &FlagMapItem{Value: true, IsDefault: false},
	}
	argMap := ArgMap{
		"file": &ArgMapItem{Value: "test.txt", IsDefault: false},
	}

	ctx := newContext(nil, cmd, flagMap, argMap)

	if ctx.App != nil {
		t.Fatal("expected App to be nil")
	}
	if ctx.Command != cmd {
		t.Fatal("expected Command to match")
	}
	if ctx.Flags == nil {
		t.Fatal("expected Flags to be non-nil")
	}
	if !ctx.Flags.Bool("verbose") {
		t.Fatal("expected verbose flag to be true")
	}
	if ctx.Args == nil {
		t.Fatal("expected Args to be non-nil")
	}
	if ctx.Args.String("file") != "test.txt" {
		t.Fatalf("expected file arg 'test.txt', got %q", ctx.Args.String("file"))
	}
}

// ---------------------------------------------------------------------------
// TestContextFields
// ---------------------------------------------------------------------------

func TestContextFields(t *testing.T) {
	// Verify that all Context struct fields are accessible and settable.
	cmd := &Command{Name: "my-cmd", Help: "my help"}
	ctx := &Context{
		App:     nil,
		Command: cmd,
		Flags:   make(FlagMap),
		Args:    make(ArgMap),
	}

	if ctx.Command.Name != "my-cmd" {
		t.Fatalf("expected command name 'my-cmd', got %q", ctx.Command.Name)
	}

	// Stop() requires a real App with a closer/readline, so we just verify
	// the method exists by confirming the Context type has it. Calling it
	// with a nil App would panic, which is expected behaviour.
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when calling Stop() with nil App")
		}
	}()
	ctx.Stop()
}

// ---------------------------------------------------------------------------
// TestParentNilByDefault
// ---------------------------------------------------------------------------

func TestParentNilByDefault(t *testing.T) {
	cmd := &Command{Name: "root", Help: "root help"}
	if cmd.Parent() != nil {
		t.Fatal("expected Parent() to be nil for a standalone command")
	}
}
