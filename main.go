package main

import (
	"RendIm/rendim"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	width  = 800
	height = 800
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	indexFile, err := os.Open("html/index.html")
	if err != nil {
		fmt.Println(err)
	}
	index, err := io.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		sceneType := r.URL.Query().Get("scene")
		if sceneType == "" {
			sceneType = "final"
		}

		samples := 10000
		if s := r.URL.Query().Get("samples"); s != "" {
			_, _ = fmt.Sscanf(s, "%d", &samples)
		}

		bucketSize := 32
		if bs := r.URL.Query().Get("bucketSize"); bs != "" {
			_, _ = fmt.Sscanf(bs, "%d", &bucketSize)
		}

		workers := 4
		if w := r.URL.Query().Get("workers"); w != "" {
			_, _ = fmt.Sscanf(w, "%d", &workers)
		}

		fmt.Printf("Client initiated a render (scene: %s, samples: %d, bucketSize: %d, workers: %d)...\n",
			sceneType, samples, bucketSize, workers)

		pixels := make(chan rendim.Pixel, width*height)
		go rendim.RenderScene(width, height, sceneType, samples, bucketSize, workers, pixels)

		for {
			time.Sleep(2 * time.Second)

			for p := range pixels {
				data, err := json.Marshal(p)
				if err != nil {
					fmt.Println(err)
					return
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "%s", string(index)); err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println("RendIm running on port 3000.")

	server := &http.Server{
		Addr:         ":3000",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
