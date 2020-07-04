package scripting

import (
	"github.com/Shopify/go-lua"
	"log"
)

func HelloWorld() {
	l := lua.NewState()
	lua.OpenLibraries(l)
	if err := lua.DoFile(l, "hello.lua"); err != nil {
		log.Fatal(err)
	}
}
