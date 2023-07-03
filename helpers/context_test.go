package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClabFN(t *testing.T) {
	assert.Equal(t, "a.labctl.yml", labctlFilename("a.clab.yml"))
	assert.Equal(t, "a.labctl.yaml", labctlFilename("a.clab.yaml"))
	assert.Equal(t, "a.labctl.yaml", labctlFilename("a.yaml"))
	assert.Equal(t, "a.labctl.yml", labctlFilename("a.yml"))
}
