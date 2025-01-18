package emitter

import (
	"os"
	"testing"

	"github.com/nexi-intra/koksmat-emit/config"
)

func TestMain(m *testing.M) {
	config.Setup()

	code := m.Run()

	os.Exit(code)
}
