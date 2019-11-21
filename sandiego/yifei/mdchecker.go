package mdchecker

import (
	"bytes"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	file     = "./greetings.md"
	name     = "Yifei"
	greeting = "greeting"
)

/**
	check if the greeting MD file has the minimum requirements: name, and the greeting word
 */
func GreetingCheck() error {

	if _, err := os.Stat(file); err != nil {
		return err
	}

	reader := strings.NewReader("Greeting file")
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}
	if !bytes.ContainsAny(data, name) {
		return errors.New("Failed to find the name!")
	}

	if !bytes.ContainsAny(data, greeting) {
		return errors.New("Failed to find the greeting!")
	}

	return nil
}
