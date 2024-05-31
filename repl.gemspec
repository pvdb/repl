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

  spec.required_ruby_version = ['>= 3.2.0', '< 4.0.0']

  spec.files = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`
      .split("\x0")
      .reject { |f| f.match(%r{^(test|spec|features)/}) }
  end
  spec.bindir        = '.'
  spec.executables   = ['repl']
  spec.require_paths = ['lib']
end
