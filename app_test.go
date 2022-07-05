package grumble_test

import (
	"testing"

	"github.com/desertbit/grumble"
	r "github.com/stretchr/testify/require"
)

func TestApp_AddCommand(t *testing.T) {
	t.Parallel()

	a := grumble.New(&grumble.Config{
		Name:        "test application",
		Description: "test description",
	})

	cases := []struct {
		cmds  []*grumble.Command
		panic bool
	}{
		{cmds: []*grumble.Command{{}}, panic: true}, // 0
		{cmds: []*grumble.Command{{Name: "test"}}, panic: true},
		{cmds: []*grumble.Command{{Name: "test", Help: "test"}}, panic: false},
		{cmds: []*grumble.Command{{Name: "-test", Help: "test"}}, panic: true},
		{cmds: []*grumble.Command{{Name: "test", Help: "test"}, {}}, panic: true},
		{cmds: []*grumble.Command{{Name: "test", Help: "test"}, {Name: "hello"}}, panic: true}, // 5
		{cmds: []*grumble.Command{{Name: "test", Help: "test"}, {Name: "test", Help: "test2"}}, panic: false},
	}

	for i, c := range cases {
		// Define the test func.
		f := func() {
			//  Remove all previous commands first.
			a.Commands().RemoveAll()
			for i := range c.cmds {
				a.AddCommand(c.cmds[i])
			}
		}

		if c.panic {
			r.Panics(t, f, "case %d", i)
		} else {
			r.NotPanics(t, f, "case %d", i)
		}
	}
}

func TestApp_RunCommand(t *testing.T) {
	t.Parallel()

	a := grumble.New(&grumble.Config{
		Name:        "test application",
		Description: "test description",
		Flags: func(f *grumble.Flags) {
			f.Bool("b", "bool", false, "a test bool flag")
		},
	})

	cmd := &grumble.Command{
		Name:      "cmd",
		Aliases:   []string{"c", "cm"},
		Help:      "a test cmd help message",
		LongHelp:  "an even longer test cmd help message",
		HelpGroup: "testgroup",
		Usage:     "a test usage message",
	}
	a.AddCommand(cmd)

	cases := []struct {
		args []string
		err  bool
	}{
		// Empty args must return an error.
		{err: true},
		// Test both name and aliases.
		{args: []string{"c"}}, {args: []string{"cm"}}, {args: []string{"cmd"}},
	}

	for i, c := range cases {
		err := a.RunCommand(c.args)
		if c.err {
			r.Error(t, err, "case %d", i)
		} else {
			r.NoError(t, err, "case %d", i)
		}
	}
}
