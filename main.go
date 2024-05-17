package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var logLevel string

func main() {
	var rootCmd = &cobra.Command{
		Use: "gitconfig-cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logLevel, _ = cmd.Flags().GetString("log")
		},
	}

	var setupGitsignCmd = &cobra.Command{
		Use:   "setup-gitsign",
		Short: "Setup git signing configuration",
		Run: func(cmd *cobra.Command, args []string) {
			setupGitsign()
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "", "Set log level")
	rootCmd.AddCommand(setupGitsignCmd)
	rootCmd.Execute()
}

func setupGitsign() {
	err := checkDirectoryForGit()
	if err != nil {
		logError(err)
		return
	}
	err = setupGitLocalCleanup()
	if err != nil {
		logError(err)
		return
	}
	err = setupGitLocalUser()
	if err != nil {
		logError(err)
		return
	}
}

func checkDirectoryForGit() error {
	_, err := os.Stat(".git")
	if os.IsNotExist(err) {
		return fmt.Errorf("this directory does not contain a git repository")
	}
	return nil
}

func setupGitLocalCleanup() error {
	configs := []string{
		"user.email",
		"commit.gpgsign",
		"tag.gpgsign",
		"gpg.x509.program",
		"gpg.format",
		"gitsign.connectorID",
		"user.signingkey",
		"gpg.ssh.program",
	}

	for _, config := range configs {
		err := runGitUnset(config)
		if err != nil {
			if !strings.Contains(fmt.Sprintf("%s", err), "exit status 5") {
				return fmt.Errorf("%s", err)
			}
		}
	}

	return nil
}

func setupGitLocalUser() error {
	fmt.Println("Enter which email address to use:")
	fmt.Println("1) lewis@denhamparry.co.uk")
	fmt.Println("2) lewis@chainguard.dev")
	fmt.Println("0) clear")

	reader := bufio.NewReader(os.Stdin)
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	switch email {
	case "1":
		err := setupGitConfigLocalUser("lewis@denhamparry.co.uk")
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	case "2":
		err := setupGitConfigLocalUser("lewis@chainguard.dev")
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	case "0":

	default:
		fmt.Fprintln(os.Stdout, []any{"Invalid option"}...)
	}
	return nil
}

func setupGitConfigLocalUser(email string) error {
	configurations := map[string]string{
		"commit.gpgsign":      "true",
		"tag.gpgsign":         "true",
		"gpg.x509.program":    "gitsign",
		"gpg.format":          "x509",
		"gitsign.connectorID": "https://accounts.google.com",
	}

	configurations["user.email"] = email

	for key, value := range configurations {
		err := runGitConfig(key, value)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}

	return nil
}

func runGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--local", key, value)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run git config for %s: %v", key, err)
	}
	return nil
}

func runGitUnset(key string) error {
	cmd := exec.Command("git", "config", "--local", "--unset", key)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if !strings.Contains(fmt.Sprintf("%s", err), "exit status 5") {
			return fmt.Errorf("failed to unset git config key %s: %v", key, err)
		}
	}
	return nil
}

func logError(err error) {
	if logLevel == "error" {
		fmt.Println(err)
	}
}
