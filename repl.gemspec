require_relative 'lib/repl/version'

Gem::Specification.new do |spec|
  spec.name          = Repl::NAME
  spec.version       = Repl::VERSION
  spec.authors       = ['Peter Vandenberk']
  spec.email         = ['pvandenberk@mac.com']

  spec.summary       = 'Sometimes you need a REPL.'
  spec.description   = 'Complete rewrite of the awesome defunct/repl'
  spec.homepage      = 'https://github.com/pvdb/repl'
  spec.license       = 'MIT'

  spec.metadata['rubygems_mfa_required'] = 'true'

  spec.required_ruby_version = ['>= 2.7.0', '< 3.0.0']

  spec.files = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`
      .split("\x0")
      .reject { |f| f.match(%r{^(test|spec|features)/}) }
  end
  spec.bindir        = '.'
  spec.executables   = ['repl']
  spec.require_paths = ['lib']

  spec.add_development_dependency 'bundler', '~> 2.0'
  spec.add_development_dependency 'minitest', '~> 5.0'
  spec.add_development_dependency 'pry', '~> 0.13'
  spec.add_development_dependency 'pry-rescue', '~> 1.5'
  spec.add_development_dependency 'rake', '~> 13.0'
  spec.add_development_dependency 'ronn', '~> 0.7'
  spec.add_development_dependency 'rubocop', '~> 1.7'
  spec.add_development_dependency 'rubocop-minitest', '~> 0.10'
  spec.add_development_dependency 'rubocop-rake', '~> 0.5'
end
