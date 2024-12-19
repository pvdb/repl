package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// script-specific logging function
func logMessage(msg string, color string) {
	const maxMsgLength = 71 // 78 - len("===[]==")

	var logMsg string
	if len(msg) > maxMsgLength {
		logMsg = msg[:maxMsgLength-3] + "..."
	} else {
		logMsg = msg
	}

	padding := strings.Repeat("=", maxMsgLength-len(logMsg))

	fmt.Printf("===[%s]==%s\n", colorize(logMsg, color), padding)
}

// script-specific colorization methods
func colorize(message string, color string) string {
	var colors = map[string]string{
		"reset":  "\033[0m",
		"bold":   "\033[1m",
		"invert": "\033[7m",

		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"blue":   "\033[34m",
	}
	return colors[color] + message + colors["reset"]
}

func bold(text string) string   { return colorize(text, "bold") }
func invert(text string) string { return colorize(text, "invert") }

func red(text string) string    { return colorize(text, "red") }
func green(text string) string  { return colorize(text, "green") }
func yellow(text string) string { return colorize(text, "yellow") }
func blue(text string) string   { return colorize(text, "blue") }


// is rlwrap utility installed?
func rlwrapInstalled() bool {
	_, err := exec.LookPath("rlwrap")
	return err == nil
}

// version of rlwrap utility
func rlwrapVersion() string {
	version, err := exec.Command("rlwrap", "--version").Output()
	if err != nil {
		return "rlwrap not found"
	}
	return strings.TrimSpace(string(version))
}

// is repl running "inside" rlwrap?
func replWrapped() bool {
	_, found := os.LookupEnv("__RLWRAP_REPL__")
	return found
}

