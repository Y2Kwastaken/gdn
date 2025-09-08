package internal

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv() error {
	return LoadEnvFrom(".env")
}

func LoadEnvFrom(environment string) error {
	file, err := os.Open(environment)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := strings.Split(scanner.Text(), "=")
		os.Setenv(text[0], text[1])
	}

	return nil
}
