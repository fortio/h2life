// Serve game of life over http (1.1 or better, 2.0 streaming).
// Use to `curl http://localhost:31337` to see the output as it is streamed.
// (or fortio's `h2cli -stream`)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"fortio.org/fortio/fhttp"
	"fortio.org/log"
	"fortio.org/progressbar"
	"fortio.org/scli"
	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/life/conway"
)

func main() {
	os.Exit(Main())
}

var (
	delayFlag   = flag.Duration("delay", 100*time.Millisecond, "Delay between frames")
	maxIterFlag = flag.Int("iter", 79, "Number of iterations per request (in addition to the initial)")
)

func Main() int {
	portFlag := flag.String("port", ":31337", "Port to listen on")
	scli.ServerMain()
	mux, _ := fhttp.HTTPServer("life", *portFlag)
	mux.HandleFunc("GET /life", log.LogAndCall("life",
		fhttp.Gzip(http.HandlerFunc(lifeHandler)).ServeHTTP))
	scli.UntilInterrupted()
	return 0
}

func isRealBrowser(userAgent string) bool {
	return strings.Contains(userAgent, "Mozilla")
}

func lifeHandler(w http.ResponseWriter, r *http.Request) {
	if isRealBrowser(r.UserAgent()) {
		http.Redirect(w, r, "https://github.com/fortio/h2life", http.StatusFound)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	// chunked is implied by multiple writes/flushes and no content-length
	w.WriteHeader(http.StatusOK)
	ww := bufio.NewWriter(w)
	game := &conway.Game{}
	game.AP = &ansipixels.AnsiPixels{Out: ww}
	game.AP.W = 80
	game.AP.H = 24
	game.C = conway.NewConway(game.AP.W, 2*game.AP.H)
	game.C.Randomize(0.1)
	cfg := progressbar.DefaultConfig()
	cfg.ScreenWriter = ww
	cfg.UpdateInterval = 0
	cfg.Width = 20 // smaller so there is room for the life game on last line too.
	pbar := cfg.NewBar()
	maxIter := *maxIterFlag
	game.Extra = func() {
		game.AP.MoveCursor(0, game.AP.H-1)
		pbar.Progress(float64(100*game.Generation) / float64(maxIter+1))
	}
	pbar.Extra = func(_ *progressbar.Bar, _ float64) string {
		return fmt.Sprintf(", Generation: %d ", game.Generation)
	}
	game.Start()
	delay := *delayFlag
	// fmt.Fprintln(ww, "Starting...")
	for i := range maxIter {
		// Add go up 1 line + \n
		// curl needs to see a \n to flush without --no-buffer
		_, _ = w.Write([]byte("\x1b[1A\n"))
		flusher.Flush()
		select {
		case <-r.Context().Done():
			log.LogVf("Client disconnected")
			return
		case <-time.After(delay):
			// fmt.Fprintln(ww, i)
			log.LogVf("Iteration %d", i)
			game.Next()
		}
	}
	_, _ = w.Write([]byte("\r\n\n"))
}
