package config

import (
	"testing"
	"time"

	"github.com/lex-r/dwl/downloader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Parse(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		c := NewConfig()
		require.NoError(t, c.Parse([]string{}))

		assert.Equal(t, c.Threads, downloader.DefaultThreadOpt)
		assert.Equal(t, c.Timeout, downloader.DefaultTimeoutOpt)
		assert.Equal(t, c.Links, []string{})
	})

	t.Run("threads option", func(t *testing.T) {
		c := NewConfig()
		require.NoError(t, c.Parse([]string{"-threads", "1"}))
		assert.Equal(t, 1, c.Threads)
	})

	t.Run("timeout option", func(t *testing.T) {
		c := NewConfig()
		require.NoError(t, c.Parse([]string{"-timeout", "1s"}))
		assert.Equal(t, time.Second, c.Timeout)
	})
}

func TestConfig_ParseInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			"invalid threads option",
			[]string{"-threads", "0"},
		},
		{
			"invalid threads option",
			[]string{"-threads", "-1"},
		},
		{
			"invalid timeout option",
			[]string{"-timeout", "0"},
		},
		{
			"invalid timeout option",
			[]string{"-timeout", "-1s"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig()
			err := c.Parse(tt.args)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid value")
		})
	}
}

func TestConfig_ParseLinks(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		links []string
	}{
		{
			"empty args",
			[]string{},
			[]string{},
		},
		{
			"one link",
			[]string{"localhost"},
			[]string{"localhost"},
		},
		{
			"two links",
			[]string{"http://localhost", "https://example.com"},
			[]string{"http://localhost", "https://example.com"},
		},
		{
			"links after other args",
			[]string{"-threads", "1", "-timeout", "5s", "http://localhost", "https://example.com"},
			[]string{"http://localhost", "https://example.com"},
		},
		{
			"args without links",
			[]string{"-threads", "1", "-timeout", "5s"},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig()
			require.NoError(t, c.Parse(tt.args))

			assert.Equal(t, tt.links, c.Links)
		})
	}
}
