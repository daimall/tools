package mobilecmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/daimall/tools/curd/customerror"
	"github.com/kballard/go-shellquote"
)

// Command add timeout support for os/exec
type Command struct {
	Args       []string
	Timeout    time.Duration
	Shell      bool
	ShellQuote bool
	Stdout     io.Writer
	Stderr     io.Writer
}

func NewCommand(args ...string) *Command {
	return &Command{
		Args: args,
	}
}

func (c *Command) shellPath() string {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh, err := exec.LookPath("sh")
		if err == nil {
			return sh
		}
		sh = "/system/bin/sh"
	}
	return sh
}

func (c *Command) computedArgs() (name string, args []string) {
	if c.Shell {
		var cmdline string
		if c.ShellQuote {
			cmdline = shellquote.Join(c.Args...)
		} else {
			cmdline = strings.Join(c.Args, " ") // simple, but works well with ">". eg Args("echo", "hello", ">output.txt")
		}
		args = append(args, "-c", cmdline)
		return c.shellPath(), args
	}
	return c.Args[0], c.Args[1:]
}

func (c Command) newCommand() *exec.Cmd {
	name, args := c.computedArgs()
	cmd := exec.Command(name, args...)
	if c.Stdout != nil {
		cmd.Stdout = c.Stdout
	}
	if c.Stderr != nil {
		cmd.Stderr = c.Stderr
	}
	return cmd
}

func (c Command) Run() error {
	cmd := c.newCommand()
	if c.Timeout > 0 {
		timer := time.AfterFunc(c.Timeout, func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		})
		defer timer.Stop()
	}
	return cmd.Run()
}

func (c Command) StartBackground() (pid int, err error) {
	cmd := c.newCommand()
	err = cmd.Start()
	if err != nil {
		return
	}
	pid = cmd.Process.Pid
	return
}

func (c Command) Output() (output []byte, err error) {
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = nil
	err = c.Run()
	return b.Bytes(), err
}

func (c Command) CombinedOutput() (output string, err customerror.CustomError) {
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b
	verr := c.Run()
	if verr != nil {
		err = customerror.New(cmdError2Code(verr), b.String())
	}
	return b.String(), err
}

func (c Command) CombinedOutputString() (output string, err error) {
	bytesOutput, err := c.CombinedOutput()
	return string(bytesOutput), err
}

func cmdError2Code(err error) int {
	if err == nil {
		return 0
	}
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 128
}
