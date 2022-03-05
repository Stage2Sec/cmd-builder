package builder

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CmdFactory allows you to create builder structs that
// use the same options
type CmdFactory struct {
	Options CmdFactoryOptions
}

// CmdFactoryOptions represents the configurable options for creating builders
// with the CmdFactory
type CmdFactoryOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Dir    string
	Env    []string
}

// NewFactory creates a new CmdFactory struct with the specified CmdFactoryOptions
func NewFactory(options CmdFactoryOptions) CmdFactory {
	return CmdFactory{
		Options: options,
	}
}

// Cmd returns the CmdBuilder struct built from the factory's options,
// that can be used to build/execute 'exec.Cmd` structs.
func (factory CmdFactory) Cmd(name string, args ...string) *CmdBuilder {
	builder := Cmd(name, args...)
	if factory.Options.Stdin != nil {
		builder.cmd.Stdin = factory.Options.Stdin
	}

	if factory.Options.Stdout != nil {
		builder.cmd.Stdout = factory.Options.Stdout
	}

	if factory.Options.Stderr != nil {
		builder.cmd.Stderr = factory.Options.Stderr
	}

	if factory.Options.Dir != "" {
		builder.cmd.Dir = factory.Options.Dir
	}

	if len(factory.Options.Env) > 0 {
		builder.cmd.Env = append(builder.cmd.Env, factory.Options.Env...)
	}

	return builder
}

// Shell is like Cmd except it passes the arg string to the OS shell.
//
// Linux: 'bash -c'
//
// macOS: 'zsh -c'
//
// Windows: 'powershell -Command'
//
// Everything else: '$SHELL -c'
func (factory CmdFactory) Shell(args string) *CmdBuilder {
	switch runtime.GOOS {
	default:
		return factory.Cmd(os.Getenv("SHELL"), "-c", args)
	case "linux":
		return factory.Cmd("bash", "-c", args)
	case "darwin":
		return factory.Cmd("zsh", "-c", args)
	case "windows":
		return factory.Cmd("powershell", "-Command", args)
	}
}

// CmdBuilder represents an 'exec.Cmd' struct using the builder design pattern
type CmdBuilder struct {
	cmd *exec.Cmd
}

// Cmd returns the CmdBuilder struct that can be used to build/execute 'exec.Cmd` structs.
func Cmd(name string, args ...string) *CmdBuilder {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return &CmdBuilder{
		cmd: cmd,
	}
}

// Shell is like Cmd except it passes the arg string to the OS shell.
//
// Linux: 'bash -c'
//
// macOS: 'zsh -c'
//
// Windows: 'powershell -Command'
//
// Everything else: '$SHELL -c'
func Shell(args string) *CmdBuilder {
	switch runtime.GOOS {
	default:
		return Cmd(os.Getenv("SHELL"), "-c", args)
	case "linux":
		return Cmd("bash", "-c", args)
	case "darwin":
		return Cmd("zsh", "-c", args)
	case "windows":
		return Cmd("powershell", "-Command", args)
	}
}

// Dir specifies the working directory of the command.
// If Dir is the empty string, the command will run in the
// in calling process's current directory.
func (cmdBuilder *CmdBuilder) Dir(dir string) *CmdBuilder {
	cmdBuilder.cmd.Dir = dir
	return cmdBuilder
}

// Stdout sets the command's stdout to the specified writer. Passing nil is the same as
// passing os.DevNull
func (cmdBuilder *CmdBuilder) Stdout(stdout io.Writer) *CmdBuilder {
	cmdBuilder.cmd.Stdout = stdout
	return cmdBuilder
}

// Stderr sets the command's stderr to the specified writer. Passing nil is the same as
// passing os.DevNull
func (cmdBuilder *CmdBuilder) Stderr(stderr io.Writer) *CmdBuilder {
	cmdBuilder.cmd.Stderr = stderr
	return cmdBuilder
}

// Stdin sets the command's stdin to the specified reader. Passing nil is the same as
// passing os.DevNull
func (cmdBuilder *CmdBuilder) Stdin(stdin io.Reader) *CmdBuilder {
	cmdBuilder.cmd.Stdin = stdin
	return cmdBuilder
}

// Interactive sets the stdin, stdout, and stderr to the OS's
// stdin, stdout, and stderr
func (cmdBuilder *CmdBuilder) Interactive() *CmdBuilder {
	cmdBuilder.cmd.Stderr = os.Stderr
	cmdBuilder.cmd.Stdin = os.Stdin
	cmdBuilder.cmd.Stdout = os.Stdout
	return cmdBuilder
}

// Env specifies the environment of the process.
// Each entry is of the form "key=value".
// If Env is nil, the new process uses the current process's
// environment.
func (cmdBuilder *CmdBuilder) Env(vars ...string) *CmdBuilder {
	cmdBuilder.cmd.Env = append(cmdBuilder.cmd.Env, vars...)
	return cmdBuilder
}

// Build returns the built *exec.Cmd struct
func (cmdBuilder *CmdBuilder) Build() *exec.Cmd {
	return cmdBuilder.cmd
}

// Start starts the specified command but does not wait for it to complete.
func (cmdBuilder *CmdBuilder) Start() error {
	return cmdBuilder.cmd.Start()
}

// Run starts the specified command and waits for it to complete.
func (cmdBuilder *CmdBuilder) Run() error {
	return cmdBuilder.cmd.Run()
}

// Output runs the command and returns its standard output.
// Any returned error will usually be of type *ExitError.
func (cmdBuilder *CmdBuilder) Output() (string, error) {
	var output []byte
	var err error

	// if cmd.Stdout is already specified then Output() errors out with "exec: Stdout already set"
	if cmdBuilder.cmd.Stdout != nil {
		var outBuf bytes.Buffer
		cmdBuilder.cmd.Stdout = io.MultiWriter(cmdBuilder.cmd.Stdout, &outBuf)
		err = cmdBuilder.cmd.Run()
		if err != nil {
			return "", err
		}
		output = outBuf.Bytes()
	} else {
		output, err = cmdBuilder.cmd.Output()
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(string(output)), nil
}

// Lines is like Output except it will split by new lines
func (cmdBuilder *CmdBuilder) Lines() ([]string, error) {
	output, err := cmdBuilder.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n"), nil
}
