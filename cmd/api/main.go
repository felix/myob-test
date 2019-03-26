package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"

	server "github.com/felix/ops-technical-test"
)

// Version is set during build
var Version string

func main() {
	m := NewMain()
	if err := m.Run(os.Args[1:]...); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Main represents the main program execution.
type Main struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Host string
	Port int
}

// NewMain returns a new instance of Main connected to the standard input/output.
func NewMain() *Main {
	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run executes the program.
func (m *Main) Run(args ...string) error {
	var err error
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	help := fs.Bool("h", false, "")

	fs.StringVar(&m.Host, "host", os.Getenv("HOST"), "")

	envPort := int64(8080)
	if port := os.Getenv("PORT"); port != "" {
		envPort, err = strconv.ParseInt(port, 10, 32)
		if err != nil {
			return err
		}
	}
	fs.IntVar(&m.Port, "port", int(envPort), "")

	fs.SetOutput(m.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	} else if *help {
		fmt.Fprintln(m.Stderr, m.Usage())
		return nil
	}

	listen := fmt.Sprintf("%s:%d", m.Host, m.Port)

	s, err := server.New(listen, server.SetVersion(Version))
	if err != nil {
		return err
	}

	// Watch for signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			select {
			case <-interrupt:
				fmt.Println("stopping")
				if err := s.Stop(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		}
	}()

	fmt.Fprintf(m.Stdout, "Starting server version %q at %s\n", Version, listen)
	s.Run()
	return nil
}

// Usage returns the help message.
func (m *Main) Usage() string {
	return strings.TrimLeft(`
A simple, small, operable web-style API.

Usage:

	server [options]

The options are:

	-h HOST
		Listen on HOST IP address. Default is 0.0.0.0.
	-p PORT
		Listen at PORT. Default is 8080.
	-h
		This help
`, "\n")
}
