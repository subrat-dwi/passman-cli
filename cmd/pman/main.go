package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/subrat-dwi/passman-cli/cmd"
	"github.com/subrat-dwi/passman-cli/internal/agent"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/service"
)

// version is set at build time via ldflags
var version = "dev"

func main() {
	// Set the version for the service package
	service.SetVersion(version)

	// start agent server when command is exactly "pman agent"
	if len(os.Args) == 2 && os.Args[1] == "agent" {
		agent.Run()
		return
	}

	ensureAgentRunning()
	app := app.New()

	rootCmd := cmd.NewRootCmd(app)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func ensureAgentRunning() {
	_, _, err := agent.Status()
	if err == nil {
		return
	}

	// start agent as a detached background process
	cmd := exec.Command(os.Args[0], "agent")
	if err := cmd.Start(); err != nil {
		log.Fatal("failed to start agent:", err)
	}

	// wait briefly (avoid race)
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		if _, _, err := agent.Status(); err == nil {
			return
		}
	}
	log.Fatal("failed to start agent")
}
