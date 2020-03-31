lib = File.expand_path('lib', __dir__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'repl/version'

Gem::Specification.new do |spec|
  spec.name          = 'repl'
  spec.version       = Repl::VERSION
  spec.authors       = ['Peter Vandenberk']
  spec.email         = ['pvandenberk@mac.com']

  spec.summary       = 'Sometimes you need a REPL.'
  spec.description   = 'Complete rewrite of the awesome defunct/repl'
  spec.homepage      = 'https://github.com/pvdb/consenter'

  spec.files         = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`
      .split("\x0")
      .reject { |f| f.match(%r{^(test|spec|features)/}) }
  end
  spec.bindir        = 'exe'
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ['lib']

  spec.add_development_dependency 'bundler', '~> 2.0'
  spec.add_development_dependency 'minitest', '~> 5.0'
  spec.add_development_dependency 'pry', '~> 0.13'
  spec.add_development_dependency 'rake', '~> 13.0'
  spec.add_development_dependency 'rubocop', '~> 0.80'
end
