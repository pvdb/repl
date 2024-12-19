#!/usr/bin/env ruby
# frozen_string_literal: true

##
# run repl as an external command for git, brew, rbenv, etc.
name = File.basename($PROGRAM_NAME)       # `git-repl`, `brew-repl`, etc.
match = name.match(/\A(?<cmd>.+)-repl\z/) # <cmd> === `git`, `brew`, etc.
exec('repl', match[:cmd], *ARGV) if match # 'git-repl *' --> 'repl git *'

require 'mkmf'
require 'English'
require 'benchmark'
require 'fileutils'
require 'shellwords'
require 'io/console'

##
# script-specific Kernel extension
module Kernel
  def which(executable)
    MakeMakefile.find_executable0(executable)
  end
end

##
# script-specific logging method
class IO
  def console.log(message, color, width: 78, debug: false, quiet: false)
    # suppress output when in quiet mode
    return if quiet

    # deduct length of '===[]==' (7 characters)
    max_msg_length = width < 78 ? 71 : (width - 7)

    # truncate message if too long
    log_msg = if message.length > max_msg_length
                "#{message[0, (max_msg_length - 3)]}..."
              else
                message
              end

    # add padding to message to fill overall width
    padding = '=' * (max_msg_length - log_msg.length)

    puts debug ? "===[#{log_msg.send(color)}]==#{padding}" : message
  end
end

