/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2018 Roland Singer [roland.singer@deserbit.com]
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package grumble

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type parseArgFunc func(args []string, res ArgMap) ([]string, bool, error)
type defaultArgFunc func(res ArgMap)

type argItem struct {
	Name            string
	Help            string
	HelpArgs        string
	HelpShowDefault bool
	Default         interface{}

	isList   bool
	optional bool
}

// Args holds all the registered args.
type Args struct {
	parsers  []parseArgFunc
	defaults []defaultArgFunc
	list     []*argItem
}

// sort the args by their name.
func (a *Args) sort() {
	sort.Slice(a.list, func(i, j int) bool {
		return a.list[i].Name < a.list[j].Name
	})
}

func (a *Args) register(
	name, help, helpArgs string,
	helpShowDefault, isList, optional bool,
	defaultValue interface{},
	df defaultArgFunc,
	pf parseArgFunc,
) {
	// Validate.
	if name == "" {
		panic("empty argument name")
	} else if help == "" {
		panic(fmt.Errorf("missing help message for argument '%s'", name))
	}

	// Ensure the name is unique.
	for _, ai := range a.list {
		if ai.Name == name {
			panic(fmt.Errorf("argument '%s' registered twice", name))
		}
	}

	if !a.empty() {
		last := a.list[len(a.list)-1]
		// Check, if a list argument has been supplied already.
		if last.isList {
			panic("list argument has been registered, nothing can come after it")
		}

		// Check, that after an optional argument no mandatory one follows.
		if !optional && last.optional {
			panic("mandatory argument after optional")
		}
	}

	a.list = append(a.list, &argItem{
		Name:            name,
		Help:            help,
		HelpShowDefault: helpShowDefault,
		HelpArgs:        helpArgs,
		Default:         defaultValue,
		isList:          isList,
		optional:        optional,
	})
	a.defaults = append(a.defaults, df)
	a.parsers = append(a.parsers, pf)
}

// empty returns true, if the args are empty.
func (a *Args) empty() bool {
	return len(a.list) == 0
}

func (a *Args) parse(args []string, res ArgMap) ([]string, error) {
	// Iterate over all arguments that have been registered.
	// There must be either a default value or a value available,
	// otherwise the argument is missing.
	var (
		err    error
		parsed bool
	)
	for i, p := range a.parsers {
		args, parsed, err = p(args, res)
		if err != nil {
			return nil, err
		} else if !parsed {
			// Assign default value to it.
			if a.defaults[i] == nil {
				return nil, fmt.Errorf("missing argument '%s'", a.list[i].Name)
			}

			a.defaults[i](res)
		}
	}

	return args, nil
}

// String registers a string argument.
func (a *Args) String(name, help, defaultValue string, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "string", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}
			res[name] = &ArgMapItem{Value: args[0]}
			return args[1:], true, nil
		},
	)
}

// StringList registers a string list argument.
func (a *Args) StringList(name, help string, defaultValue []string, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "string list", true, true, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}
			res[name] = &ArgMapItem{Value: args}
			return []string{}, true, nil
		},
	)
}

// Bool registers a bool argument.
func (a *Args) Bool(name, help string, defaultValue, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "bool", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			b, err := strconv.ParseBool(args[0])
			if err != nil {
				return nil, false, fmt.Errorf("invalid bool value '%s' for argument: %s", args[0], name)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: b}
			return args[1:], true, nil
		},
	)
}

// BoolList registers a bool list argument.
func (a *Args) BoolList(name, help string, defaultValue []bool, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "bool list", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			var (
				err error
				bs  = make([]bool, len(args))
			)
			for i, a := range args {
				bs[i], err = strconv.ParseBool(a)
				if err != nil {
					return nil, false, fmt.Errorf("invalid bool value '%s' for argument: %s", a, name)
				}
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: bs}
			return []string{}, true, nil
		},
	)
}

// Int registers an int argument.
func (a *Args) Int(name, help string, defaultValue int, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "int", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			i, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, false, fmt.Errorf("invalid int value '%s' for argument: %s", args[0], name)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: i}
			return args[1:], true, nil
		},
	)
}

// IntList registers an int list argument.
func (a *Args) IntList(name, help string, defaultValue []int, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "int list", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			var (
				err error
				is  = make([]int, len(args))
			)
			for i, a := range args {
				is[i], err = strconv.Atoi(a)
				if err != nil {
					return nil, false, fmt.Errorf("invalid int value '%s' for argument: %s", a, name)
				}
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: is}
			return []string{}, true, nil
		},
	)
}

