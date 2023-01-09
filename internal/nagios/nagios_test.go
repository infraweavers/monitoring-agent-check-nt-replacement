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

	t.Run("Outside a defined range", func(t *testing.T) {
		parsedThing := ParseRangeString("5:33")
		assert.Equal(t, parsedThing.Start, 5.0)
		assert.Equal(t, parsedThing.End, 33.0)

		shouldAlert := parsedThing.CheckRange("54")
		assert.Equal(t, shouldAlert, true)
		shouldAlert = parsedThing.CheckRange("4")
		assert.Equal(t, shouldAlert, true)

	})

	t.Run("Within a defined range", func(t *testing.T) {
		parsedThing := ParseRangeString("10:200")
		assert.Equal(t, parsedThing.Start, 10.0)
		assert.Equal(t, parsedThing.End, 200.0)
		shouldAlert := parsedThing.CheckRange("54")

		assert.Equal(t, shouldAlert, false)
	})

}
