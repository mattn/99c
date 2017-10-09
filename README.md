# Table of Contents

1. [99c](#99c)
     1. Usage
     1. Installation
     1. Changelog
     1. Supported platforms and operating systems
     1. Project status
     1. Executing compiled programs
     1. Compiling a simple program
     1. Setting the output file name
     1. Obtaining the preprocessor output
     1. Multiple C files projects
     1. Using object files
     1. Stack traces
     1. Argument passing
     1. Executing a C program embedded in a Go program
     1. Calling into an embedded C library from Go
     1. Loading C plugins at run-time
     1. Inserting defines
     1. Specifying include paths
1. [99run](#99run)
     1. Usage
     1. Installation
     1. Changelog
1. [99trace](#99trace)
     1. Usage
     1. Installation
     1. Changelog
     1. Sample
1. [99strace](#99strace)
     1. Usage
     1. Installation
     1. Changelog
     1. Sample

# 99c

Command 99c is a c99 compiler targeting a virtual machine.

### Usage

Output of 99c -h

     Usage of 99c:
       -99lib
                Library link mode.
       -Dname
                Equivalent to inserting '#define name 1' at the start of the
                translation unit.
       -Dname=definition
                Equivalent to inserting '#define name definition' at the start
                of the translation unit.
       -E       Copy C-language source files to standard output, executing all
                preprocessor directives; no compilation shall be performed. If
                any operand is not a text file, the effects are unspecified.
       -Ipath
                Add path to the include files search paths.
       -Olevel
                Optimization setting, ignored.
       -Wwarn
                Warning level, ignored.
       -c       Suppress the link-edit phase of the compilation, and do not
                remove any object files that are produced.
       -g       Produce debug information, ignored.
       -o pathname
                Use the specified pathname, instead of the default a.out, for
                the executable file produced. If the -o option is present with
                -c or -E, the result is unspecified.

Rest of the input is a list of file names, either C (.c) files or object (.o) files.

### Installation

To install or update the compiler and the virtual machine

     $ go get [-u] github.com/cznic/99c github.com/cznic/99c/99run

Online documentation: [godoc.org/github.com/cznic/99c](http://godoc.org/github.com/cznic/99c)

### Changelog

2017-10-07: Initial public release.

### Supported platforms and operating systems

See: https://godoc.org/github.com/cznic/ccir#hdr-Supported_platforms_and_architectures

At the time of this writing, in GOOS_GOARCH form

     linux_386
     linux_amd64
     windows_386
     windows_amd64

Porting to other platforms/architectures is considered not difficult.

### Project status

Both the compiler and the C runtime library implementation are known to be
incomplete.  Missing pieces are added as needed. Please fill an issue when
you run into problems and be patient. Only limited resources can be
allocated to this project and to the related parts of the tool chain.

Also, contributions are welcome.

### Executing compiled programs

Running a binary on Linux

     $ ./a.out
     hello world
     $

Running a binary on Windows

     C:\> 99run a.out
     hello world
     C:\>

### Compiling a simple program

All in just a single C file.

     $ cd examples/hello/
     $ cat hello.c
     #include <stdio.h>

     int main() {
     	printf("hello world\n");
     }
     $ 99c hello.c && ./a.out
     hello world
     $

### Setting the output file name

If the output is a single file, use -o to set its name. (POSIX option)

     $ cd examples/hello/
     $ cat hello.c
     #include <stdio.h>

     int main() {
     	printf("hello world\n");
     }
     $ 99c -o hello hello.c && ./hello
     hello world
     $
     
### Obtaining the preprocessor output

Use -E to produce the cpp results. (POSIX option)

     $ cd examples/hello/
     /home/jnml/src/github.com/cznic/99c/examples/hello
     $ 99c -E hello.c
     # 24 "/home/jnml/src/github.com/cznic/ccir/libc/predefined.h"
     typedef   char   * __builtin_va_list ;
     # 41 "/home/jnml/src/github.com/cznic/ccir/libc/builtin.h"
     typedef   __builtin_va_list   __gnuc_va_list ;
     ...
     extern   void   flockfile ( FILE   *   __stream ) ;
     extern   int   ftrylockfile ( FILE   *   __stream ) ;
     extern   void   funlockfile ( FILE   *   __stream ) ;
     # 5 "hello.c"
     int   main ( )   {
     printf ( "hello world\n" ) ;
     }
     # 7 "/home/jnml/src/github.com/cznic/ccir/libc/crt0.c"
     int   main ( ) ;
     __FILE_TYPE__   __stdfiles [ 3 ] ;
     void   * stdin   =   & __stdfiles [ 0 ] ,   * stdout   =   & __stdfiles [ 1 ] ,   * stderr   =   & __stdfiles [ 2 ] ;
     void   _start ( int   argc ,   char   * * argv )
     {
     __register_stdfiles ( stdin ,   stdout ,   stderr ) ;
     __builtin_exit ( ( ( int   ( * ) ( ) ) main )   ( argc ,   argv ) ) ;
     }
     $

### Multiple C files projects

A translation unit may consist of multiple source files.

     $ cd examples/multifile/
     /home/jnml/src/github.com/cznic/99c/examples/multifile
     $ cat main.c
     char *hello();

     #include <stdio.h>
     int main() {
     	printf("%s\n", hello());
     }
     $ cat hello.c
     char *hello() {
     	return "hello world";
     }
     $ 99c main.c hello.c && ./a.out
     hello world
     $

### Using object files

Use -c to output object files. (POSIX option)

     $ cd examples/multifile/
     /home/jnml/src/github.com/cznic/99c/examples/multifile
     $ 99c -c hello.c main.c
     $ 99c hello.o main.o && ./a.out
     hello world
     $ 99c hello.o main.c && ./a.out
     hello world
     $ 99c hello.c main.o && ./a.out
     hello world
     $

### Stack traces

If the program source(s) are available at the same location(s) as when the
program was compiled, then any stack trace produced is annotated using the
source code lines.

     $ cd examples/stack/
     /home/jnml/src/github.com/cznic/99c/examples/stack
     $ cat stack.c
     void f(int n) {
     	if (n) {
     		f(n-1);
     		return;
     	}

     	*(char *)n;
     }

     int main() {
     	f(4);
     }
     $ 99c stack.c && ./a.out
     panic: runtime error: invalid memory address or nil pointer dereference [recovered]
     	panic: runtime error: invalid memory address or nil pointer dereference
     stack.c.f(0x0)
     	stack.c:7:1	0x0002c		load8          0x0	; -	// *(char *)n;
     stack.c.f(0x1)
     	stack.c:3:1	0x00028		call           0x21	; -	// f(n-1);
     stack.c.f(0x2)
     	stack.c:3:1	0x00028		call           0x21	; -	// f(n-1);
     stack.c.f(0x3)
     	stack.c:3:1	0x00028		call           0x21	; -	// f(n-1);
     stack.c.f(0x7f3400000004)
     	stack.c:3:1	0x00028		call           0x21	; -	// f(n-1);
     stack.c.main(0x7f3400000001, 0x7f34c9400030)
     	stack.c:11:1	0x0001d		call           0x21	; -	// f(4);
     /home/jnml/src/github.com/cznic/ccir/libc/crt0.c._start(0x1, 0x7f34c9400030)
     	/home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1	0x0000d		call           0x16	; -	// __builtin_exit(((int (*)())main) (argc, argv));

     [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x511e46]

     goroutine 1 [running]:
     github.com/cznic/virtual.(*cpu).run.func1(0xc42009e2a0)
     	/home/jnml/src/github.com/cznic/virtual/cpu.go:222 +0x26e
     panic(0x555340, 0x66a270)
     	/home/jnml/go/src/runtime/panic.go:491 +0x283
     github.com/cznic/virtual.readI8(...)
     	/home/jnml/src/github.com/cznic/virtual/cpu.go:74
     github.com/cznic/virtual.(*cpu).run(0xc42009e2a0, 0x2, 0x0, 0x0, 0x0)
     	/home/jnml/src/github.com/cznic/virtual/cpu.go:993 +0x2116
     github.com/cznic/virtual.New(0xc4201d2000, 0xc42000a090, 0x1, 0x1, 0x656540, 0xc42000c010, 0x656580, 0xc42000c018, 0x656580, 0xc42000c020, ...)
     	/home/jnml/src/github.com/cznic/virtual/virtual.go:73 +0x2be
     github.com/cznic/virtual.Exec(0xc4201d2000, 0xc42000a090, 0x1, 0x1, 0x656540, 0xc42000c010, 0x656580, 0xc42000c018, 0x656580, 0xc42000c020, ...)
     	/home/jnml/src/github.com/cznic/virtual/virtual.go:84 +0xe9
     main.main()
     	/home/jnml/src/github.com/cznic/99c/99run/main.go:37 +0x382
     $

### Argument passing

Command line arguments are passed the standard way.

     $ cd examples/args/
     /home/jnml/src/github.com/cznic/99c/examples/args
     $ cat args.c
     #include <stdio.h>

     int main(int argc, char **argv) {
     	for (int i = 0; i < argc; i++) {
     		printf("%i: %s\n", i, argv[i]);
     	}
     }
     $ 99c args.c && ./a.out foo bar -x -y - qux
     0: ./a.out
     1: foo
     2: bar
     3: -x
     4: -y
     5: -
     6: qux
     $

### Executing a C program embedded in a Go program

This example requires installation of additional tools

     $ go get -u github.com/cznic/assets github.com/cznic/httpfs
     $ cd examples/embedding/
     /home/jnml/src/github.com/cznic/99c/examples/embedding
     $ ls *
     main.c  main.go

     assets:
     keepdir
     $ cat main.c
     // +build ignore

     #include <stdio.h>

     int main() {
     	int c;
     	while ((c = getc(stdin)) != EOF) {
     		printf("%c", c >= 'a' && c <= 'z' ? c^' ' : c);
     	}
     }
     $ cat main.go
     //go:generate 99c -o assets/a.out main.c
     //go:generate assets

     package main

     import (
     	"bytes"
     	"fmt"
     	"strings"
     	"time"

     	"github.com/cznic/httpfs"
     	"github.com/cznic/virtual"
     )

     func main() {
     	fs := httpfs.NewFileSystem(assets, time.Now())
     	f, err := fs.Open("/a.out")
     	if err != nil {
     		panic(err)
     	}

     	var bin virtual.Binary
     	if _, err := bin.ReadFrom(f); err != nil {
     		panic(err)
     	}

     	var out bytes.Buffer
     	exitCode, err := virtual.Exec(&bin, nil, strings.NewReader("Foo Bar"), &out, &out, 0, 1<<20, "")
     	if err != nil {
     		panic(err)
     	}

     	fmt.Printf("%s\n%v\n", out.Bytes(), exitCode)
     }
     $ go generate && go build && ./embedding
     FOO BAR
     0
     $

### Calling into an embedded C library from Go

It's possible to call individual C functions from Go.

This example requires installation of additional tools

     $ go get -u github.com/cznic/assets github.com/cznic/httpfs
     $ cd examples/ffi/
     /home/jnml/src/github.com/cznic/99c/examples/ffi
     $ ls *
     lib42.c  main.go

     assets:
     keepdir
     $ cat lib42.c
     // +build ignore

     static int answer;

     int main() {
     	// Any library initialization comes here.
     	answer = 42;
     }

     // Use the -99lib option to prevent the linker from eliminating this function.
     int f42(int arg) {
     	return arg*answer;
     }
     $ cat main.go
     //go:generate 99c -99lib -o assets/a.out lib42.c
     //go:generate assets

     package main

     import (
     	"fmt"
     	"time"

     	"github.com/cznic/httpfs"
     	"github.com/cznic/ir"
     	"github.com/cznic/virtual"
     	"github.com/cznic/xc"
     )

     func main() {
     	fs := httpfs.NewFileSystem(assets, time.Now())
     	f, err := fs.Open("/a.out")
     	if err != nil {
     		panic(err)
     	}

     	var bin virtual.Binary
     	if _, err := bin.ReadFrom(f); err != nil {
     		panic(err)
     	}

     	m, _, err := virtual.New(&bin, nil, nil, nil, nil, 0, 1<<10, "")
     	if err != nil {
     		panic(err)
     	}

     	defer m.Close()

     	pc, ok := bin.Sym[ir.NameID(xc.Dict.SID("f42"))]
     	if !ok {
     		panic("symbol not found")
     	}

     	t, err := m.NewThread(1 << 10)
     	if err != nil {
     		panic(err)
     	}

     	for _, v := range []int{-1, 0, 1} {
     		var y int32
     		_, err := t.FFI1(pc, virtual.Int32Result{&y}, virtual.Int32(int32(v)))
     		if err != nil {
     			panic(err)
     		}

     		fmt.Println(y)
     	}
     }
     $ go generate && go build && ./ffi
     -42
     0
     42
     $

### Loading C plugins at run-time

It's possible to load C plugins at run-time.

     $ cd examples/plugin/
     /home/jnml/src/github.com/cznic/99c/examples/plugin
     $ ls *
     lib42.c  main.go
     $ cat lib42.c
     // +build ignore

     static int answer;

     int main() {
     	// Any library initialization comes here.
     	answer = 42;
     }

     // Use the -99lib option to prevent the linker from eliminating this function.
     int f42(int arg) {
     	return arg*answer;
     }
     $ cat main.go
     //go:generate 99c -99lib lib42.c

     package main

     import (
     	"fmt"
     	"os"

     	"github.com/cznic/ir"
     	"github.com/cznic/virtual"
     	"github.com/cznic/xc"
     )

     func main() {
     	f, err := os.Open("a.out")
     	if err != nil {
     		panic(err)
     	}

     	var bin virtual.Binary
     	if _, err := bin.ReadFrom(f); err != nil {
     		panic(err)
     	}

     	m, _, err := virtual.New(&bin, nil, nil, nil, nil, 0, 1<<10, "")
     	if err != nil {
     		panic(err)
     	}

     	defer m.Close()

     	pc, ok := bin.Sym[ir.NameID(xc.Dict.SID("f42"))]
     	if !ok {
     		panic("symbol not found")
     	}

     	t, err := m.NewThread(1 << 10)
     	if err != nil {
     		panic(err)
     	}

     	for _, v := range []int{1, 2, 3} {
     		var y int32
     		_, err := t.FFI1(pc, virtual.Int32Result{&y}, virtual.Int32(int32(v)))
     		if err != nil {
     			panic(err)
     		}

     		fmt.Println(y)
     	}
     }
     $ go generate && go run main.go
     42
     84
     126
     $

### Inserting defines

Use the -D flag to define additional macros on the command line.

     $ cd examples/define/
     /home/jnml/src/github.com/cznic/99c/examples/define
     $ ls *
     main.c
     $ cat main.c
     #include <stdio.h>

     int main() {
     #ifdef VERBOSE
     	printf(GREETING);
     #endif
     }
     $ 99c -DVERBOSE -DGREETING=\"hello\\n\" main.c && ./a.out
     hello
     $ 99c -DGREETING=\"hello\\n\" main.c && ./a.out
     $


### Specifying include paths

The -I flag defines additional include files search path(s).

     $ cd examples/include/
     /home/jnml/src/github.com/cznic/99c/examples/include
     $ ls *
     main.c
     
     foo:
     main.h
     $ cat main.c
     #include <stdio.h>
     #include "main.h"
     
     int main() {
     	printf(HELLO);
     }
     $ cat foo/main.h
     #ifndef _MAIN_H_
     #define _MAIN_H_
     
     #define HELLO "hello\n"
     
     #endif
     $ 99c main.c && ./a.out
     99c: main.c:2:10: include file not found: main.h (and 2 more errors)
     $ 99c -Ifoo main.c && ./a.out
     hello
     $

# 99run

Command 99run executes binary programs produced by the 99c compiler.

### Usage

    $ 99run a.out

On Linux a.out can be executed directly.

### Installation

To install or update 99run

     $ go get [-u] github.com/cznic/99c/99run

Online documentation: [godoc.org/github.com/cznic/99c/99run](http://godoc.org/github.com/cznic/99c/99run)

### Changelog

2017-10-07: Initial public release.

# 99trace

Command 99trace traces execution of binary programs produced by the 99c compiler.

The trace is written to stderr.

### Usage

    99trace a.out [arguments]

### Installation

To install or update 99trace

     $ go get [-u] -tags virtual.trace github.com/cznic/99c/99trace

Online documentation: [godoc.org/github.com/cznic/99c/99trace](http://godoc.org/github.com/cznic/99c/99trace)

### Changelog

2017-10-09: Initial public release.

### Sample

    $ cd examples/hello/
    /home/jnml/src/github.com/cznic/99c/examples/hello
    $ 99c hello.c && 99trace a.out 2>log
    hello world
    $ cat log
    # _start
    0x00002	func	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:13:1
    0x00003		arguments      		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00004		push64         (ds)	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00005		push64         (ds+0x10); /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00006		push64         (ds+0x20); /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00007		#register_stdfiles	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00008		arguments      		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x00009		sub            sp, 0x8	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000a		arguments      		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000b		push32         (ap-0x8)	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000c		push64         (ap-0x10); /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000d		call           0x16	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    # main
    0x00016	func	; hello.c:3:1
    0x00017		push           ap	; hello.c:3:1
    0x00018		zero32         		; hello.c:3:1
    0x00019		store32        		; hello.c:3:1
    0x0001a		arguments      		; hello.c:3:1
    0x0001b		push           ts+0x0	; hello.c:4:1
    0x0001c		#printf        		; hello.c:4:1
    0x0001d		return         		; hello.c:4:1
    
    0x0000e	#exit          	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    
    $ 

# 99strace

Command 99strace traces system calls of programs produced by the 99c compiler.

The trace is written to stderr.

### Usage

    99strace a.out [arguments]

### Installation

To install or update 99strace

     $ go get [-u] -tags virtual.strace github.com/cznic/99c/99strace

Online documentation: [godoc.org/github.com/cznic/99c/99strace](http://godoc.org/github.com/cznic/99c/99strace)

### Changelog

2017-10-09: Initial public release.

### Sample

    $ cd example/strace/
    /home/jnml/src/github.com/cznic/99c/example/strace
    $ ls *
    data.txt  main.c
    $ cat main.c 
    #include <stdlib.h>
    #include <fcntl.h>
    #include <unistd.h>
    
    #define BUFSIZE 1<<16
    
    int main(int argc, char **argv) {
    	char *buf = malloc(BUFSIZE);
    	if (!buf) {
    		return 1;
    	}
    
    	for (int i = 1; i < argc; i++) {
    		int fd = open(argv[i], O_RDWR);
    		if (fd < 0) {
    			return 1;
    		}
    
    		ssize_t n;
    		while ((n = read(fd, buf, BUFSIZE)) > 0) {
    			write(0, buf, n);
    		}
    	}
    	free(buf);
    }
    $ cat data.txt 
    Lorem ipsum.
    $ 99c main.c && ./a.out data.txt
    Lorem ipsum.
    $ 99strace a.out data.txt
    malloc(0x10000) 0x7f83a9400020
    open("data.txt", O_RDWR, 0) 5 errno 0
    read(5, 0x7f83a9400020, 65536) 13 errno 0
    Lorem ipsum.
    write(0, 0x7f83a9400020, 13) 13 errno 0
    read(5, 0x7f83a9400020, 65536) 0 errno 0
    freep(0x7f83a9400020)
    $ 
