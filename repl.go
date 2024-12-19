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
)

const replNameStr = "repl"
const replVersionStr = "1.0.0"

func rlwrapInstalled() bool {
	_, err := exec.LookPath("rlwrap")
	return err == nil
}

func rlwrapVersion() string {
	version, err := exec.Command("rlwrap", "--version").Output()
	if err != nil {
		return "rlwrap not found"
	}
	return strings.TrimSpace(string(version))
}

func runtimeVersion() string {
	return runtime.Version()
}

func replVersion() string {
	return fmt.Sprintf("%s %s (%s, %s)", replNameStr, replVersionStr, rlwrapVersion(), runtimeVersion())
}

func replHelp() string {
	return "Coming soon!"
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

func replWrapped() bool {
	_, found := os.LookupEnv("__RLWRAP_REPL__")
	return found
}

func replGetEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func replPrompt() string {
	return replGetEnv("REPL_PROMPT", ">>")
}

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
	// run repl as an external command for git, brew, etc.
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

	// run repl within rlwrap
	if !replWrapped() && rlwrapInstalled() {
		os.Setenv("__RLWRAP_REPL__", replPid())
		os.Args[0] = replFullPath()
		replaceProcess("rlwrap", os.Args)
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
		if !quiet {
			fmt.Print(fullPrompt, " ")
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			eof = errors.Is(err, io.EOF)
		}
		input = strings.TrimSpace(input)
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
			logMessage("sh -c '"+expanded+"'", "blue")
		}

		command := exec.Command("sh", "-c", expanded)

		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			logMessage(err.Error(), "red")

			match, _ := regexp.MatchString("(?i)quit|exit", input)

			if match {
				logMessage("use ^C or ^D to exit repl", "yellow")
			}
		}
	}
}
