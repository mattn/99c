// Copyright 2017 The 99c Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command 99run executes binary programs produced by the 99c compiler.
//
// Usage
//
// To execute a compiled binary named a.out
//
//	99run a.out [arguments]
//
// Installation
//
// To install or update
//
//      $ go get [-u] github.com/cznic/99c/99run
//
// Online documentation: [godoc.org/github.com/cznic/99c/99run](http://godoc.org/github.com/cznic/99c/99run)
//
// Changelog
//
// 2017-01-07: Initial public release.
package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cznic/virtual"
)

func winshebang(fn string) bool {
	if runtime.GOOS == "windows" {
		ext := filepath.Ext(fn)
		if ext == ".bat" || ext == ".cmd" {
			return true
		}
	}
	return false
}

func skipWinshebang(br *bufio.Reader) {
	shebang, err := br.Peek(7)
	if err != nil {
		return
	}
	if string(shebang) != "@99run " {
		return
	}
	for {
		_, err = br.ReadBytes('\n')
		if err != nil {
			return
		}
		bb, err := br.Peek(1)
		if err != nil {
			return
		}
		if bb[0] != '@' {
			break
		}
	}
}

func exit(code int, msg string, arg ...interface{}) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, os.Args[0]+": "+msg, arg...)
	}
	os.Exit(code)
}

func main() {
	if len(os.Args) < 2 {
		exit(2, "invalid arguments %v\n", os.Args)
	}

	bin, err := os.Open(os.Args[1])
	if err != nil {
		exit(1, "%v\n", err)
	}

	var r io.Reader

	br := bufio.NewReader(bin)
	if winshebang(os.Args[1]) {
		skipWinshebang(br)
		r = base64.NewDecoder(base64.StdEncoding, br)
	} else {
		r = br
	}

	var b virtual.Binary
	if _, err := b.ReadFrom(r); err != nil {
		exit(1, "%v\n", err)
	}

	code, err := virtual.Exec(&b, os.Args[1:], os.Stdin, os.Stdout, os.Stderr, 0, 8<<20, "")
	if err != nil {
		if code == 0 {
			code = 1
		}
		exit(code, "%v\n", err)
	}

	exit(code, "")
}
