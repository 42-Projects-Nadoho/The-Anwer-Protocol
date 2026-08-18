package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gorilla/websocket"
)

var (
	addr           string
	gameServerAddr string
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	httpAddr := flag.String("http", ":8081", "HTTP server address to listen on")
	gAddr := flag.String("addr", "localhost:8080", "Game server address to connect to")
	flag.Parse()

	addr = *httpAddr
	gameServerAddr = *gAddr

	http.Handle("/", noCache(http.FileServer(http.Dir("cmd/gui/web"))))
	http.Handle("/assets/", noCache(http.StripPrefix("/assets/", http.FileServer(http.Dir("data")))))
	http.HandleFunc("/ws", handleWS)

	browserURL := addr
	if len(addr) > 0 && addr[0] == ':' {
		browserURL = "localhost" + addr
	}
	fmt.Println("Listening on " + addr)
	go openBrowser("http://" + browserURL)
	http.ListenAndServe(addr, nil)
}

// handleWS bridges one browser WebSocket connection to one fresh TCP
// connection against the real game server. It doesn't parse or understand
// any of the RFC 42TAP traffic — it only relays lines in both directions;
// the actual game server on the other end does all the real work.
func handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer ws.Close()

	game, err := net.Dial("tcp", gameServerAddr)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("ERR 900 CONNECTION_FAILED"))
		return
	}
	defer game.Close()

	done := make(chan struct{})

	// game -> browser
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(game)
		for scanner.Scan() {
			if err := ws.WriteMessage(websocket.TextMessage, scanner.Bytes()); err != nil {
				return
			}
		}
	}()

	// browser -> game
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Fprintf(game, "%s\n", msg)
	}

	game.Close()
	<-done
}

// noCache forces the browser to always re-fetch instead of serving a
// stale cached copy of the page/assets while we're actively editing them.
func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// openBrowser best-effort launches the system's default browser. If it
// fails, the server keeps running; the user can still open the URL by hand.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Could not open browser automatically, visit " + url)
	}
}
