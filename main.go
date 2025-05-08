package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	yaml "github.com/goccy/go-yaml"
)

func main() {
	CheckArgs("<repo/namespace>")
	name := os.Args[1]

	authMethod, err := ssh.NewSSHAgentAuth("git")
	CheckIfError(err)

	// Clone the inventory-mapping repository to a temp directory
	url := "git@github.com:TheJumpCloud/inventory-mapping.git"
	directory, err := os.MkdirTemp("", "inventory-mapping")
	CheckIfError(err)
	defer func() {
		// Info("Removing temp directory %s", directory)
		if err := os.RemoveAll(directory); err != nil {
			Warning("Failed to remove temp directory %s: %v", directory, err)
		}
	}()

	// Info("Cloning inventory-mapping to %s", directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth:  authMethod,
		URL:   url,
		Depth: 1,
	})
	CheckIfError(err)

	// Retrieve the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// Retrieve the head object
	head, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	Info("Inventory last updated on %s", head.Committer.When.Local().Format(time.RFC1123))

	inventory, err := os.ReadFile(directory + "/inventory.yaml")
	CheckIfError(err)

	var inv Inventory
	err = yaml.Unmarshal(inventory, &inv)
	CheckIfError(err)

	fmt.Println()
	for owner, ownership := range inv {
		for _, ownedNS := range ownership.Namespaces {
			if strings.Contains(ownedNS, name) {
				fmt.Printf("Found k8s namespace %q owned by %s\n", ownedNS, owner)
			}
		}
		for _, ownedRepo := range ownership.Repos {
			if strings.Contains(ownedRepo, name) {
				fmt.Printf("Found repo %s owned by %s\n", ownedRepo, owner)
			}
		}
	}
}
