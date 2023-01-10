package nagios

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Best documentation for this is: https://www.monitoring-plugins.org/doc/guidelines.html#THRESHOLDFORMAT

func TestThis(t *testing.T) {
	t.Run("Test 0 to N or alert", func(t *testing.T) {
		parsedThing := ParseRangeString("10")
		assert.Equal(t, parsedThing.End, 10.0)
		assert.Equal(t, parsedThing.CheckRange("54"), true)
		assert.Equal(t, parsedThing.CheckRange("-1"), true)
	})

	t.Run("Test N to infinity or alert", func(t *testing.T) {
		parsedThing := ParseRangeString("10:")
		assert.Equal(t, parsedThing.Start, 10.0)
		assert.Equal(t, parsedThing.End_Infinity, true)
		assert.Equal(t, parsedThing.CheckRange("10"), false)
		assert.Equal(t, parsedThing.CheckRange("9"), true)
		assert.Equal(t, parsedThing.CheckRange("-1"), true)
		assert.Equal(t, parsedThing.CheckRange("11"), false)
	})

	t.Run("Within a range involving -inf", func(t *testing.T) {
		parsedThing := ParseRangeString("~:30")
		assert.Equal(t, parsedThing.Start_Infinity, true)
		assert.Equal(t, parsedThing.End, 30.0)
		assert.Equal(t, parsedThing.CheckRange("5"), false)
		assert.Equal(t, parsedThing.CheckRange("-10"), false)
		assert.Equal(t, parsedThing.CheckRange("-100"), false)
		assert.Equal(t, parsedThing.CheckRange("30"), false)
		assert.Equal(t, parsedThing.CheckRange("31"), true)
	})

	t.Run("Outside a defined range", func(t *testing.T) {
		parsedThing := ParseRangeString("5:33")
		assert.Equal(t, parsedThing.Start, 5.0)
		assert.Equal(t, parsedThing.End, 33.0)
		assert.Equal(t, parsedThing.CheckRange("33"), false)
		assert.Equal(t, parsedThing.CheckRange("34"), true)
		assert.Equal(t, parsedThing.CheckRange("4"), true)
		assert.Equal(t, parsedThing.CheckRange("5"), false)

	})

	t.Run("Within a defined range", func(t *testing.T) {
		parsedThing := ParseRangeString("10:200")
		assert.Equal(t, parsedThing.Start, 10.0)
		assert.Equal(t, parsedThing.End, 200.0)
		assert.Equal(t, parsedThing.CheckRange("54"), false)
		assert.Equal(t, parsedThing.CheckRange("10"), false)
		assert.Equal(t, parsedThing.CheckRange("9"), true)
		assert.Equal(t, parsedThing.CheckRange("200"), false)
		assert.Equal(t, parsedThing.CheckRange("201"), true)
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

	t.Run("If invalid range is provided (with positive infinity) parsing should return nil", func(t *testing.T) {
		parsedThing := ParseRangeString("50:~")
		assert.Nil(t, parsedThing)
	})

	t.Run("Alert in 0-32", func(t *testing.T) {
		parsedThing := ParseRangeString("@32")

		assert.Equal(t, parsedThing.CheckRange("32"), true)
		assert.Equal(t, parsedThing.CheckRange("31"), true)
		assert.Equal(t, parsedThing.CheckRange("0"), true)
		assert.Equal(t, parsedThing.CheckRange("33"), false)
		assert.Equal(t, parsedThing.CheckRange("-32"), false)
		assert.Equal(t, parsedThing.CheckRange("-1"), false)
	})
	t.Run("Alert on value 32", func(t *testing.T) {
		parsedThing := ParseRangeString("@32:32")

		assert.Equal(t, parsedThing.CheckRange("32"), true)
		assert.Equal(t, parsedThing.CheckRange("31"), false)
		assert.Equal(t, parsedThing.CheckRange("0"), false)
		assert.Equal(t, parsedThing.CheckRange("33"), false)
		assert.Equal(t, parsedThing.CheckRange("-32"), false)
		assert.Equal(t, parsedThing.CheckRange("-1"), false)
	})
}
