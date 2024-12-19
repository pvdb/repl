# repl(1) -- sometimes you _really_ need a REPL

`repl` wraps a non-interactive command in an interactive REPL _(read-eval-print loop)_.

## Synopsis

    repl [options] command <...>

`command` is executed using the lines you type into `repl`'s prompt as either command-line arguments or else as standard input, and anything written by `command` to standard output and standard error is displayed.

This is repeated until you exit out of `repl`'s interactive loop by using either `CTRL-C` or `CTRL-D` at the prompt (default is `>> `):

    $ repl git ls-files
    >> exe
    exe/repl
    
    >> lib
    lib/repl.rb
    lib/repl/version.rb
    
    >> bin
    bin/console
    bin/setup
    
    >> ^D
    $ _

It can be used with commands you wish had an interactive mode, but don't, like `gem`:

    $ repl gem
    >> --version
    3.0.3
    
    >> sources
    *** CURRENT SOURCES ***
    
    https://rubygems.org/
    
    >> search '^repl$'
    
    *** REMOTE GEMS ***
    
    repl (1.0.0)
    
    >> ^D
    $ _

Or system utilities like `host`:

    $ repl host -t A
    >> google.com
    google.com has address 216.58.212.206
    
    >> google.co.uk
    google.co.uk has address 172.217.169.35
    
    >> ^D
    $ _

## Usage

`repl` is meant to wrap programs which accept and process command line arguments or which read from and process standard input, and which in turn print to standard output or standard error, but which don't have an interactive REPL of their own.

If you have [`rlwrap(1)`][rlwrap] installed you'll automatically get the full benefits of [`readline`][readline]: persistent history, reverse history searches, command-specific tab completions, rich command-line editing, etc.

Combined with `rlwrap`, `repl` can provide a much richer interactive environment even when wrapping commands that have their own basic one, by supporting `readline` features the command itself doesn't.

## Limitation

`repl` keeps no state between subsequent `command` invocations and, as such, cannot be used to replace things like the Ruby and Python REPLs _(`irb`, `pry`, etc.)_, or other language shells.

## Credits

This version borrows _very, very_ heavily from but is a ground-up rewrite of [Chris Wanstrath][defunkt]'s awesome [original version][defunkt/repl].

[defunkt/repl]: https://github.com/defunkt/repl
[pvdb/repl]: https://github.com/pvdb/repl
[defunkt]: https://github.com/defunkt
[rlwrap]: https://github.com/hanslub42/rlwrap
[readline]: https://tiswww.case.edu/php/chet/readline/rltop.html

## Installation

`repl` is easily installed as a standalone utility somewhere in your `${PATH}`; the default install location is `/usr/local/bin`, but it can be anything on your `${PATH}`.

This way you can run `repl` without any changes to your system's `$PATH`, but you can adjust `REPL_INSTALL_DIR` to match your shell environment.

It comes in two flavors: a Ruby version and a Go version. The Ruby version is the original one, while the Go version is a complete rewrite of the Ruby version, but with the same functionality.

### ![Ruby logo](ruby_small.png "Ruby logo") Ruby version

