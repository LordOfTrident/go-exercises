package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"bufio"
	"regexp"
	"strconv"
)

const (
	defaultPrompt = "$ "
	rcFile        = "./shellrc"
)

var in *bufio.Reader

func escape(p_str string) (string, error) {
	result := ""

	escape := false
	for i := 0; i < len(p_str); i ++ {
		if escape {
			switch p_str[i] {
			case 'n':  result += "\n";   break
			case 'r':  result += "\r";   break
			case 't':  result += "\t";   break
			case 'e':  result += "\x1b"; break
			case '\\': result += "\\";   break
			case '$':  result += "\\$";  break
			case 's':  result += " ";    break

			default: return "", fmt.Errorf("Unknown escape sequence '\\%v'", string(p_str[i]))
			}

			escape = false
		} else {
			switch p_str[i] {
			case '\\': escape = true; break

			default: result += string(p_str[i])
			}
		}
	}

	return result, nil
}

func escapeCmd(p_cmd []string) error {
	for i, _ := range p_cmd {
		var err error
		p_cmd[i], err = escape(p_cmd[i])
		if err != nil {
			return err
		}

		if len(p_cmd[i]) < 2 {
			continue
		}

		if p_cmd[i][0] == '$' {
			p_cmd[i] = os.Getenv(p_cmd[i][1:])
		} else if p_cmd[i][0:2] == "\\$" {
			p_cmd[i] = p_cmd[i][1:]
		}
	}

	return nil
}

func evalCmd(p_cmd string) {
	re  := regexp.MustCompile("\\s+")
	cmd := re.Split(p_cmd, -1)

	if err := escapeCmd(cmd[1:]); err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}

	if len(cmd[0]) == 0 {
		return
	}

	switch {
	case cmd[0] == "halt":
		if len(cmd) <= 1 {
			os.Exit(0)
		} else if len(cmd) > 2 {
			fmt.Printf("Error: Unexpected argument '%v' in halt\n", cmd[2])

			return
		}

		num, err := strconv.Atoi(cmd[1])
		if err != nil {
			fmt.Printf("Error: Expected integer as argument to halt, got '%v'\n", cmd[1])

			return
		}

		os.Exit(num)


	case cmd[0][len(cmd[0]) - 1] == ':':
		name, newVal := cmd[0][:len(cmd[0]) - 1], ""

		for i, val := range cmd[1:] {
			if i > 0 {
				newVal += " "
			}

			newVal += val
		}

		os.Setenv(name, newVal)

	default:
		process := exec.Command(cmd[0], cmd[1:]...)

		out, err := process.Output()
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
		} else {
			fmt.Print(string(out))
		}
	}
}

func escapePrompt(p_prompt string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return strings.Replace(p_prompt, "\\w", wd, -1)
}

func prompt() {
	quit := false
		for !quit {
		fmt.Print(escapePrompt(os.Getenv("PS1")))

		input, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if input[len(input) - 1] == '\n' {
			input = input[:len(input) - 1]
		}

		if input == "exit" {
			quit = true
		} else {
			evalCmd(input)
		}
	}
}

func evalFile(p_path string) error {
	file, err := os.Open(p_path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		evalCmd(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Printf("Usage: %v [SCRIPT...]\n", os.Args[0])
}

func init() {
	os.Setenv("SHELL", os.Args[0])
	os.Setenv("PS1", defaultPrompt)
}

func main() {
	if err := evalFile(rcFile); err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		fmt.Printf("Creating '%v'\n", rcFile)

		file, err := os.Create(rcFile)
		if err != nil {
			panic(err)
		}

		file.WriteString("PS1: \\e[0m\\e[94m\\\\w\\e[0m\\e[92m $ \\e[0m")
		file.Close()
	}

	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-h" {
				usage()

				return
			}

			if err := evalFile(arg); err != nil {
				fmt.Printf("Error: %v\n", err.Error())
			}
		}
	} else {
		in = bufio.NewReader(os.Stdin)
		prompt()
	}
}
