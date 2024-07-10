// ref: https://github.com/gofiber/recipes/blob/master/sse/main.go
package api

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func sseHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		logger.Info("WRITER")
		var i int
		for {
			i++
			msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)

			logger.Info(msg)

			err := w.Flush()
			if err != nil {
				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				logger.Error("error while flushing, closing http connection", "err", err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}

var shouldReload = true

func sseReloadHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		logger.Info("connection created")
		for {
			if shouldReload {
				logger.Info("reload", "shouldReload", shouldReload)

				// send message if should reload
				msg := fmt.Sprintf("reload")
				fmt.Fprintf(w, "data: Message: %s\n\n", msg)
				shouldReload = false
			}

			err := w.Flush()
			if err != nil {
				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				logger.Info("Error while flushing, closing http connection", "err", err)
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
	}))

	return nil
}
