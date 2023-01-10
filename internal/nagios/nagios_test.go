package nagios

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Best documentation for this is: https://www.monitoring-plugins.org/doc/guidelines.html#THRESHOLDFORMAT

func TestThis(t *testing.T) {
	t.Run("Test This Baby", func(t *testing.T) {
		parsedThing := ParseRangeString("10")
		assert.Equal(t, parsedThing.End, 10.0)

		shouldAlert := parsedThing.CheckRange("54")
		assert.Equal(t, shouldAlert, true)
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

		assert.Equal(t, parsedThing.CheckRange("54"), false)
		assert.Equal(t, parsedThing.CheckRange("65535.7"), false)
		assert.Equal(t, parsedThing.CheckRange("-54"), true)
		assert.Equal(t, parsedThing.CheckRange("49"), true)
	})

	t.Run("Alert in 0-32", func(t *testing.T) {
		parsedThing := ParseRangeString("@32")

		assert.Equal(t, parsedThing.CheckRange("32"), true)
		assert.Equal(t, parsedThing.CheckRange("31"), true)
		assert.Equal(t, parsedThing.CheckRange("33"), false)
		assert.Equal(t, parsedThing.CheckRange("-32"), false)
	})

	t.Run("InsideRange", func(t *testing.T) {
		parsedThing := ParseRangeString("@32:64")
		assert.Equal(t, parsedThing.CheckRange("32"), true)
		assert.Equal(t, parsedThing.CheckRange("33"), true)
		assert.Equal(t, parsedThing.CheckRange("64"), true)
		assert.Equal(t, parsedThing.CheckRange("63"), true)
		assert.Equal(t, parsedThing.CheckRange("31"), false)
		assert.Equal(t, parsedThing.CheckRange("65"), false)
	})

}
