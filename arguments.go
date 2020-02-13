package grumble

import (
	"fmt"
)

// ExpectedArgs defines the expected arguments for a command.
type ExpectedArgs func(args []string) error

// NoArgs returns an error if any args are included.
func NoArgs(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown command: %q", args[0])
	}
	return nil
}

// MinimumNArgs returns an error if there is not at least N args.
func MinimumNArgs(n int) ExpectedArgs {
	return func(args []string) error {
		if len(args) < n {
			return fmt.Errorf("requires at least %d arg(s), only received %d", n, len(args))
		}
		return nil
	}
}

// MaximumNArgs returns an error if there are more than N args.
func MaximumNArgs(n int) ExpectedArgs {
	return func(args []string) error {
		if len(args) > n {
			return fmt.Errorf("accepts at most %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

// ExactArgs returns an error if there are not exactly N args.
func ExactArgs(n int) ExpectedArgs {
	return func(args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

// RangeArgs returns an error if the number of args is not within the expected range.
func RangeArgs(min int, max int) ExpectedArgs {
	return func(args []string) error {
		if len(args) < min || len(args) > max {
			return fmt.Errorf("accepts between %d and %d arg(s), received %d", min, max, len(args))
		}
		return nil
	}
}
