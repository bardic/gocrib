//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ServerUp() error {
	fmt.Println("Starting Server...")
	cmd := exec.Command("dagger", "call", "server", "--src=.", "--migrate", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func DbUp() error {
	fmt.Println("Starting Server...")
	cmd := exec.Command("dagger", "call", "db-up", "--src=.", "--with-port", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func GenQueries() error {
	fmt.Println("Building Queries...")
	cmd := exec.Command("dagger", "call", "build-queries", "--src=.", "export", "--path=sql/queries")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test() error {
	fmt.Println("Testing...")
	cmd := exec.Command("dagger", "call", "test-http", "--src=.", "--is-ci")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func LocalTest() error {
	fmt.Println("Local Testing...")
	os.Setenv("GOCRIB_HOST", "localhost")

	cmd := exec.Command("migrate", "-database", "postgres://postgres:example@localhost:5432/cribbage?sslmode=disable", "-path", "sql/migrations", "down", "-all")
	o, err := cmd.CombinedOutput()
	fmt.Println(string(o))

	cmd = exec.Command("migrate", "-database", "postgres://postgres:example@localhost:5432/cribbage?sslmode=disable", "-path", "sql/migrations", "up")
	o, err = cmd.CombinedOutput()
	fmt.Println(string(o))

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

	f = append(f, "-e", "local", "-v", "./http/http-client.env.json")

	cmd = exec.Command("ijhttp", f...)
	o, err = cmd.CombinedOutput()
	fmt.Println(string(o))

	return err
}

func MigrateLocal() error {
	fmt.Println("Migrating...")
	cmd := exec.Command("migrate", "-database", "postgres://postgres:example@localhost:5432/cribbage?sslmode=disable", "-path", "sql/migrations", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BuildCLI() error {
	fmt.Println("Building CLI...")
	cmd := exec.Command("dagger", "call", "test-http", "--src=.")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BuildServer() error {
	fmt.Println("Building Server...")
	cmd := exec.Command("dagger", "call", "build-server", "--src=.", "export", "--path=builds/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BuildGame() error {
	fmt.Println("Building Game...")
	cmd := exec.Command("dagger", "call", "build-game", "--src=.", "export", "--path=builds/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
