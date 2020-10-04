# repl(1) -- sometimes you *really* need a REPL

`repl` wraps non-interactive commands in an interactive REPL _(read-eval-print loop)_.

## Synopsis

    repl [options] command ...

`command` is executed using the lines you type into `repl`'s prompt as either command-line arguments or else as standard input, and anything written by `command` to standard output and standard error is displayed.

This is repeated until you exit out of `repl`'s interactive loop by using either `CTRL-C` or `CTRL-D` at the prompt (default is `>> `):

    $ repl host -t A
    >> google.com
    google.com has address 216.58.212.206
    
    >> google.co.uk
    google.co.uk has address 172.217.169.35
    
    >> ^D
    $ _

## Usage

`repl` is meant to wrap programs which accept and process command line arguments or which read from and process standard input, and which in turn print to standard output or standard error, but which don't have an interactive REPL of their own.

If you have [`rlwrap(1)`][rlwrap] installed you'll automatically get the full benefits of [`readline`][readline]: persistent history, reverse history searches, tab completion, command-specific completions, etc.

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

COMING SOON

## Features

COMING SOON

## Configuration

The following environment variables can be used to override `repl`'s defaults:

* `REPL_PROMPT`: prompt to use in `repl`'s read-eval-print loop (default: `>>`)
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
