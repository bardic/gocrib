//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// A build step that requires additional params, or platform specific steps for example
func GameBuild() error {
	fmt.Println("Building Game...")
	cmd := exec.Command("dagger", "call", "build-game", "--src=.", "export", "--path=/builds")
	return cmd.Run()
}

func ServerBuild() error {
	fmt.Println("Building Server...")
	cmd := exec.Command("dagger", "call", "build-server", "--src=.", "export", "--path=/builds")
	return cmd.Run()
}

func ServerUp() error {
	fmt.Println("Starting Server...")
	cmd := exec.Command("dagger", "call", "game-server", "--src=.", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func Queries() error {
	fmt.Println("Building Queries...")
	cmd := exec.Command("dagger", "call", "sqlc", "--src=.", "export", "--path=sql/queries")
	return cmd.Run()
}

func Test() error {
	fmt.Println("Testing...")
	cmd := exec.Command("dagger", "call", "http", "--src=.")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func LocalTest() error {
	fmt.Println("Local Testing...")

	entries, err := os.ReadDir("./http")
	if err != nil {
		return err
	}

	f := make([]string, 0)
	for _, file := range entries {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".http") {
			f = append(f, "./http/"+file.Name())
		}
	}

	cmd := exec.Command("ijhttp", f...)
	o, err := cmd.CombinedOutput()
	fmt.Println(string(o))

	return err
}
