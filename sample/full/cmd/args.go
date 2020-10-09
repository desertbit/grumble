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

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/desertbit/grumble"
)

func init() {
	App.AddCommand(&grumble.Command{
		Name: "args",
		Help: "test args",
		Args: func(a *grumble.Args) {
			a.String("s", "test string")
			a.Duration("d", "test duration", grumble.Default(time.Second))
			a.Int("i", "test int", grumble.Default(5))
			a.Int64("i64", "test int64", grumble.Default(int64(-88)))
			a.Uint("u", "test uint", grumble.Default(uint(66)))
			a.Uint64("u64", "test uint64", grumble.Default(uint64(8888)))
			a.Float64("f64", "test float64", grumble.Default(float64(5.889)))
			a.StringList("sl", "test string list", grumble.Default([]string{"first", "second", "third"}), grumble.Max(3))
		},
		Run: func(c *grumble.Context) error {
			fmt.Println("s  ", c.Args.String("s"))
			fmt.Println("d  ", c.Args.Duration("d"))
			fmt.Println("i  ", c.Args.Int("i"))
			fmt.Println("i64", c.Args.Int64("i64"))
			fmt.Println("u  ", c.Args.Uint("u"))
			fmt.Println("u64", c.Args.Uint64("u64"))
			fmt.Println("f64", c.Args.Float64("f64"))
			fmt.Println("sl ", strings.Join(c.Args.StringList("sl"), ","))
			return nil
		},
	})
}
