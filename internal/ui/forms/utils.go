package forms

import (
	"errors"
	"strconv"
	"strings"
)

func validateIsNotEmpty(fieldName string) func(s string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errors.New(fieldName + " is required.")
		}
		return nil
	}
}

func validatePort(s string) error {
	if strings.TrimSpace(s) != "" {
		p, err := strconv.Atoi(s)
		if err != nil || p < 1 || p > 65535 {
			return errors.New("Incorrect port, range: 1..65535")
		}
	}
	return nil
}
