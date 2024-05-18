package main

import (
	"log"
	"os"

	"github.com/Shopify/go-lua"
)

func main() {
	l := lua.NewState()
	lua.OpenLibraries(l)

	// run snippet
	lua.DoString(l, "print('Hello from Lua!')")

	if len(os.Args) == 1 {
		log.Println("No file to run")
		return
	}

	filename := os.Args[1]
	if err := lua.DoFile(l, filename); err != nil {
		panic(err)
	}
}