// is repl running "inside" pipeline?
func interactive() bool {
	// https://stackoverflow.com/a/26567513/525415
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// version of Go runtime
func runtimeVersion() string {
	return runtime.Version()
}

// version of repl script
func replVersion() string {
	return fmt.Sprintf("repl 1.0.0 (%s, %s)", rlwrapVersion(), runtimeVersion())
}

// short help message for repl script
func replHelp() string {
	return "Coming soon!"
}

// directory containing command-specific history files
func replHistoryDir() (historyDir string) {
	defaultDir := os.Getenv("HOME")
	historyDir = replGetEnv("REPL_HISTORY_DIR", defaultDir)
	historyDir, _ = filepath.Abs(historyDir)
	return
}

// command-specific rlwrap history file
func historyFileFor(command string) string {
	historyDir := replHistoryDir()
	historyFile := filepath.Join(historyDir, "." + command + "_history")

	// check if history *directory* exists
	_, err := os.Stat(historyDir)

	if err == nil {
		return historyFile
	} else {
		return ""
	}
}

// directory containing command-specific completion files
func replCompletionDir() (completionDir string) {
	defaultDir := filepath.Join(os.Getenv("HOME"), ".repl")
	completionDir = replGetEnv("REPL_COMPLETION_DIR", defaultDir)
	completionDir, _ = filepath.Abs(completionDir)
	return
}

// command-specific rlwrap completion file
func completionFileFor(command string) string {
	completionDir := replCompletionDir()
	completionFile := filepath.Join(completionDir, command)

	// check if completion *file* exists
	_, err := os.Stat(completionFile)

	if err == nil {
		return completionFile
	} else {
		return ""
	}
}

// command-specific rlwrap options
func rlwrapOptionsFor(command string) (options []string) {
	options = []string{}

	// suppress all default rlwrap break characters
	// specifically the '-' (hyphen/dash) character
	// note that whitespace is always word-breaking
	options = append(options, "-b", "''")

	historyFile := historyFileFor(command)
	if historyFile != "" {
		options = append(options, "-H", historyFile)
	}

	completionFile := completionFileFor(command)
	if completionFile != "" {
		options = append(options, "-f", completionFile)
	}

	return
}

// path to the REPL configuration file
func replConf() (confPath string) {
	defaultPath := filepath.Join(os.Getenv("HOME"), ".repl.conf")
	confPath = replGetEnv("REPL_CONF", defaultPath)
	confPath, _ = filepath.Abs(confPath)
	return
}

// update ENV with config from ~/.repl.conf
func processConf() {
	if config, err := os.ReadFile(replConf()); err == nil {
		lines := strings.Split(string(config), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, value, _ := strings.Cut(line, "=")
			replSetEnv(key,	strings.Trim(value, "\""))
		}
	}
}


func replPid() string {
	return strconv.Itoa(os.Getpid())
}

func replFullPath() (fullPath string) {
	fullPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func replInstallDir() string {
	return filepath.Dir(replFullPath())
}

func replSetEnv(key string, value string) {
	// ENV takes precedence over ~/.repl.conf
	if _, found := os.LookupEnv(key); !found {
		os.Setenv(key, value)
	}
}

func replGetEnv(key string, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defaultValue
}

func replPrompt() string {
	return replGetEnv("REPL_PROMPT", ">>")
}

func shellEscape(str string) string {
	// https://github.com/alessio/shellescape/blob/v1.4.2/shellescape.go
	var pattern *regexp.Regexp
	pattern = regexp.MustCompile(`[^\w@%+=:,./-]`)

	if len(str) == 0 {
		return "''"
	}

	if pattern.MatchString(str) {
		return "'" + strings.ReplaceAll(str, "'", "'\"'\"'") + "'"
	}

	return str
}

func replaceProcess(cmd string, args []string) {
	cmdPath, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal(err)
	}

	// prepend cmd to the args array
	// for correct Exec() invocation
	args = append([]string{cmd}, args...)

	err = syscall.Exec(cmdPath, args, os.Environ())
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// run repl as an external command for git, brew, rbenv, etc.
	name := filepath.Base(os.Args[0])
	re := regexp.MustCompile(`\A(?P<cmd>.+)-repl\z`)
	if match := re.FindStringSubmatch(name); match != nil {
		// if:   name == `git-repl`, `brew-repl`, etc.
		// then: <cmd> == `git`, `brew`, etc.
		// and:  'git-repl [args]' >=> 'repl git [args]', etc.
		args := append([]string{match[1]}, os.Args[1:]...)
		replaceProcess("repl", args)
	}

	// show repl version and exit
	if slices.Contains(os.Args, "--version") {
		fmt.Println(replVersion())
		os.Exit(0)
	}

	// show repl help and exit
	if slices.Contains(os.Args, "--help") {
		fmt.Println(replHelp())
		os.Exit(0)
	}

	// show repl man page (and "exit")
	if slices.Contains(os.Args, "--man") {
		manPath := filepath.Join(replInstallDir(), "repl.1")
		replaceProcess("man", []string{manPath})
	}

	// show repl html page (and "exit")
	if slices.Contains(os.Args, "--html") {
		htmlPath := filepath.Join(replInstallDir(), "repl.1.html")
		replaceProcess("open", []string{htmlPath})
	}

	// process ~/.repl.conf file
	if !replWrapped() {
		processConf()
	}

	// replace process with `rlwrap`-ed version
	// if `rlwrap` is installed and also `repl`
	// is running interactively (and not piped)
	if interactive() && !replWrapped() && rlwrapInstalled() {
		command := os.Args[1] // FIXME: remove repl options first

		// replace: "repl command [args]"
		// with: "rlwrap [options] /usr/local/bin/repl command [args]"
		rlwrapArgs := slices.Clone(os.Args)
		rlwrapArgs[0] = replFullPath()
		rlwrapOptions := rlwrapOptionsFor(command)
		rlwrapArgs = slices.Insert(rlwrapArgs, 0, rlwrapOptions...)

		os.Setenv("__RLWRAP_REPL__", replPid())
		replaceProcess("rlwrap", rlwrapArgs)
	}
}

func main() {
	var stdin, printf, escape, debug, quiet bool

	debug, _ = strconv.ParseBool(replGetEnv("REPL_DEBUG", "false"))
	quiet, _ = strconv.ParseBool(replGetEnv("REPL_QUIET", "false"))

	os.Args = slices.DeleteFunc(os.Args, func(arg string) bool {
		switch arg {
		case "--stdin":
			stdin = true
			return true
		case "--printf":
			printf = true
			return true
		case "--escape":
			escape = true
			return true
		case "--debug":
			debug = true
			return true
		case "--quiet":
			quiet = true
			return true
		default:
			return false
		}
	})

	var interactive bool = interactive()

	cmdString := strings.TrimSpace(strings.Join(os.Args[1:], " "))

	if cmdString == "" {
		fmt.Println("No command specified... use `--help`")
		os.Exit(1)
	}

	var cmdTemplate string
	if stdin {
		pipeCmd := "echo"
		if printf {
			pipeCmd = "printf"
		}
		cmdTemplate = fmt.Sprintf("%s \"%%s\" | %s", pipeCmd, cmdString)
	} else if strings.Contains(cmdString, "%s") {
		cmdTemplate = cmdString
	} else {
		cmdTemplate = cmdString + " %s"
	}

	var cmdPrompt string
	if debug {
		if replWrapped() {
			cmdPrompt = "rlwrap(repl(\"" + blue(cmdTemplate) + "\"))"
		} else {
			cmdPrompt = "repl(\"" + blue(cmdTemplate) + "\")"
		}
	} else {
		cmdPrompt = "\"" + blue(cmdTemplate) + "\""
	}

	fullPrompt := fmt.Sprintf("%s %s", cmdPrompt, replPrompt())

	reader := bufio.NewReader(os.Stdin)

	eof := false
	for !eof {
		// prompt user for cmd arguments
		if interactive || !quiet {
			fmt.Print(fullPrompt, " ")
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			eof = errors.Is(err, io.EOF)
		}
		input = strings.TrimSpace(input)

		// echo input if read from piped stdin
		if !interactive && !quiet {
			fmt.Println(input)
		}

		if input == "" {
			continue
		}
		if strings.HasPrefix(input, "#") {
			continue
		}

		if escape {
			input = shellEscape(input)
		}

		expanded := strings.ReplaceAll(cmdTemplate, "%s", input)
		if debug {
			// print "expanded" command to be executed
			logMessage("sh -c '"+expanded+"'", "blue")
		}

		command := exec.Command("sh", "-c", expanded)

		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		startTime := time.Now()
		err = command.Run()
		duration := time.Since(startTime)

		if err == nil {
			if debug {
				// print elapsed real time to execute command
				logMessage(fmt.Sprintf("Command took %.2fs to execute", duration.Seconds()), "green")
			}
		} else {
			// print exception message when command fails
			logMessage(err.Error(), "red")

			match, _ := regexp.MatchString("(?i)quit|exit", input)

			if match {
				// print message when command fails due to 'quit'/'exit'
				logMessage("use ^C or ^D to exit repl", "yellow")
			}
		}

		// empty separator line
		if interactive || !quiet {
			fmt.Println()
		}
	}
}
