package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Mirrors the LOOK JSON shape from internal/server/handlers.go.
type lookPayload struct {
	Room struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Exits       map[string]string `json:"exits"`
	} `json:"room"`
	Players []string `json:"players"`
}

func main() {
	a := app.New()
	w := a.NewWindow("TAP - The Answer Protocol")
	w.Resize(fyne.NewSize(700, 500))
	showConnectForm(a, w)
	w.ShowAndRun()
}

func showConnectForm(a fyne.App, w fyne.Window) {
	addr := widget.NewEntry()
	addr.SetText("localhost:8080")
	user := widget.NewEntry()
	user.SetPlaceHolder("username")
	status := widget.NewLabel("")

	var connectBtn *widget.Button
	connectBtn = widget.NewButton("Connect", func() {
		if strings.TrimSpace(user.Text) == "" {
			status.SetText("Username is required.")
			return
		}
		connectBtn.Disable()
		status.SetText("Connecting...")
		go connect(w, addr.Text, user.Text, func(errMsg string) {
			status.SetText(errMsg)
			connectBtn.Enable()
		})
	})

	w.SetContent(container.NewCenter(container.NewVBox(
		widget.NewLabel("Server address"), addr,
		widget.NewLabel("Username"), user,
		connectBtn, status,
	)))
}

// Runs off the main goroutine, so every UI touch goes through fyne.Do.
func connect(w fyne.Window, addr, username string, onFail func(string)) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fyne.Do(func() { onFail("Connection failed: " + err.Error()) })
		return
	}

	roomLabel := widget.NewLabel("")
	roomLabel.Wrapping = fyne.TextWrapWord
	exitsBox := container.NewHBox()
	logBox := widget.NewMultiLineEntry()
	logBox.Disable()
	logBox.Wrapping = fyne.TextWrapWord

	appendLog := func(line string) {
		fyne.Do(func() { logBox.SetText(logBox.Text + line + "\n") })
	}
	send := func(cmd string) { fmt.Fprintf(conn, "%s\n", cmd) }

	lookBtn := widget.NewButton("LOOK", func() { send("LOOK") })
	quitBtn := widget.NewButton("QUIT", func() { send("QUIT") })

	gameShown := false

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "OK hello proto=1":
			send("CONNECT " + username)

		case line == "OK connected":
			fyne.Do(func() {
				w.SetContent(container.NewBorder(
					container.NewVBox(roomLabel, exitsBox, container.NewHBox(lookBtn, quitBtn)),
					nil, nil, nil,
					container.NewScroll(logBox),
				))
			})
			gameShown = true
			send("LOOK")

		case strings.HasPrefix(line, "ERR 201 NAME_IN_USE"):
			conn.Close()
			fyne.Do(func() { onFail("Username already taken.") })
			return

		case strings.HasPrefix(line, "OK {"):
			var payload lookPayload
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "OK ")), &payload); err == nil {
				fyne.Do(func() {
					roomLabel.SetText(fmt.Sprintf(
						"%s\n%s\nPlayers here: %s",
						payload.Room.Name, payload.Room.Description,
						strings.Join(payload.Players, ", "),
					))
					exitsBox.RemoveAll()
					for direction := range payload.Room.Exits {
						direction := direction
						exitsBox.Add(widget.NewButton(direction, func() { send("MOVE " + direction) }))
					}
					exitsBox.Refresh()
				})
			}

		case strings.HasPrefix(line, "OK room="):
			send("LOOK")

		case line == "OK bye":
			conn.Close()

		default:
			appendLog(line)
		}
	}

	if gameShown {
		appendLog("--- disconnected ---")
	} else {
		fyne.Do(func() { onFail("Disconnected.") })
	}
}
