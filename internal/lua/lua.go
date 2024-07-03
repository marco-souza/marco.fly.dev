package lua

import (
	"io"
	"log"
	"os"

	"github.com/Shopify/go-lua"
	"github.com/marco-souza/marco.fly.dev/internal/discord"
)

func Run(snippet string) (string, error) {
	// setup lua runtime
	l := lua.NewState()

	lua.OpenLibraries(l)
	pushRuntimeLibraries(l)

	// running lua snippet
	outputReader, outputWriter, _ := os.Pipe()
	rescueStdout := os.Stdout // save the actual stdout
	// FIXME: this is a shared variable, and if causes some jobs to print using the original stdout
	//	  instead of pipeing it to the execution output. As we're running concurrently, we'd need
	//    to lock this resource and execute code snippets sync.
	os.Stdout = outputWriter //  redirect stdout to pipe

	err := lua.DoString(l, snippet)
	if err != nil {
		log.Println("error running lua snippet: ", err)
		return "", err
	}

	outputWriter.Close()     // close pipe writer
	os.Stdout = rescueStdout // restore stdout

	output, err := io.ReadAll(outputReader)
	if err != nil {
		log.Println("error reading lua output: ", err)
		return "", err
	}

	return string(output), nil
}

func pushRuntimeLibraries(l *lua.State) {
	// add discord to Runtime
	discord.DiscordService.PushClient(l)
}