// Int64 registers an int64 argument.
func (a *Args) Int64(name, help string, defaultValue int64, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "int64", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			i, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return nil, false, fmt.Errorf("invalid int64 value '%s' for argument: %s", args[0], name)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: i}
			return args[1:], true, nil
		},
	)
}

// Int64List registers an int64 list argument.
func (a *Args) Int64List(name, help string, defaultValue []int64, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "int64 list", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			var (
				err error
				is  = make([]int64, len(args))
			)
			for i, a := range args {
				is[i], err = strconv.ParseInt(a, 10, 64)
				if err != nil {
					return nil, false, fmt.Errorf("invalid int64 value '%s' for argument: %s", a, name)
				}
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: is}
			return []string{}, true, nil
		},
	)
}

// Uint registers an uint argument.
func (a *Args) Uint(name, help string, defaultValue uint, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "uint", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			i, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return nil, false, fmt.Errorf("invalid uint value '%s' for argument: %s", args[0], name)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: i}
			return args[1:], true, nil
		},
	)
}

// UintList registers an uint list argument.
func (a *Args) UintList(name, help string, defaultValue []uint, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "uint list", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			var (
				err error
				u   uint64
				is  = make([]uint, len(args))
			)
			for i, a := range args {
				u, err = strconv.ParseUint(a, 10, 64)
				if err != nil {
					return nil, false, fmt.Errorf("invalid uint value '%s' for argument: %s", a, name)
				}
				is[i] = uint(u)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: is}
			return []string{}, true, nil
		},
	)
}

// Uint64 registers an uint64 argument.
func (a *Args) Uint64(name, help string, defaultValue uint64, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "uint64", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			i, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, false, fmt.Errorf("invalid uint64 value '%s' for argument: %s", args[0], name)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: i}
			return args[1:], true, nil
		},
	)
}

// Uint64List registers an uint64 list argument.
func (a *Args) Uint64List(name, help string, defaultValue []uint64, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "uint64 list", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			var (
				err error
				is  = make([]uint64, len(args))
			)
			for i, a := range args {
				is[i], err = strconv.ParseUint(a, 10, 64)
				if err != nil {
					return nil, false, fmt.Errorf("invalid uint64 value '%s' for argument: %s", a, name)
				}
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: is}
			return []string{}, true, nil
		},
	)
}

// Float64 registers a float64 argument.
func (a *Args) Float64(name, help string, defaultValue float64, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "float64", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			i, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return nil, false, fmt.Errorf("invalid float64 value '%s' for argument: %s", args[0], name)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: i}
			return args[1:], true, nil
		},
	)
}

// Float64List registers an float64 list argument.
func (a *Args) Float64List(name, help string, defaultValue []float64, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "float64 list", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			var (
				err error
				is  = make([]float64, len(args))
			)
			for i, a := range args {
				is[i], err = strconv.ParseFloat(a, 64)
				if err != nil {
					return nil, false, fmt.Errorf("invalid float64 value '%s' for argument: %s", a, name)
				}
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: is}
			return []string{}, true, nil
		},
	)
}

// Duration registers a duration argument.
func (a *Args) Duration(name, help string, defaultValue time.Duration, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "duration", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			if len(args) == 0 {
				return args, false, nil
			}

			i, err := time.ParseDuration(args[0])
			if err != nil {
				return nil, false, fmt.Errorf("invalid duration value '%s' for argument: %s", args[0], name)
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: i}
			return args[1:], true, nil
		},
	)
}

// DurationList registers an duration list argument.
func (a *Args) DurationList(name, help string, defaultValue []time.Duration, optional bool) {
	var df defaultArgFunc
	if optional {
		df = func(res ArgMap) { res[name] = &ArgMapItem{Value: defaultValue, IsDefault: true} }
	}

	a.register(name, help, "duration list", true, false, optional, defaultValue, df,
		func(args []string, res ArgMap) ([]string, bool, error) {
			var (
				err error
				is  = make([]time.Duration, len(args))
			)
			for i, a := range args {
				is[i], err = time.ParseDuration(a)
				if err != nil {
					return nil, false, fmt.Errorf("invalid duration value '%s' for argument: %s", a, name)
				}
			}

			// Set next argument as value.
			res[name] = &ArgMapItem{Value: is}
			return []string{}, true, nil
		},
	)
}
