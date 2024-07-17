package lua

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/Shopify/go-lua"
	"github.com/marco-souza/marco.fly.dev/internal/currency"
	"github.com/marco-souza/marco.fly.dev/internal/discord"
	"github.com/marco-souza/marco.fly.dev/internal/telegram"
)

var (
	stdoutLock = &sync.Mutex{}
	logger     = slog.With("lua:")
)

func Run(snippet string) (string, error) {
	// setup lua runtime
	l := lua.NewState()

	lua.OpenLibraries(l)
	pushRuntimeLibraries(l)

	// running lua snippet
	stdoutLock.Lock()
	defer stdoutLock.Unlock()

	outputReader, outputWriter, _ := os.Pipe()
	rescueStdout := os.Stdout // save the actual stdout
	os.Stdout = outputWriter  //  redirect stdout to pipe

	err := lua.DoString(l, snippet)
	if err != nil {
		logger.Error("error running lua snippet", "err", err)
		return "", err
	}

	outputWriter.Close()     // close pipe writer
	os.Stdout = rescueStdout // restore stdout

	output, err := io.ReadAll(outputReader)
	if err != nil {
		logger.Error("error reading lua output", "err", err)
		return "", err
	}

	return string(output), nil
}

func pushRuntimeLibraries(l *lua.State) {
	// add discord to Runtime
	discord.DiscordService.PushClient(l)
	telegram.PushClient(l)
	currency.PushClient(l)
}
