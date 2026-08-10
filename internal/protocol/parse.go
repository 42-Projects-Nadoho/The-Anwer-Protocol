package protocol

import (
	"strings"
)

func Parse(input string) Command {
	cleanInput := strings.TrimSpace(input)

	if cleanInput == "" {
		return Command{
			Action: "UNKNOWN",
			Args:   []string{},
		}
	}

	sepIdx := strings.IndexAny(cleanInput, " \t")
	if sepIdx == -1 {
		return Command{
			Action: strings.ToUpper(cleanInput),
			Args:   []string{},
		}
	}

	action := strings.ToUpper(cleanInput[:sepIdx])
	rest := strings.TrimSpace(cleanInput[sepIdx+1:])

	return Command{
		Action: action,
		Args:   strings.Fields(rest),
		Raw:    rest,
	}
}
