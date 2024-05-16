package main

import (
	"os"
	"testing"
)

func TestCheckDirectoryForGit(t *testing.T) {
	err := checkDirectoryForGit()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSetupGitLocalCleanup(t *testing.T) {
	err := setupGitLocalCleanup()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCheckDirectoryForNoGit(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	// Change the working directory to the temporary directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(oldDir) // change back to the original directory

	// Call the function
	err = checkDirectoryForGit()
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
