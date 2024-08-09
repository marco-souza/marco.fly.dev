package lua

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/Shopify/go-lua"
	"github.com/marco-souza/marco.fly.dev/internal/binance"
	"github.com/marco-souza/marco.fly.dev/internal/currency"
	"github.com/marco-souza/marco.fly.dev/internal/discord"
	"github.com/marco-souza/marco.fly.dev/internal/telegram"
)

var logger = slog.With("lua:")

type LuaService struct {
	*lua.State
	outputLock *sync.Mutex
}

func NewLuaService() *LuaService {
	state := lua.NewState()

	lua.OpenLibraries(state)

	return &LuaService{
		State:      state,
		outputLock: &sync.Mutex{},
	}
}

func (ls *LuaService) Run(snippet string) (string, error) {
	// running lua snippet
	ls.outputLock.Lock()
	defer ls.outputLock.Unlock()

	outputReader, outputWriter, _ := os.Pipe()
	rescueStdout := os.Stdout // save the actual stdout
	os.Stdout = outputWriter  //  redirect stdout to pipe

	err := lua.DoString(ls.State, snippet)
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

func (ls *LuaService) Start() error {
	// add discord to Runtime
	discord.PushClient(ls.State)
	telegram.PushClient(ls.State)
	currency.PushClient(ls.State)
	binance.PushClient(ls.State)

	return nil
}

func (ls *LuaService) Stop() error {
	return nil
}
