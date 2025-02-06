package utils

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
	"golang.org/x/term"
)

type Provider interface {
	ReadInput() (string, error)
	ReadPassword() (string, error)
	Clear()
}

type osCredentialProvider struct {
	userOs string
}

func New(userOs string) Provider {
	return &osCredentialProvider{userOs: userOs}
}

func (o *osCredentialProvider) ReadInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}

func (o *osCredentialProvider) ReadPassword() (string, error) {
	bytePass, err := term.ReadPassword(int(syscall.Stdin))
	return strings.TrimSpace(string(bytePass)), err
}

func (o *osCredentialProvider) Clear() {
	switch o.userOs {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	default:
		return
	}
}

func AskInput(title, placeholder string) string {
	input, _ := readline.New(title)
	defer func(input *readline.Instance) {
		_ = input.Close()
	}(input)

	data2 := []byte(placeholder)
	_, _ = input.WriteStdin(data2)

	userText, _ := input.Readline()
	return strings.TrimSpace(userText)
}
