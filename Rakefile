# rubocop:disable Style/SymbolArray
# rubocop:disable Style/HashSyntax

require 'bundler/setup'
require 'bundler/gem_tasks'

task :validate_gemspec do
  Bundler.load_gemspec('repl.gemspec').validate
end

task :version => :validate_gemspec do
  puts Repl.version
  puts `#{Bundler.root.join('repl')} --version`.chomp
end

require 'rubocop/rake_task'

RuboCop::RakeTask.new(:rubocop)

require 'rake/testtask'

Rake::TestTask.new(:test) do |t|
  t.libs << 'test'
  t.libs << 'lib'
  t.test_files = FileList['test/**/*_test.rb']
end

task :default => [:version, :rubocop, :test]

task :documentation do
  # update script version to match gem version
  repl = Bundler.root.join('repl')

  # rubocop:disable Style/RedundantStringEscape
  sed = "s/repl [0-9]+\.[0-9]+\.[0-9]+/#{Repl.repl_version}/"
  system "sed -E -i \'\' -e \"#{sed}\" #{repl}"
  # rubocop:enable Style/RedundantStringEscape

  # generate ROFF and HTML versions of man page
  ronn = Bundler.root.join('repl.1.ronn')
  # rubocop:disable Style/TrailingCommaInArrayLiteral
  options = [
    '--organization=PVDB',
    '--manual="Awesome Utilities"',
    "--date=#{File.mtime(Bundler.root).strftime('%F')}",
  ].join(' ')
  # rubocop:enable Style/TrailingCommaInArrayLiteral
  system "ronn --roff --html #{options} #{ronn}"
end

task :ready => :documentation do
  sh('bundle --quiet') # regenerate Gemfile.lock e.g. if version has changed
  sh('git update-index --really-refresh') # refresh touched but unchanged docs
  sh('git diff-index --quiet HEAD --') # https://stackoverflow.com/a/2659808
end

Rake::Task['build'].enhance([:default, :ready])

# rubocop:enable Style/HashSyntax
# rubocop:enable Style/SymbolArray
