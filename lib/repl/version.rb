module Repl
  NAME = 'repl'.freeze
  VERSION = '0.8.0'.freeze

  def self.repl_version
    "#{NAME} #{VERSION}"
  end

  def self.rlwrap_version
    `rlwrap --version`.chomp
  rescue Errno::ENOENT
    'rlwrap not installed'
  end

  def self.version
    "#{repl_version} (#{rlwrap_version})"
  end

  V_REGEXP = 'repl [0-9]+\.[0-9]+(\.[0-9]+)? \(rlwrap [0-9]+\.[0-9]+(\.[0-9]+)?\)'.freeze

  def self.v_regexp
    @v_regexp ||= Regexp.compile(V_REGEXP)
  end
end
