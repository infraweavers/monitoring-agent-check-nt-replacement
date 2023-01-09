package nagios

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThis(t *testing.T) {
	t.Run("Test This Baby", func(t *testing.T) {
		parsedThing := ParseRangeString("10")
		assert.Equal(t, parsedThing.End, 10.0)

		isInRange := parsedThing.CheckRange("54")
		assert.Equal(t, isInRange, true)
	})

	t.Run("Within a defined range", func(t *testing.T) {
		parsedThing := ParseRangeString("10:20")
		assert.Equal(t, parsedThing.Start, 10.0)
		assert.Equal(t, parsedThing.End, 20.0)

		isInRange := parsedThing.CheckRange("54")
		assert.Equal(t, isInRange, false)
	})

}
