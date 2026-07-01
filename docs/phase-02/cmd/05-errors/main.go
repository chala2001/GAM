package main

import (
	"errors"
	"fmt"
)

func validateAPIName(name string) error {
	if name == "" {
		return errors.New("api name cannot be empty")
	}
	return nil
}

func main() {
	names := []string{"orders-api", ""}

	for _, name := range names {
		if err := validateAPIName(name); err != nil {
			fmt.Println("rejected:", err)
			continue
		}
		fmt.Println("accepted:", name)
	}
}
