package main

import (
	"RendIm/rendim"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	width  = 300
	height = 200
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
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Client initiated a render...")

		pixels := make(chan rendim.Pixel, width*height)
		go rendim.Render(width, height, pixels)

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
		fmt.Fprintf(w, string(index))
	})

	fmt.Println("RendIm running on port 3000.")

	http.ListenAndServe(":3000", nil)
}
