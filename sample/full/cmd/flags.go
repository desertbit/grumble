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
	"time"

	"github.com/desertbit/grumble"
)

func init() {
	App.AddCommand(&grumble.Command{
		Name: "flags",
		Help: "test flags",
		Flags: func(f *grumble.Flags) {
			f.Bool("b", "bool", false, "test bool")
			f.Int("i", "int", 1, "test int")
			f.Int8L("int8", -8, "test int8")
			f.Int16L("int16", -16, "test int16")
			f.Int32L("int32", -32, "test int32")
			f.Int64L("int64", -64, "test int64")
			f.Uint("u", "uint", 3, "test uint")
			f.Uint8L("uint8", 8, "test uint8")
			f.Uint16L("uint16", 16, "test uint16")
			f.Uint32L("uint32", 32, "test uint32")
			f.Uint64L("uint64", 64, "test uint64")
			f.Float32L("float32", 5.55, "test float32")
			f.Float64("f", "float64", 5.55, "test float64")
			f.Duration("d", "duration", time.Second, "duration test")
		},
		Run: func(c *grumble.Context) error {
			fmt.Println("bool     ", c.Flags.Bool("bool"))
			fmt.Println("int      ", c.Flags.Int("int"))
			fmt.Println("int8     ", c.Flags.Int8("int8"))
			fmt.Println("int16    ", c.Flags.Int16("int16"))
			fmt.Println("int32    ", c.Flags.Int32("int32"))
			fmt.Println("int64    ", c.Flags.Int64("int64"))
			fmt.Println("uint     ", c.Flags.Uint("uint"))
			fmt.Println("uint8    ", c.Flags.Uint8("uint8"))
			fmt.Println("uint16   ", c.Flags.Uint16("uint16"))
			fmt.Println("uint32   ", c.Flags.Uint32("uint32"))
			fmt.Println("uint64   ", c.Flags.Uint64("uint64"))
			fmt.Println("float32  ", c.Flags.Float32("float32"))
			fmt.Println("float64  ", c.Flags.Float64("float64"))
			fmt.Println("duration ", c.Flags.Duration("duration"))
			return nil
		},
	})
}
