module Repl
  NAME = 'repl'.freeze
  VERSION = '1.0.0'.freeze

  def self.repl_version
    "#{NAME} #{VERSION}"
  end

  def self.rlwrap_version
    `rlwrap --version`.chomp
  rescue Errno::ENOENT
    'rlwrap not installed'
  end

  def self.version
    @version ||= "#{repl_version} (#{rlwrap_version})"
  end

  SEMVER_REGEXP = '\d+\.\d+(\.\d+)?'.freeze

  REPL_VERSION = "repl #{SEMVER_REGEXP}".freeze
  RLWRAP_VERSION = "rlwrap (#{SEMVER_REGEXP}|not installed)".freeze

  VERSION_REGEXP = "#{REPL_VERSION} \\(#{RLWRAP_VERSION}\\)".freeze

  def self.version_regexp
    @version_regexp ||= Regexp.compile(VERSION_REGEXP)
  end
end
