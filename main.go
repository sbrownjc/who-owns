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

	sshAuth, err := ssh.NewSSHAgentAuth("git")
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
		URL:   url,
		Auth:  sshAuth,
		Depth: 1,
		Tags:  git.NoTags,
	})
	CheckIfError(err)

	// Retrieve the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// Retrieve the head object
	head, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	Info("Inventory last updated on %s", head.Committer.When.Local().Format(time.RFC1123))

	inventoryYaml, err := os.ReadFile(directory + "/inventory.yaml")
	CheckIfError(err)

	var inventory Inventory
	err = yaml.Unmarshal(inventoryYaml, &inventory)
	CheckIfError(err)

	fmt.Println()
	for owner, ownership := range inventory {
		for _, namespace := range ownership.Namespaces {
			if strings.Contains(namespace, name) {
				fmt.Printf("Found k8s namespace %q owned by %s\n", namespace, owner)
			}
		}
		for _, repo := range ownership.Repos {
			if strings.Contains(repo, name) {
				fmt.Printf("Found repo %s owned by %s\n", repo, owner)
			}
		}
	}
}
