package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/subrat-dwi/passman-cli/cmd"
	"github.com/subrat-dwi/passman-cli/internal/agent"
	"github.com/subrat-dwi/passman-cli/internal/app"
)

func main() {

	// Only start agent server when command is exactly "pman agent" (no subcommands)
	if len(os.Args) == 2 && os.Args[1] == "agent" {
		agent.Run()
		return
	}

	ensureAgentRunning()
	app := app.New()

	rootCmd := cmd.NewRootCmd(app)
	rootCmd.Execute()
}

func ensureAgentRunning() {
	_, _, err := agent.Status()
	if err == nil {
		return
	}

	// start agent as a detached background process
	cmd := exec.Command(os.Args[0], "agent")
	// if runtime.GOOS == "windows" {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{
	// 		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	// 	}
	// }), 600
	cmd.Start()

	// wait briefly (avoid race)
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		if _, _, err := agent.Status(); err == nil {
			return
		}
	}
	log.Fatal("failed to start agent")
}
