package main

import (
	"github.com/desertbit/grumble"
	"github.com/desertbit/readline"
	"time"
)

func main() {
	handleFunc := func(rl *readline.Instance) {
		app := grumble.NewWithReadline(&grumble.Config{
			Name:        "app",
			Description: "short app description",
			InterruptHandler: func(a *grumble.App, count int) {
				// do nothing
			},
			Flags: func(f *grumble.Flags) {
				f.String("d", "directory", "DEFAULT", "set an alternative directory path")
				f.Bool("v", "verbose", false, "enable verbose mode")
			},
		}, rl)

		app.AddCommand(&grumble.Command{
			Name:    "daemon",
			Help:    "run the daemon",
			Aliases: []string{"run"},

			Flags: func(f *grumble.Flags) {
				f.Duration("t", "timeout", time.Second, "timeout duration")
			},

			Args: func(a *grumble.Args) {
				a.String("service", "which service to start", grumble.Default("server"))
			},

			Run: func(c *grumble.Context) error {
				// Parent Flags.
				c.App.Println("directory:", c.Flags.String("directory"))
				c.App.Println("verbose:", c.Flags.Bool("verbose"))
				// Flags.
				c.App.Println("timeout:", c.Flags.Duration("timeout"))
				// Args.
				c.App.Println("service:", c.Args.String("service"))
				return nil
			},
		})

		adminCommand := &grumble.Command{
			Name:     "admin",
			Help:     "admin tools",
			LongHelp: "super administration tools",
		}
		app.AddCommand(adminCommand)

		app.Println("Welcome to remote shell, type <help> to get list of commands")
		app.Run()
	}

	cfg := &readline.Config{}
	readline.ListenRemote("tcp", ":5555", cfg, handleFunc)
}
