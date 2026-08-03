package protocol


import (
	"strings"
)

func Parse(input string) Command {
	cleanInput:= strings.TrimSpace(input)
	parts:= strings.Fields(cleanInput)

	if len(parts) == 0 {
		return Command{
			Action:	"UNKNOWN",
			Args:	[]string{},
		}
	}
	action := strings.ToUpper(parts[0])
	args := parts[1:]

	return Command{
		Action:	action,
		Args:	args,
	}
}
