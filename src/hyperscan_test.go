//go:build chimera
// +build chimera

package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hs_makeAuditLogQueryRegex(t *testing.T) {
	assert.NotPanics(t, func() { hs_makeAuditLogQueryRegex([]string{"db1", "db2"}) })
}
