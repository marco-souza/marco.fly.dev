package lua

import (
	"io"
	"log"
	"os"

	"github.com/Shopify/go-lua"
)

type luaRuntime struct {
	l *lua.State
}

var fMap map[int]int

func fib(n int) int {
	if val, ok := fMap[n]; ok {
		return val
	}

	if n <= 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func new() *luaRuntime {
	l := lua.NewState()

	l.PushGoFunction(func(s *lua.State) int {
		n, ok := s.ToInteger(1) // pop stack first value
		if !ok {
			return 0
		}

		res := fib(n)

		s.PushInteger(res) // push result to the stack
		return 1           // number of results
	})
	l.SetGlobal("go_fib")

	lua.OpenLibraries(l)
	return &luaRuntime{l}
}

func (r *luaRuntime) Run(snippet string) (string, error) {
	log.Println("Running Lua snippet: ", snippet)

	outputReader, outputWriter, _ := os.Pipe()
	rescueStdout := os.Stdout // save the actual stdout
	os.Stdout = outputWriter  //  redirect stdout to pipe

	err := lua.DoString(r.l, snippet)
	if err != nil {
		log.Println("Error running Lua snippet", err)
		return "", err
	}

	outputWriter.Close()     // close pipe writer
	os.Stdout = rescueStdout // restore stdout

	output, err := io.ReadAll(outputReader)
	if err != nil {
		log.Println("Error reading Lua output", err)
		return "", err
	}

	log.Println("Lua output: ", string(output))
	return string(output), nil
}

var Runtime = new()
