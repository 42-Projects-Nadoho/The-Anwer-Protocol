package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

const addr = ":8081"

func main() {
	http.Handle("/", noCache(http.FileServer(http.Dir("cmd/gui/web"))))
	http.Handle("/assets/", noCache(http.StripPrefix("/assets/", http.FileServer(http.Dir("data")))))

	fmt.Println("Listening on " + addr)
	go openBrowser("http://localhost" + addr)
	http.ListenAndServe(addr, nil)
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
