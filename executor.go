package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	. "github.com/Bo0mer/executor/model"
	. "github.com/Bo0mer/executor/utils"
)

type Executor interface {
	Execute(c Command) (stdout, stderr string, exitCode int, err error)
	ExecuteAsync(c Command) <-chan CommandResult
}

type executor struct{}

func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) Execute(c Command) (string, string, int, error) {
	result := <-e.executeCommandAsync(c)
	return result.Stdout, result.Stderr, result.ExitCode, result.Error
}

func (e *executor) ExecuteAsync(c Command) <-chan CommandResult {
	return e.executeCommandAsync(c)
}

func (e *executor) executeCommandAsync(c Command) <-chan CommandResult {
	result := make(chan CommandResult)

	go func() {
		execCmd, stdout, stderr := e.buildExecCommand(c)

		err := execCmd.Start()
		if err != nil {
			result <- CommandResult{
				Stdout:   "",
				Stderr:   "",
				ExitCode: -1,
				Error:    err,
			}
		}

		err = execCmd.Wait()

		exitCode := -1

		waitStatus := execCmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = waitStatus.ExitStatus()

		result <- CommandResult{
			Stdout:   string(stdout.SniffedData()),
			Stderr:   string(stderr.SniffedData()),
			ExitCode: exitCode,
			Error:    err,
		}
	}()

	return result

}

func (e *executor) buildExecCommand(c Command) (*exec.Cmd, Sniffer, Sniffer) {
	execCmd := exec.Command(c.Name, c.Args...)

	stdoutSniffer := e.buildSniffer(c.Stdout)
	execCmd.Stdout = stdoutSniffer
	stderrSniffer := e.buildSniffer(c.Stderr)
	execCmd.Stderr = stderrSniffer

	if c.WorkingDir != "" {
		execCmd.Dir = c.WorkingDir
	}

	env := []string{}

	if !c.UseIsolatedEnv {
		env = os.Environ()
	}

	for name, value := range c.Env {
		env = append(env, fmt.Sprintf("%s=%s", name, value))
	}
	execCmd.Env = env

	return execCmd, stdoutSniffer, stderrSniffer
}

func (e *executor) buildSniffer(writer io.Writer) Sniffer {
	if writer == nil {
		return NewSniffer(ioutil.Discard)
	}
	return NewSniffer(writer)
}
