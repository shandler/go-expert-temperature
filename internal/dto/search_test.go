package dto_test

import (
	"testing"

	"github.com/shandler/go-expert-temperature/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Run("Test SearchRequest", func(t *testing.T) {
		request := dto.SearchRequest{ZipCode: "07987110"}

		assert.Equal(t, "07987110", request.ZipCode)
	})

	t.Run("Test SearchResponse", func(t *testing.T) {
		response := dto.SearchResponse{Status: 200, Body: "Test Body"}

		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "Test Body", response.Body)
	})
}
