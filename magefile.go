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
	os.Setenv("GOCRIB_HOST", "localhost")
	cmd := exec.Command("dagger", "call", "postgres", "--src=.", "--with-port", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func GenQueries() error {
	fmt.Println("Building Queries...")
	cmd := exec.Command("dagger", "call", "gen", "--src=.", "export", "--path=sql/queries")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test() error {
	fmt.Println("Testing...")
	cmd := exec.Command("dagger", "call", "test-http", "--src=.")
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

func MigrateLocal() error {
	fmt.Println("Migrating...")
	cmd := exec.Command("migrate", "-database", "postgres://postgres:example@localhost:5432/cribbage?sslmode=disable", "-path", "sql/migrations", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
