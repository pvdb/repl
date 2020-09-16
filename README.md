# repl(1) -- sometimes you need a repl

`repl` is an interactive program which tenderly wraps another,
non-interactive program.

This version borrows _(very!)_ heavily from but is a ground-up rewrite of
[Chris Wanstrath][defunkt]'s awesome original version found in [defunkt/repl][].

[defunkt/repl]: https://github.com/defunkt/repl
[pvdb/repl]: https://github.com/pvdb/repl
[defunkt]: https://github.com/defunkt

## Installation

`repl` can be installed as a RubyGem:

    $ gem install repl

That way you can run `repl` without any changes to your system's `$PATH`.

## Usage

See [defunkt/repl][] for usage instructions and examples

## Features

COMING SOON

## TODO

COMING SOON

## Development

After checking out the repo, run `bin/setup` to install dependencies. Then, run `rake test` to run the tests. You can also run `bin/console` for an interactive prompt that will allow you to experiment.

To install this gem onto your local machine, run `bundle exec rake install`. To release a new version, update the version number in `version.rb`, and then run `bundle exec rake release`, which will create a git tag for the version, push git commits and tags, and push the `.gem` file to [rubygems.org](https://rubygems.org).

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/pvdb/repl>.
