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

	t.Run("Within a range involving -inf", func(t *testing.T) {
		parsedThing := ParseRangeString("~:30")
		assert.Equal(t, parsedThing.Start_Infinity, true)
		assert.Equal(t, parsedThing.End, 30.0)

		shouldAlert := parsedThing.CheckRange("5")
		assert.Equal(t, shouldAlert, false)

		shouldAlert = parsedThing.CheckRange("-10")
		assert.Equal(t, shouldAlert, false)

		shouldAlert = parsedThing.CheckRange("-100")
		assert.Equal(t, shouldAlert, false)

		shouldAlert = parsedThing.CheckRange("31")
		assert.Equal(t, shouldAlert, true)
	})

	t.Run("Within a range involving +inf", func(t *testing.T) {
		parsedThing := ParseRangeString("50:~")
		assert.Equal(t, parsedThing.Start, 50.0)
		assert.Equal(t, parsedThing.End_Infinity, true)

		shouldAlert := parsedThing.CheckRange("54")
		assert.Equal(t, shouldAlert, false)

		shouldAlert = parsedThing.CheckRange("65535")
		assert.Equal(t, shouldAlert, false)

		shouldAlert = parsedThing.CheckRange("-54")
		assert.Equal(t, shouldAlert, true)

		shouldAlert = parsedThing.CheckRange("49")
		assert.Equal(t, shouldAlert, true)
	})

	//TODO check this, see what is the correct state
	t.Run("Single Value Only", func(t *testing.T) {
		parsedThing := ParseRangeString("@32")
		//assert.Equal(t, parsedThing.Start, 32.0)
		//assert.Equal(t, parsedThing.End, 32.0)

		shouldAlert := parsedThing.CheckRange("32")
		assert.Equal(t, shouldAlert, true)

		shouldAlert = parsedThing.CheckRange("31")
		assert.Equal(t, shouldAlert, true)

		shouldAlert = parsedThing.CheckRange("33")
		assert.Equal(t, shouldAlert, true)

		shouldAlert = parsedThing.CheckRange("-32")
		assert.Equal(t, shouldAlert, true)
	})

}
