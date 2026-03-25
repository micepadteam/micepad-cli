package terminalwire

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/term"
)

func (c *Client) handleResource(msg Message) error {
	name, _ := msg["name"].(string)
	command, _ := msg["command"].(string)
	params, _ := msg["parameters"].(map[string]interface{})

	switch name {
	case "stdout":
		return c.handleIO(os.Stdout, name, command, params)
	case "stderr":
		return c.handleIO(os.Stderr, name, command, params)
	case "stdin":
		return c.handleStdin(command)
	case "browser":
		return c.handleBrowser(params)
	case "file":
		return c.handleFile(command, params)
	case "directory":
		return c.handleDirectory(command, params)
	case "environment_variable":
		return c.handleEnvVar(command, params)
	default:
		return fmt.Errorf("unknown resource: %s", name)
	}
}

// stdout / stderr — identical logic, different file descriptor
func (c *Client) handleIO(w *os.File, name, command string, params map[string]interface{}) error {
	data, _ := params["data"].(string)
	switch command {
	case "print":
		fmt.Fprint(w, data)
	case "print_line":
		fmt.Fprintln(w, data)
	}
	return c.succeed(name, nil)
}

// stdin
func (c *Client) handleStdin(command string) error {
	switch command {
	case "read_line":
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			return c.fail("stdin", err.Error())
		}
		return c.succeed("stdin", line)
	case "read_password":
		pw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println() // newline after hidden input
		if err != nil {
			return c.fail("stdin", err.Error())
		}
		return c.succeed("stdin", string(pw))
	default:
		return c.succeed("stdin", nil)
	}
}

// browser
func (c *Client) handleBrowser(params map[string]interface{}) error {
	if urlStr, ok := params["url"].(string); ok {
		openBrowser(urlStr)
	}
	return c.succeed("browser", nil)
}

// file
func (c *Client) handleFile(command string, params map[string]interface{}) error {
	pathStr := expandPath(paramStr(params, "path"))

	switch command {
	case "read":
		data, err := os.ReadFile(pathStr)
		if err != nil {
			return c.fail("file", err.Error())
		}
		return c.succeed("file", string(data))
	case "write":
		os.MkdirAll(filepath.Dir(pathStr), 0755)
		if err := os.WriteFile(pathStr, []byte(paramStr(params, "content")), 0644); err != nil {
			return c.fail("file", err.Error())
		}
		return c.succeed("file", nil)
	case "append":
		f, err := os.OpenFile(pathStr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return c.fail("file", err.Error())
		}
		defer f.Close()
		f.WriteString(paramStr(params, "content"))
		return c.succeed("file", nil)
	case "delete":
		if err := os.Remove(pathStr); err != nil {
			return c.fail("file", err.Error())
		}
		return c.succeed("file", nil)
	case "exist":
		_, err := os.Stat(pathStr)
		return c.succeed("file", err == nil)
	case "change_mode":
		mode, _ := params["mode"].(int64)
		if err := os.Chmod(pathStr, os.FileMode(mode)); err != nil {
			return c.fail("file", err.Error())
		}
		return c.succeed("file", nil)
	default:
		return c.succeed("file", nil)
	}
}

// directory
func (c *Client) handleDirectory(command string, params map[string]interface{}) error {
	pathStr := expandPath(paramStr(params, "path"))

	switch command {
	case "list":
		matches, err := filepath.Glob(pathStr)
		if err != nil {
			return c.fail("directory", err.Error())
		}
		return c.succeed("directory", matches)
	case "create":
		if err := os.MkdirAll(pathStr, 0755); err != nil {
			return c.fail("directory", err.Error())
		}
		return c.succeed("directory", nil)
	case "exist":
		info, err := os.Stat(pathStr)
		return c.succeed("directory", err == nil && info.IsDir())
	case "delete":
		if err := os.Remove(pathStr); err != nil {
			return c.fail("directory", err.Error())
		}
		return c.succeed("directory", nil)
	default:
		return c.succeed("directory", nil)
	}
}

// environment_variable
func (c *Client) handleEnvVar(command string, params map[string]interface{}) error {
	if command == "read" {
		return c.succeed("environment_variable", os.Getenv(paramStr(params, "name")))
	}
	return c.succeed("environment_variable", nil)
}

// helpers

func paramStr(params map[string]interface{}, key string) string {
	v, _ := params[key].(string)
	return v
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	}

	args = append(args, url)
	exec.Command(cmd, args...).Start()
}
