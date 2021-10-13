package cmd_test

import (
	"github.com/speechly/cli/cmd"
	"testing"
)

func TestCompareUtterances(t *testing.T) {
	annotatedExamples := []string{
		`*order [large|L](size) [coffee|coffee](coffee) with [cream|432](addition) and [cream|432](addition)`,
		`*order [small|S](size) [coffee|coffee](coffee) with [sugar|21](addition) and [milk|98](addition)`,
		`*order order [latte|latte](coffee) with [cream|432](addition) please`,
		`*order order [latte|latte](coffee) with [cream|432](addition) please`,
		`*order order [latte|latte](coffee) with [cream|432](addition) please`,
		`*help can I see the menu`,
		`*help can I see the menu`,
		`*order i'd like to have an [cappuccino|cappuccino](coffee) with [sugar|21](addition) please`,
		`*order i'd like to have an [cappuccino|cappuccino](coffee) with [sugar|21](addition) please`,
		`*order i'd like to have an [cappuccino|cappuccino](coffee) with [sugar|21](addition) please`,
		`*order may i have [3|triple](shot) [cafe latte|latte](coffee)`,
	}
	groundTruthExamples := []string{
		`*order [large|large](size) [coffee|coffee](coffee) with [cream|432](addition) and [cream|432](addition)`,
		`*order [small|S](size) [coffee|coffee](coffee) with [sugar|21](size) and [milk|98](addition)`,
		`*order order [latte|latte](coffee) with [cream|432](addition) please`,
		`*order order [latte|latte](coffee) with [cream|432](addition)`,
		`*order order [latte|latte](coffee)`,
		`*help can I see the menu`,
		`*help can I see the menu please`,
		`*order i would like to have an [cappuccino|cappuccino](coffee) with [sugar|21](addition) please`,
		`*order i'd like to have an [cappuccino|](coffee) with [sugar|221](addition) please`,
		`*order i'd like to have an [cappuccino|cappuccino](coffee) with [sugar|21](coffee) please`,
		`*order may i have [triple|triple](shot) [cafe latte|latte](coffee)`,
	}
	expectedResults := []bool{true, false, true, false, false, true, false, false, true, false, false}
	comp := cmd.CreateComparator()
	for i, a := range annotatedExamples {
		b := groundTruthExamples[i]
		exp := expectedResults[i]
		res := comp.Equal(a, b)
		if res != exp {
			t.Errorf("comparing %v and %v should yield %t", a, b, exp)
		}
	}
}