The Ruby version of the `repl` utility requires [Ruby](https://www.ruby-lang.org/en/) to be installed on your system, but has no other external dependencies _(that is: it only uses classes and modules from [Ruby's standard library](https://docs.ruby-lang.org/en/master/standard_library_md.html))_.

    REPL_INSTALL_DIR=/usr/local/bin
    curl -sO https://raw.githubusercontent.com/pvdb/repl/master/repl.rb
    mv repl.rb "${REPL_INSTALL_DIR}/repl"
    chmod 755 "${REPL_INSTALL_DIR}/repl"

### ![Go logo](go_small.png "Go logo") Go version

The Go version of the `repl` utility requires [Go](https://go.dev/) to be installed on your system, but has no other external dependencies _(that is: it only uses packages from [Go's standard library](https://pkg.go.dev/std))_.

    REPL_INSTALL_DIR=/usr/local/bin
    curl -s -O https://raw.githubusercontent.com/pvdb/repl/main/repl.go
    go build -o "${REPL_INSTALL_DIR}/repl" repl.go
    rm repl.go
    chmod 755 "${REPL_INSTALL_DIR}/repl"

## Documentation

The `repl` documentation is maintained in the `repl.1.ronn` Markdown file, which is converted into a man page, including an HTML version, as follows:

    brew install groff
    gem install ronn-ng

    ronn --roff --html --date="$(date -Idate)" --organization PVDB --manual="Awesome Utilities" repl.1.ronn

## Examples

In order to show what's going on, the following examples use the `--debug` option, which makes `repl` print the expanded command line - prefixed with `$ ` - just before executing in.

### basic

By default, anything you enter on `repl`'s prompt is passed on to the command as positional command line arguments:

    $ repl --debug git
    >> version
    $ git version
    git version 2.28.0
    
    >> status
    $ git status
    fatal: not a git repository (or any of the parent directories): .git
    
    >> ^D
    $ _

You can also enter multiple command line arguments at once:

    $ repl --debug git config --global
    >> --get user.name
    $ git config --global --get user.name
    Peter Vandenberk
    
    >> --get user.email
    $ git config --global --get user.email
    pvandenberk@mac.com
    
    >> ^D
    $ _

### placeholders

You can control where `repl` inserts what you enter on the prompt by using a `%s` placeholder, similar to `printf(3)`:

    $ repl --debug grep %s ~/.gitconfig
    >> name
    $ grep name ~/.gitconfig
            name = Peter Vandenberk
    
    >> email
    $ grep email ~/.gitconfig
            email = pvandenberk@mac.com
    
    >> ^D
    $ _

You can also use a `%s` placeholder for the command itself, as opposed to its arguments:

    $ repl --debug %s /tmp
    >> file
    $ file /tmp
    /tmp: sticky, directory
    
    >> ls -ld
    $ ls -ld /tmp
    lrwxr-xr-x@ 1 root  admin  11 15 Jun 20:19 /tmp -> private/tmp
    
    >> ^D
    $ _

Multiple `%s` placeholders are also supported:

    $ repl --debug file /etc/%s /etc/%s_Apple_Terminal
    >> bashrc
    $ file /etc/bashrc /etc/bashrc_Apple_Terminal
    /etc/bashrc:                ASCII text
    /etc/bashrc_Apple_Terminal: ASCII text
    
    >> zshrc
    $ file /etc/zshrc /etc/zshrc_Apple_Terminal
    /etc/zshrc:                ASCII text
    /etc/zshrc_Apple_Terminal: ASCII text
    
    >> ^D
    $ _

### standard input

Using the `--stdin` option tells `repl` to write what you enter on the prompt to the command's standard input _(instead of providing it as command arguments)_:

    $ repl --debug --stdin bc
    >> 21 * 2
    $ /bin/echo "21 * 2" | bc
    42
    
    >> 14 * 3
    $ /bin/echo "14 * 3" | bc
    42
    
    >> ^D
    $ _

Use the `--printf` option _(in addition to `--stdin`)_ to make `repl` use `printf(1)` instead of `echo(1)`, thus suppressing newline characters if and when these are superfluous:

    $ repl --debug --stdin --printf wc -c
    >> one
    $ /usr/bin/printf "one" | wc -c
           3
    
    >> two
    $ /usr/bin/printf "two" | wc -c
           3
    
    >> three
    $ /usr/bin/printf "three" | wc -c
           5
    
    >> ^D
    $ _

### pipelines

By single-quoting the command, you can wrap `repl` around entire shell pipelines:

    $ repl --debug '%s blegga|wc -c'
    >> echo
    $ echo blegga|wc -c
           7
    
    >> printf
    $ printf blegga|wc -c
           6
    
    >> ^D
    $ _

### escaping

In order to create these pipelines, `repl` doesn't by default escape shell constructs, which can cause issues:

    $ repl --debug echo
    >> foo|bar
    $ echo foo|bar
    sh: bar: command not found
    
    >> <blegga>
    $ echo <blegga>
    sh: -c: line 0: syntax error near unexpected token `newline'
    sh: -c: line 0: `echo <blegga>'
    
    >> ^D
    $ _

By using the `--escape` option, `repl` can be made to escape shell constructs in its input:

    $ repl --debug --escape echo
    >> foo|bar
    $ echo foo\|bar
    foo|bar
    
    >> <blegga>
    $ echo \<blegga\>
    <blegga>
    
    >> ^D
    $ _

## Features

Features of and improvements in this new version of `repl`:

### new options

1. `--version` - print the `repl` version info
   * example: `repl --version`
1. `--html` - open the `repl` man page in web browser
   * example: `repl --html`
1. `--escape` - shell escape user's input
   * example: compare `repl echo <<<'Peter "The Rock" V.'` with `repl --escape echo <<<'Peter "The Rock" V.'`
1. `--printf` - causes `--stdin` to use `/usr/bin/printf` instead of `/bin/echo` to avoid adding superfluous trailing newline characters
   * example: compare `repl --debug --stdin wc -c <<<"blegga"` with `repl --debug --stdin --printf wc -c <<<"blegga"`

### improved options
1. `--debug` - now works correctly in conjunction with `--stdin`
   * example: `repl --stdin --debug wc -c` using [defunkt/repl][] doesn't work

### `rlwrap` improvements

1. explicitly ignore options/flags for calculating history and completion files
1. use `MakeMakefile.find_executable0()` instead of `which(1)` to find `rlwrap`
1. set `__REPL_RLWRAP__` to the PID of the parent `repl` process, instead of `0`

### command processing

1. all prompts, debug and error messages are written to `IO.console` (not `STDOUT`) for pipelining purposes
   * example: `repl --debug echo > blegga; cat blegga`
1. support for multiple embedded `%s` placeholders in the command string, not just one
   * example: `repl diff source_dir/%s target_dir/%s`

## Configuration

The following environment variables can be used to override `repl`'s defaults:

* `REPL_PROMPT`: prompt to use in `repl`'s read-eval-print loop (default: `>>`)
* `REPL_DEBUG` : equivalent to the `--debug` option if set to true (default: false)
* `REPL_QUIET` : equivalent to the `--quiet` option if set to true (default: false)
* `REPL_HISTORY_DIR`: directory in which history files are kept (default: `${HOME}`)
* `REPL_COMPLETION_DIR`: directory in which completion files are kept (default: `${HOME}/.repl`)

These options can also be set permantently in `${HOME}/.repl.conf`, instead of "polluting" your shell environment with them; use the [`repl.conf`](repl.conf) template file as a starting point.

## TODO

Potential improvements to this new version:

1. update and improve the usage instructions and examples
1. update and improve the documentation (man page & HTML)
1. decide on best implementation language (Ruby, Go, Rust, etc.)

## Development

After checking out the repo, you can make changes to the `repl.rb` Ruby version or the `repl.go` Go version to change its behaviour, and/or to the `repl.1.ronn` Markdown file to update the documentation.

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/pvdb/repl>.

## Credits

Chris Wanstrath :: [@defunkt](https://github.com/defunkt)

Check out his _(awesome, but unmaintained)_ [original version](https://github.com/defunkt/repl) on which this one is based!
