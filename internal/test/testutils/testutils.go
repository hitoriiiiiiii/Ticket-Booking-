// Package testutils provides utilities for testing
package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContext returns a context for testing
func TestContext() context.Context {
	return context.Background()
}

// AssertNoError fails the test if there's an error
func AssertNoError(t *testing.T, err error) {
	assert.NoError(t, err)
}

// AssertError fails the test if there's no error
func AssertError(t *testing.T, err error) {
	assert.Error(t, err)
}

// AssertEqual fails if values aren't equal
func AssertEqual[T any](t *testing.T, expected, actual T) {
	assert.Equal(t, expected, actual)
}

// SkipIfShort skips the test if -short flag is provided
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}
