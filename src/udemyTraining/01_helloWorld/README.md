
[Source](http://blog.hashbangbash.com/2014/04/linking-golang-statically/ "Permalink to Linking golang statically | blog.hashbangbash.com")

# Linking golang statically | blog.hashbangbash.com

If you are not familiar with [Golang][1], do take the [go tour][2] or read some of the [docs][3] first.

There are a number of reasons that folks are in love with golang. One the most mentioned is the static linking.

As long as the source being compiled is native go, the go compiler will statically link the executable. Though when you need to use [cgo][4], then the compiler has to use its external linker.

## Pure go

`

    // code-pure.go
    package main

    import "fmt"

    func main() {
            fmt.Println("hello, world!")
    }

`

Straight forward example. Let's compile it.

`

    $&gt; go build ./code-pure.go
    $&gt; ldd ./code-pure
            not a dynamic executable
    $&gt; file ./code-pure
    ./code-pure: ELF 64-bit LSB  executable, x86-64, version 1 (SYSV), statically linked, not stripped

`

## cgo

Using a contrived, but that passes through the C barrier:
`

    // code-cgo.go
    package main

    /*
    char* foo(void) { return "hello, world!"; }
    */
    import "C"

    import "fmt"

    func main() {
      fmt.Println(C.GoString(C.foo()))
    }

`

Seems simple enough. Let's compile it

`

    $&gt; go build ./code-cgo.go
    $&gt; file ./code-cgo
    ./code-cgo: ELF 64-bit LSB  executable, x86-64, version 1 (SYSV), dynamically linked (uses shared libs), not stripped
    $&gt; ldd ./code-cgo
            linux-vdso.so.1 (0x00007fff07339000)
            libpthread.so.0 =&gt; /lib64/libpthread.so.0 (0x00007f5e62737000)
            libc.so.6 =&gt; /lib64/libc.so.6 (0x00007f5e6236e000)
            /lib64/ld-linux-x86-64.so.2 (0x00007f5e62996000)
    $&gt; ./code-cgo
    hello, world!

`

## wait, what?

That code that is using cgo is _not_ statically linked. Why not?

The compile for this does not wholly use golang's internal linker, and has to use the external linker. So, this is not surprising, since this is not unlike simple ``gcc -o hello-world.c``, which is default to dynamically linked.

`

    // hello-world.c

    int main() {
            puts("hello, world!");
    }

`

    $&gt; gcc -o ./hello-world ./hello-world.c
    $&gt; ./hello-world
    hello, world!
    $&gt; file ./he
    hello-world    hello-world.c  hex.go
    $&gt; file ./hello-world
    ./hello-world: ELF 64-bit LSB  executable, x86-64, version 1 (SYSV), dynamically linked (uses shared libs), not stripped
    $&gt; ldd ./hello-world
            linux-vdso.so.1 (0x00007fff5f109000)
            libc.so.6 =&gt; /lib64/libc.so.6 (0x00007f0906e53000)
            /lib64/ld-linux-x86-64.so.2 (0x00007f090725e000)

`

But for that example, we just have to add the '-static' flag to gcc (and ensure that glibc-static package is available).

`

    $&gt; gcc -o ./hello-world -static ./hello-world.c
    $ file ./hello-world
    ./hello-world: ELF 64-bit LSB  executable, x86-64, version 1 (GNU/Linux), statically linked, not stripped
    $ ldd ./hello-world
            not a dynamic executable

`

Let's apply that same logic to our ``go build``

## static cgo

Using same `code-cgo.go` source, let's apply that gcc flag, but using the `go build` command.

`

    $&gt; go build --ldflags '-extldflags "-static"' ./code-cgo.go
    $&gt; file ./code-cgo
    ./code-cgo: ELF 64-bit LSB  executable, x86-64, version 1 (GNU/Linux), statically linked, not stripped
    $&gt; ldd ./code-cgo
            not a dynamic executable
    $&gt; ./code-cgo
    hello, world!

`

Cool! Here we've let the go compiler use the external linker, and that linker linked statically from libc.

An explanation of the flags here.

`\--ldflags` is passed to the go linker. It takes a string of arguments.
On my linux-x86_64 machine, that is the `6l` tool (`5l` for arm and `8l` for ix86). To see the tools available on your host, call `go tool`, and then get help on that tool with `go tool 6l --help`

`'-extldflags ...'` is a flag for the `6l` linker, to pass additional flags to the external linker (in my situation, that is `gcc`).

`"-static"` is the argument to `gcc` (also to `ld`) to link statically.

## gccgo love

Say you have a use case for using/needing the `gccgo` compiler, instead of the go compiler.
Again, using our same `code-cgo.go` source, let's compile the code statically using `gccgo`.

`

    $&gt; go build -compiler gccgo --gccgoflags "-static" ./code-cgo.go
    $&gt; file ./code-cgo
    ./code-cgo: ELF 64-bit LSB  executable, x86-64, version 1 (GNU/Linux), statically linked, not stripped
    $&gt; ldd ./code-cgo
            not a dynamic executable
    $&gt; ./code-cgo
    hello, world!

`

Huzzah! Still a static binary, using cgo _and_ giving gccgo a whirl. A quick run-down of these flags.

`-compiler gccgo` instructs the build to use `gccgo` instead of the go compiler (`gc`).

`\--gccgoflags ...` are additional arguments passed to the `gccgo` command. See also the `gccgo` man page.

`"-static"` similar to `gcc`, this is the same flag the instructs `gccgo` to link the executable statically.

## more info

If you're curious to what is happening behind the scenes of your compile, the `go build` has an `-x` flag that prints out the commands that it is running. This is often helpful if you are following what is linked in, or to see the commands used during compile such that you can find where and how to insert arguments needed.

`

[1]: http://golang.org/
[2]: http://tour.golang.org/
[3]: http://golang.org/doc/
[4]: http://golang.org/cmd/cgo/
