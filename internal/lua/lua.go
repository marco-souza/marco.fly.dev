package lua

import (
	"log"

	"github.com/Shopify/go-lua"
)

type luaRuntime struct {
	l *lua.State
}

func new() *luaRuntime {
	l := lua.NewState()
	lua.OpenLibraries(l)
	return &luaRuntime{l}
}

func (r *luaRuntime) Run(snippet string) (string, error) {
	log.Println("Running Lua snippet", snippet)
	// TODO: get lua execution output: https://github.com/Shopify/go-lua/pull/43
	return "", lua.DoString(r.l, snippet)
}

var Runtime = new()
