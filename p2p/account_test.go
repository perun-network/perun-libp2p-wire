package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pkgtest "polycry.pt/poly-go/test"
)

func TestNewAccount(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
}

func getHost(t *testing.T) *Account {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	return acc
}
