package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPath(t *testing.T) {
	pattern := Get.WithPath("/acme")

	assert.Equal(t, "GET /acme", pattern)
}

func TestWithPathParts(t *testing.T) {
	pattern := Put.WithPath("/prefix/", "/acme//", "/{id}")

	assert.Equal(t, "PUT /prefix/acme/{id}", pattern)
}
