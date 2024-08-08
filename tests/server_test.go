package main

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateServer(t *testing.T) {
	t.Setenv("VIEWS", "../views/")

	done := make(chan bool)
	s := server.New(&done)
	s.Start()

	go s.Run()

	defer s.Stop()
	values := fmt.Sprintf("%v", s)
	baseUrl := "http://localhost:3001"

	t.Run("validate env configs", func(t *testing.T) {
		assert.NotNil(t, s)
		assert.Contains(t, values, "localhost")
		assert.Contains(t, values, "3001")
	})

	t.Run("can start server", func(t *testing.T) {
		assert.True(t, <-done)
	})

	t.Run("can fetch home /", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseUrl+"/", nil)

		resp, err := s.Test(req, 5000)
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, 200)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body), "Marco")
	})

	t.Run("can fetch home /resume", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseUrl+"/resume", nil)

		resp, err := s.Test(req, 5000)
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, 200)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body), "Marco Antônio")
	})

	t.Run("can shutdown server", func(t *testing.T) {
		s.Stop()
		fmt.Println("server stopped")
	})
}

// req := httptest.NewRequest("GET", "http://google.com", nil)
// req.Header.Set("X-Custom-Header", "hi")
