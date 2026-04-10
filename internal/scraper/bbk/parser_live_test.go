package bbk

import (
	"testing"

	"github.com/FaintLocket424/rc-timing-api/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestParseRaceResult(t *testing.T) {
	raw := "HTML HERE"
	result := parseRaceResult(raw)
	expected := models.LiveTiming{}

	assert.Equal(t, expected, result, "test failed lol")
}
