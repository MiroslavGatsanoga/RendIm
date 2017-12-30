package main

import (
	"RendIm/rendim"
	"fmt"
	"image/png"
	"os"
	"time"
)

func main() {
	start := time.Now()
	img := rendim.Render()
	elapsed := time.Since(start)
	fmt.Println("Image rendering took:", elapsed)

	f, err := os.Create("out.png")
	defer f.Close()
	if err != nil {
		panic("cannot create out.png")
	}

	png.Encode(f, img)
}
