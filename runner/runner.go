package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func run() bool {
	runnerLog("Running...")

	appArgs := os.Getenv("APP_ARGS")
	fmt.Println(appArgs)

	args := []string{}
	if appArgs != "" {
		args = strings.Split(appArgs, " ")
	}

	cmd := exec.Command(buildPath(), args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
