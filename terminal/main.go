package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	perms := info.Mode()
	return perms&0111 != 0 // Checking if any exec bits are set
}

func findInPath(filename string) (string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH environment variable is empty")
	}

	directories := filepath.SplitList(pathEnv)
	for _, dir := range directories {
		fullPath := filepath.Join(dir, filename)
		if isExecutable(fullPath) {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("file %s not found in PATH", filename)
}

func executeCommand(executablePath string, originalName string, args []string) error {
	// 1. Create the command using the full absolute path so the OS can find and run it
	cmd := exec.Command(executablePath, args...)

	// 2. Force Arg #0 (cmd.Args[0]) to be what the user typed (e.g., "custom_exe_2590")
	// instead of the absolute path (/tmp/fox/custom_exe_2590)
	cmd.Args = append([]string{originalName}, args...)

	// 3. Connect standard streams
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}

func main() {
	// 1. Move reader outside the loop to prevent memory allocation leaks
	reader := bufio.NewReader(os.Stdin)
	builtIns := []string{"echo", "exit", "type", "pwd"}

mainLoop:
	for {
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

		switch command[0] {
		case "type":
			if len(command) == 2 {
				target := command[1]
				match := false
				for _, builtIn := range builtIns {
					if target == builtIn {
						fmt.Printf("%s is a shell builtin\n", target)
						match = true
						break
					}
				}
				if !match {
					result, err := findInPath(target)
					if err != nil {
						fmt.Printf("%s: not found\n", target)
					} else {
						fmt.Printf("%s is %s\n", target, result)
					}
				}
			} else {
				fmt.Printf("%s: command not found\n", input)
			}
		case "exit":
			break mainLoop
		case "echo":
			remainingWords := strings.Join(command[1:], " ")
			fmt.Println(remainingWords)
		case "pwd":
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				return
			}
			fmt.Println(cwd)
		default:
			// 2. Use your custom findInPath or exec.LookPath consistently
			path, err := findInPath(command[0])
			if err != nil {
				fmt.Printf("%s: command not found\n", command[0])
				continue
			}

			// 3. Exec.Command handles args correctly whether slice is empty or has elements
			err = executeCommand(path, command[0], command[1:])
			if err != nil {
				continue
			}
		}
	}
}
