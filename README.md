# cmd-builder
Simple Go library to build `exec.Cmd` structs and execute them using the builder design pattern.

## Installation
Install with Go:
```
go get github.com/NoF0rte/cmd-builder@latest
```

## Usage
Here are some usage examples. (Explanations will come later. Didn't have time to create detailed documentation at the time)

```go
import "github.com/NoF0rte/cmd-builder"

// Example 1
output, err := builder.
	Cmd("git", cloneArgs...).
	Dir(dir).
	Env("GIT_SSH_COMMAND=ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no").
	Output()

// Example 2
err = builder.Cmd("sshfs", mountArgs...).Run()

// Example 3
err = builder.Cmd(args[0], args[1:]...).Interactive().Run()

// Example 4
cmd := builder.Cmd("scp", args...).
	Stdin(os.Stdin).
	Stdout(os.Stdout).
	Build()

// Example 5
factory := builder.NewFactory(builder.CmdFactoryOptions{
	Stdout: os.Stdout,
	Dir:    dir,
})
err = factory.Cmd("terraform", "init").Run()

// Example 6 
builder.Shell(fmt.Sprintf("sleep %d; vncviewer %s %s > /dev/null 2>&1", options.Delay, options.PasswordFile, options.Host)).Start()
```
