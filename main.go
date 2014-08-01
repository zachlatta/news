package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"pkg/text/template"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// A Command is an implementation of a news command.
// like news latest or news fetch.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'news help' output.
	Short string

	// Long is the short description shown int he 'news help' output.
	Long string

	// Flag is the set of flags specific to this command.
	Flag flag.FlagSet
}

// Name returns the command's name: the first word in the usage line
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'news help'.
var commands = []*Command{}

var exitStatus = 0
var exitMu sync.Mutex

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}
}

var usageTemplate = `news is a utility for reading Hacker News.

Usage:

	news command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "news help [command]" for more information about a command.

`

var helpTemplate = `usage: news {{.UsageLine}}

{{.Long | trim}}
`

// tmpl executes the given template text on data, writing the results to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace,
		"capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'news helper'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr,
			"usage: news help command\n\nTo many arguments given.\n")
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2 : succeeced at 'go help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q. Run 'news help'.\n", arg)
	os.Exit(2) // failed at 'news help cmd'
}
