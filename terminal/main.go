package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Hope it's a normal file
	if info.IsDir() {
		return false
	}
	// permissions checking
	perms := info.Mode()
	return perms&0111 != 0 // Checking if any exec bits is set

}

func findInPath(filename string) (string, error) {
	// Fetch Path var
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH environment variable is empty")
	}

	// split path env bt os-specific seperators (: or ;)
	directories := filepath.SplitList(pathEnv)
	for _, dir := range directories {
		fullPath := filepath.Join(dir, filename)

		// existence and permissions check
		if isExecutable(fullPath) {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("file %s not found in PATH", filename)

}

func main() {

mainLoop:
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("$ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		input = strings.TrimSpace(input)
		command := strings.Fields(input)

		if len(command) == 0 {
			continue
		}
		builtIns := []string{"echo", "exit", "type"}

		switch command[0] {
		case "type":
			if len(command) == 2 {
				match := false
				for _, builtIn := range builtIns {
					if command[1] == builtIn {
						fmt.Println(builtIn, "is a shell builtin")
						match = true
						break
					}
				}
				if !match {
					// Path variable exec
					// Get/collect paths
					// Find find in paths[]
					// Check file permission if exists
					result, err := findInPath(command[1])
					if err != nil {
						fmt.Printf("%s: not found \n", command[1])
					} else {
						fmt.Printf("%s is %s\n", command[1], result)
					}
				}
			} else {
				fmt.Println(input, ": command not found")
			}
		case "exit":
			break mainLoop
		case "echo":
			remainingWords := strings.Join(command[1:], " ")
			fmt.Println(remainingWords)
		default:
			fmt.Printf("%s: command not found \n", input)
		}
	}
}
