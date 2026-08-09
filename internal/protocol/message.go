package protocol

import "fmt"

type Command struct {
	Action string
	Args   []string
	// Raw preserves multi-word arguments that Args would split apart.
	Raw string
}

func FormatOK(data string) string {
	if data == "" {
		return "OK\n"
	}
	return "OK " + data + "\n"
}

func FormatErr(code int, message string) string {
	return fmt.Sprintf("ERR %d %s\n", code, message)
}

func FormatEvt(category, evtType, data string) string {
	if data == "" {
		return fmt.Sprintf("EVT %s %s\n", category, evtType)
	}
	return fmt.Sprintf("EVT %s %s %s\n", category, evtType, data)
}
