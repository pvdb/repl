# repl(1) -- sometimes you *really* need a REPL

`repl` wraps non-interactive commands in an interactive REPL _(read-eval-print loop)_.

## Synopsis

    repl [options] command ...

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

It can be used with commands you wish had an interactive mode, but don't:

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

`repl` is easily installed as a standalone script somewhere in your `${PATH}`; the default install location is `/usr/local/bin`, but it can be anything on your `${PATH}`.

Run the following commands in your preferred shell:

    REPL_INSTALL_DIR=/usr/local/bin
    curl -s https://raw.githubusercontent.com/pvdb/repl/master/exe/repl -o "${REPL_INSTALL_DIR}/repl"
    chmod 755 "${REPL_INSTALL_DIR}/repl"

This way you can run `repl` without any changes to your system's `$PATH`.

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

By "double quoting" the command, you can wrap `repl` around entire command pipelines:

    $ repl --debug "'%s blegga|wc -c'"
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

COMING SOON

## Configuration

The following environment variables can be used to override `repl`'s defaults:

* `REPL_PROMPT`: prompt to use in `repl`'s read-eval-print loop (default: `>>`)
* `REPL_ESCAPE`: equivalent to the `--escape` option if set to true (default: false)
* `REPL_DEBUG` : equivalent to the `--debug` option if set to true (default: false)
* `REPL_CLEAR` : equivalent to the `--clear` option if set to true (default: false)
* `REPL_HISTORY_DIR`: directory in which history files are kept (default: `${HOME}`)
* `REPL_COMPLETION_DIR`: directory in which completion files are kept (default: `${HOME}/.repl`)

These options can also be set permantently in `${HOME}/.repl.rc`, instead of "polluting" your shell environment with them; use the [`repl.rc`](repl.rc) template file as a starting point.

## TODO

COMING SOON

## Development

After checking out the repo, run `bin/setup` to install dependencies. Then, run `rake test` to run the tests. You can also run `bin/console` for an interactive prompt that will allow you to experiment.

To install this gem onto your local machine, run `bundle exec rake install`. To release a new version, update the version number in `version.rb`, and then run `bundle exec rake release`, which will create a git tag for the version, push git commits and tags, and push the `.gem` file to [rubygems.org](https://rubygems.org).

## Build Status

[![Build Status](https://travis-ci.org/pvdb/repl.svg?branch=master)](https://travis-ci.org/pvdb/repl)

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/pvdb/repl>.
