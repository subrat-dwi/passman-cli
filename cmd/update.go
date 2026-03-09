package cmd

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
	"github.com/subrat-dwi/passman-cli/internal/service"
	"github.com/subrat-dwi/passman-cli/internal/ui/styles"
)

const (
	repoOwner = "subrat-dwi"
	repoName  = "passman-cli"
)

var checkOnly bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update pman to the latest version",
	Long:  "Check for updates and optionally download and install the latest version of pman.",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "Only check for updates without installing")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := service.Version()

	// Strip 'v' prefix if present for comparison
	cleanVersion := strings.TrimPrefix(currentVersion, "v")

	fmt.Printf("%s %s\n", styles.Dim.Render("Current version:"), styles.Bold.Render(currentVersion))
	fmt.Println(styles.Dim.Render("Checking for updates..."))

	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return fmt.Errorf("failed to create update source: %w", err)
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source:    source,
		Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"},
	})
	if err != nil {
		return fmt.Errorf("failed to create updater: %w", err)
	}

	latest, found, err := updater.DetectLatest(context.Background(), selfupdate.NewRepositorySlug(repoOwner, repoName))
	if err != nil {
		return fmt.Errorf("failed to detect latest version: %w", err)
	}

	if !found {
		fmt.Println(styles.WarningMsg("No releases found."))
		return nil
	}

	fmt.Printf("%s %s\n", styles.Dim.Render("Latest version: "), styles.Info.Render(latest.Version()))

	if !latest.GreaterThan(cleanVersion) {
		fmt.Println(styles.SuccessMsg("You are already running the latest version!"))
		return nil
	}

	fmt.Printf("\n%s %s %s %s\n",
		styles.Highlight.Render("New version available:"),
		styles.Dim.Render(currentVersion),
		styles.Dim.Render("→"),
		styles.Info.Render(latest.Version()))
	fmt.Printf("%s https://github.com/%s/%s/releases/tag/v%s\n",
		styles.Dim.Render("Release notes:"),
		repoOwner, repoName, latest.Version())

	if checkOnly {
		fmt.Printf("\n%s\n", styles.InfoMsg("Run 'pman update' (without --check) to install the update."))
		return nil
	}

	fmt.Printf("\n%s %s for %s/%s...\n",
		styles.Dim.Render("Downloading"),
		styles.Info.Render(latest.Version()),
		runtime.GOOS, runtime.GOARCH)

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	if err := updater.UpdateTo(context.Background(), latest, exe); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	fmt.Printf("\n%s\n", styles.SuccessMsg(fmt.Sprintf("Successfully updated to version %s!", latest.Version())))
	return nil
}
