package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPkgManager(t *testing.T) {
	assert := assert.New(t)
	mgrTypes := []string{"apt", "yum", "apk", "pip", "gem"}
	for _, mgrType := range mgrTypes {
		m, err := NewManager(mgrType)
		assert.NoError(err)
		assert.Equal(mgrType, m.Type())
	}
	// Unsupported package manager
	_, err := NewManager("random")
	assert.Error(err)
}