##
# script-specific colorization methods
class String
  def colorize(color_code)
    "\e[#{color_code}m#{self}\e[0m"
  end

  def strip_ansi
    gsub(/\e\[(\d+)(;\d+)*m/, '')
  end

  # rubocop:disable Style/SingleLineMethods
  def bold()   colorize('1'); end
  def invert() colorize('7'); end

  def red()    colorize('31'); end
  def green()  colorize('32'); end
  def yellow() colorize('33'); end
  def blue()   colorize('34'); end
  # rubocop:enable Style/SingleLineMethods

  def true?
    strip == 'true'
  end

  def comment?
    strip.start_with? '#'
  end
end

##
# is rlwrap utility installed?
def rlwrap_installed?
  !which('rlwrap').nil?
end

##
# version of rlwrap utility
def rlwrap_version
  `rlwrap --version`.chomp
rescue Errno::ENOENT
  'rlwrap not found'
end

##
# is repl running "inside" rlwrap?
def repl_wrapped?
  !ENV['__RLWRAP_REPL__'].nil?
end

##
# is repl running "inside" a pipeline?
def interactive?
  $stdin.tty?
end

##
# version of Ruby runtime
def runtime_version
  "#{RUBY_ENGINE} #{RUBY_VERSION}"
end

##
# version of repl script
def repl_version
  "repl 1.0.0 (#{rlwrap_version}, #{runtime_version})"
end

##
# short help message for repl script
def repl_help
  DATA.read
end

##
# directory containing command-specific history files
def repl_history_dir
  default_dir = Dir.home
  history_dir = ENV.fetch('REPL_HISTORY_DIR', default_dir)
  File.expand_path(history_dir)
end

##
# command-specific rlwrap history file
def history_file_for(command)
  history_dir = repl_history_dir
  history_file = File.join(history_dir, ".#{command}_history")

  # check if history_dir exists
  File.directory?(history_dir) ? history_file : nil
end

##
# directory containing command-specific completion files
def repl_completion_dir
  default_dir = File.join(Dir.home, '.repl')
  completion_dir = ENV.fetch('REPL_COMPLETION_DIR', default_dir)
  File.expand_path(completion_dir)
end

##
# command-specific rlwrap completion file
def completion_file_for(command)
  completion_dir = repl_completion_dir
  completion_file = File.join(completion_dir, command)

  # check if completion_file exists
  File.exist?(completion_file) ? completion_file : nil
end

##
# command-specific rlwrap options
def rlwrap_options_for(command)
  [].tap do |rlwrap_options|
    # suppress all default rlwrap break characters
    # specifically the '-' (hyphen/dash) character
    # note that whitespace is always word-breaking
    (rlwrap_options << '-b' << "''")

    history_file = history_file_for(command)
    (rlwrap_options << '-H' << history_file) unless history_file.nil?

    completion_file = completion_file_for(command)
    (rlwrap_options << '-f' << completion_file) unless completion_file.nil?
  end
end

##
# path to the REPL configuration file
def repl_conf
  default_path = File.join(Dir.home, '.repl.conf')
  conf_path = ENV.fetch('REPL_CONF', default_path)
  File.expand_path(conf_path)
end

##
# update ENV with config from ~/.repl.conf
def process_conf
  return unless File.file?(repl_conf)

  File.readlines(repl_conf, chomp: true).each do |line|
    next if line.empty?
    next if line.comment?

    key, value = line.split(/\s*=\s*/, 2).map(&:strip)

    # strip surrounding quotes
    value.delete_prefix!('"')
    value.delete_suffix!('"')

    # ENV takes precedence over ~/.repl.conf
    ENV[key] ||= value
  end
end

if ARGV.delete('--version')
  IO.console.puts(repl_version)
  exit(0)
end

if ARGV.delete('--help')
  IO.console.puts(repl_help)
  exit(0)
end

if ARGV.delete('--man')
  exec('man', File.join(__dir__, 'repl.1'))
  exit(0) # superfluous, really...
end

if ARGV.delete('--html')
  exec('open', File.join(__dir__, 'repl.1.html'))
  exit(0) # superfluous, really...
end

# list of remaining options (pipeline and logging)
OPTIONS = %w[--stdin --printf --escape --debug --quiet].freeze

# process ~/.repl.conf file
process_conf unless repl_wrapped?

# replace process with `rlwrap`-ed version
# if `rlwrap` is installed and also `repl`
# is running interactively (and not piped)
if interactive? && !repl_wrapped? && rlwrap_installed?
  rlwrap_options = []

  # rubocop:disable Style/SoleNestedConditional
  unless (base_cmd = (ARGV - OPTIONS).map(&:strip).first).nil?
    if File.executable?(base_cmd) || !which(base_cmd).nil?
      base_cmd = File.basename(base_cmd)
      rlwrap_options = rlwrap_options_for(base_cmd)
    end
  end
  # rubocop:enable Style/SoleNestedConditional

  ENV['__RLWRAP_REPL__'] = Process.pid.to_s
  exec(which('rlwrap'), *rlwrap_options, $PROGRAM_NAME, *ARGV)
end

# process the pipeline options
stdin  = ARGV.delete('--stdin')
printf = ARGV.delete('--printf')
escape = ARGV.delete('--escape')

# process the logging options (can be set in ~/.repl.conf or ENV)
debug  = ARGV.delete('--debug') || ENV['REPL_DEBUG']&.true?
quiet  = ARGV.delete('--quiet') || ENV['REPL_QUIET']&.true?

# command is whatever's left
if (cmd_string = ARGV.join(' ').strip).empty?
  IO.console.puts('No command specified... use `--help`')
  exit(0)
end

cmd_template = if stdin
                 pipe_cmd = which(printf ? 'printf' : 'echo')
                 "#{pipe_cmd} \"%s\" | #{cmd_string}"
               elsif ARGV.grep(/%s/).any?
                 cmd_string # the '%s' is embedded
               else
                 "#{cmd_string} %s"
               end

cmd_prompt = if debug
               if repl_wrapped?
                 "rlwrap(repl(\"#{cmd_template.blue}\"))"
               else
                 "repl(\"#{cmd_template.blue}\")"
               end
             else
               "\"#{cmd_template.blue}\""
             end

log_options = {
  debug: debug, quiet: quiet,
  width: cmd_prompt.strip_ansi.length
}.freeze

repl_prompt = ENV.fetch('REPL_PROMPT', '>>')
full_prompt = "#{cmd_prompt} #{repl_prompt}"

# rubocop:disable Metrics/BlockLength
loop do
  # prompt user for cmd arguments
  IO.console.print(full_prompt, ' ') if interactive? || !quiet

  line = begin
    $stdin.gets&.strip # nil when ^D
  rescue Interrupt
    nil # Interrupt is raised for ^C
  end

  # echo input if read from piped stdin
  IO.console.puts(line) unless interactive? || quiet

  # terminate `repl` loop on (^C|^D|EOF)
  break unless line

  unless line.empty? || line.comment?
    line = Shellwords.escape(line) if escape

    # command = format(cmd_template, line)  # expand single '%s' placeholder
    command = cmd_template.gsub('%s', line) # expand _all_ '%s' placeholders

    begin
      # print "expanded" command to be executed
      message = "Executing: '#{command}'"
      IO.console.log(message, :blue, **log_options)

      tms = Benchmark.measure do
        system(command, exception: true)
      rescue Interrupt
        # print message when command is interrupted
        message = 'Command was interrupted'
        IO.console.log(message, :red, **log_options)
      end

      # print elapsed real time to execute command
      message = format('Command took %.2fs to execute', tms.real)
      IO.console.log(message, :green, **log_options)
    rescue RuntimeError, Errno::ENOENT
      # print exception message when command fails
      message = $ERROR_INFO.message
      IO.console.log(message, :red, **log_options)

      next unless line =~ /(quit|exit)/i

      # print message when command fails due to 'quit'/'exit'
      message = 'Use ^C or ^D to exit repl'
      IO.console.log(message, :yellow, **log_options)
    end
  end

  # empty separator line
  IO.console.puts if interactive? || !quiet
end
# rubocop:enable Metrics/BlockLength

exit(0) # cleanly exit repl

__END__
Usage: repl [options] command ...

Options:
  --version Display repl version information
  --help    Display repl usage information
  --man     Display the repl man page
  --html    Open HTML version of man page
  --stdin   Pipe input to command's STDIN
  --printf  Avoid newline chars in STDIN
  --escape  Shell escape user's input
  --debug   Display each command being executed
  --quiet   Don't echo the prompt in pipelines

Homepage:

  http://github.com/pvdb/repl

Bug reports, suggestions, updates:

  http://github.com/pvdb/repl/issues

That's all Folks!
