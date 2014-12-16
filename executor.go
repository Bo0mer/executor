package executor

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
}

type executor struct{}

func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) Execute(c Command) (string, string, int, error) {
	execCmd, stdout, stderr := e.buildExecCommand(c)

	err := execCmd.Start()
	if err != nil {
		return "", "", -1, err
	}

	err = execCmd.Wait()

	exitStatus := -1

	waitStatus := execCmd.ProcessState.Sys().(syscall.WaitStatus)
	exitStatus = waitStatus.ExitStatus()

	return string(stdout.SniffedData()), string(stderr.SniffedData()), exitStatus, err
}

func (e *executor) buildExecCommand(c Command) (*exec.Cmd, Sniffer, Sniffer) {
	execCmd := &exec.Cmd{
		Path: c.Name,
		Args: c.Args,
	}

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
