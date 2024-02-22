package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough args")
	}

	scanner := bufio.NewScanner(os.Stdin)

	args := os.Args[1:]
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}
		cmd := exec.Command(args[0], append(args[1:], input)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			log.Println(cmd, "Invalid command")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
