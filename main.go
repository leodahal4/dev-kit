package main

import (
	"github.com/leodahal4/dev-kit/cmd"
	"github.com/sirupsen/logrus"
)

type SimpleFormatter struct{}

func (f *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil // Only output the message
}

func main() {
	logrus.SetFormatter(&SimpleFormatter{})
	cmd.Execute()
}
