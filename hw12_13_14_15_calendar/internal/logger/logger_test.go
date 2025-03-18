package logger

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("new logger", func(t *testing.T) {
		var buf bytes.Buffer
		r, w, err := os.Pipe()
		require.NoError(t, err)

		originalStdout := os.Stdout
		os.Stdout = w
		defer func() {
			os.Stdout = originalStdout

			_ = w.Close()

			err = r.Close()
			require.NoError(t, err)
		}()

		require.NotPanics(t, func() {
			logger := New("INFO")
			logger.Info("test")
			logger.Warn("test")
			logger.Error("test")
		})

		_ = w.Close()

		_, err = io.Copy(&buf, r)
		require.NoError(t, err)

		output := buf.String()
		require.Contains(t, output, `"severity":"INFO"`)
		require.Contains(t, output, `"severity":"WARN"`)
		require.Contains(t, output, `"severity":"ERROR"`)
	})
}
