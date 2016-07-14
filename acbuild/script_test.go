// Copyright 2016 The appc Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"
)

func TestJoinLines(t *testing.T) {
	type testcase struct {
		input  []string
		output []string
	}
	cases := []testcase{
		testcase{
			[]string{"this is", "a test"},
			[]string{"this is", "a test"},
		},
		testcase{
			[]string{`this is \`, `a test`},
			[]string{"", "this is  a test"},
		},
		testcase{
			[]string{`this is\`, `another    \`, `test`},
			[]string{"", "", "this is another     test"},
		},
		testcase{
			[]string{`this is a test \`},
			[]string{`this is a test \`},
		},
		testcase{
			[]string{`this\`, `is\`, `a`, `test`},
			[]string{"", "", "this is a", "test"},
		},
	}
	for _, c := range cases {
		output := joinLines(c.input)
		if !equal(output, c.output) {
			t.Errorf("output, expected:%v actual:%v", c.output, output)
		}
	}
}

func TestTokenizeLine(t *testing.T) {
	type testcase struct {
		input  string
		output []string
		err    error
	}
	cases := []testcase{
		// nothing special here
		testcase{
			"",
			[]string{},
			nil,
		},
		testcase{
			"run -- test",
			[]string{"run", "--", "test"},
			nil,
		},

		// whitespace shenanigans
		testcase{
			"run --  test",
			[]string{"run", "--", "test"},
			nil,
		},
		testcase{
			"run -- 	 test",
			[]string{"run", "--", "test"},
			nil,
		},
		testcase{
			"   run --  test   ",
			[]string{"run", "--", "test"},
			nil,
		},
		testcase{
			"	run -- test",
			[]string{"run", "--", "test"},
			nil,
		},

		// comments
		testcase{
			"#this is a test",
			[]string{},
			nil,
		},
		testcase{
			"    #this is a test",
			[]string{},
			nil,
		},
		testcase{
			"run  -- #this is a test",
			[]string{"run", "--"},
			nil,
		},
		testcase{
			`this is a "comment in #a string"`,
			[]string{"this", "is", "a", "comment in #a string"},
			nil,
		},
		testcase{
			`this is a 'comment in #another string'`,
			[]string{"this", "is", "a", "comment in #another string"},
			nil,
		},

		// can it escape slashes?
		testcase{
			`set-name \\'hello'\\""`,
			[]string{"set-name", `\hello\`},
			nil,
		},

		// do single quotes behave as expected?
		testcase{
			`set-name \'hello\'`,
			[]string{"set-name", "'hello'"},
			nil,
		},
		testcase{
			`set-name 'hello'`,
			[]string{"set-name", "hello"},
			nil,
		},
		testcase{
			`set-name 'he llo'`,
			[]string{"set-name", "he llo"},
			nil,
		},
		testcase{
			`set-name \'he llo\'`,
			[]string{"set-name", "'he", "llo'"},
			nil,
		},

		// do double quotes behave as expected?
		testcase{
			`set-name \"hello\"`,
			[]string{"set-name", `"hello"`},
			nil,
		},
		testcase{
			`set-name "hello"`,
			[]string{"set-name", "hello"},
			nil,
		},
		testcase{
			`set-name "he llo"`,
			[]string{"set-name", "he llo"},
			nil,
		},
		testcase{
			`set-name \"he llo\"`,
			[]string{"set-name", `"he`, `llo"`},
			nil,
		},

		// can I properly nest quotes?
		testcase{
			`set-name "t'es't"`,
			[]string{"set-name", "t'es't"},
			nil,
		},
		testcase{
			`set-name "t'est"`,
			[]string{"set-name", "t'est"},
			nil,
		},
		testcase{
			`set-name 't"es"t'`,
			[]string{"set-name", `t"es"t`},
			nil,
		},
		testcase{
			`set-name 'tes"t'`,
			[]string{"set-name", `tes"t`},
			nil,
		},

		// will it error when it should?
		testcase{
			`set-name 'test"`,
			[]string{},
			errSingleQuote,
		},
		testcase{
			`set-name "test'`,
			[]string{},
			errDoubleQuote,
		},
		testcase{
			`set-name test\`,
			[]string{},
			errEscape,
		},

		// are quotes in the middle of words handled?
		testcase{
			`set-name t"es"t`,
			[]string{"set-name", "test"},
			nil,
		},
		testcase{
			`set-name t'es't`,
			[]string{"set-name", "test"},
			nil,
		},
		testcase{
			`set-name "te"'st'`,
			[]string{"set-name", "test"},
			nil,
		},

		// let's try something crazy
		testcase{
			`set-name 't'h"i's "i\'s" 'a' "t''''e""""s'""'t\\`,
			[]string{"set-name", `thi's i's 'a' tes""t\`},
			nil,
		},
	}

	for _, c := range cases {
		output, err := tokenizeLine(c.input)
		if err != c.err {
			t.Errorf("error, expected:%v actual:%v", c.err, err)
		}
		if !equal(output, c.output) {
			t.Errorf("output, expected:%v actual:%v", c.output, output)
		}
	}
}

func TestInsertRunTacks(t *testing.T) {
	type testcase struct {
		input  []string
		output []string
	}
	cases := []testcase{
		testcase{
			[]string{"run", "--", "test"},
			[]string{"run", "--", "test"},
		},
		testcase{
			[]string{"run", "test"},
			[]string{"run", "--", "test"},
		},
		testcase{
			[]string{"run", "test", "--foo"},
			[]string{"run", "--", "test", "--foo"},
		},
		testcase{
			[]string{"run", "--foo", "test"},
			[]string{"run", "--foo", "--", "test"},
		},
		testcase{
			[]string{"run", "--foo", "test", "--bar"},
			[]string{"run", "--foo", "--", "test", "--bar"},
		},
		testcase{
			[]string{"run", "--foo", "--", "test", "--bar"},
			[]string{"run", "--foo", "--", "test", "--bar"},
		},
	}
	for _, c := range cases {
		output := insertRunTacks(c.input)
		if !equal(output, c.output) {
			t.Errorf("output, expected:%v actual:%v", c.output, output)
		}
	}
}

// no really guys, this language is _great_
func equal(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
