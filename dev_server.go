package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func isPortOpen(port string) bool {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return false // Port is in use
	}
	ln.Close()
	return true // Port is available
}

func runDevServer() {
	staticDir := "./build"
	port := ":8080"

	if !isPortOpen(port) {
		fmt.Println("port already in use.")
		return
	}

	fs := http.FileServer(http.Dir(staticDir))
	noCacheHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Surrogate-Control", "no-store")

		path := staticDir + "/" + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// If the file doesn't exist, serve index.html
			http.ServeFile(w, r, path+".html")
			return
		}

		fs.ServeHTTP(w, r)
	})

	http.Handle("/", http.StripPrefix("/", noCacheHandler))

	fmt.Printf("serving static files from %s on http://localhost%s\n", staticDir, port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
