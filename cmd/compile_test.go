package cmd_test

	
import (
	"testing"
	"github.com/speechly/cli/cmd"
	"math"
)

func checkResultRowSliceEqual(t *testing.T, a []cmd.ResultRow, b []cmd.ResultRow) {
	if len(a) != len(b) {
		t.Errorf("Not equal length") 
	}
	for i := 0 ; i<len(a); i++ {
		if a[i].Name != b[i].Name {
			t.Errorf("Elem %d should have name %s but had %s",i, a[i].Name, b[i].Name)
		}
		if a[i].Count != b[i].Count {
			t.Errorf("Elem %d should have count %d but had %d",i, a[i].Count, b[i].Count)
		}
		eps := 0.001
		if math.Abs(float64(a[i].Distrib) - float64(b[i].Distrib)) > eps {
			t.Errorf("Elem %d should have distib %f but had %f",i, a[i].Distrib, b[i].Distrib)
		}
		if math.Abs(float64(a[i].Proportion) - float64(b[i].Proportion)) > eps {
			t.Errorf("Elem %d should have proportion %f but had %f",i, a[i].Proportion, b[i].Proportion)
		}
	}
}

func TestGetIntentAndEntityCounts(t *testing.T) {
	examples := []string{
		`*turn_on turn off the [lights](device) in the [kitchen](room)`,
		`*turn_off turn off the [air conditioner](device) in the [living room](room) 
		and *turn_on turn on the [heating](device) in the [garage](room)`,
		`*switch_off switch off the [lights](device) in the [kitchen](room)`,
		`*turn_off turn off the [air conditioner](device) in the [kitchen](room)`,
		`*turn_off turn off the [air conditioner](device) in the [kitchen](room)`,
		`*switch_on switch on something [tomorrow](time) please`,
	}
	counter := cmd.CreateCounter(examples)
	
	expected := []cmd.ResultRow{
		cmd.ResultRow{Name: "turn_off", Count: 3, Distrib: 3.0/7.0, Proportion: 3.0/7.0},
		cmd.ResultRow{Name: "turn_on", Count: 2, Distrib: 2.0/7.0, Proportion: 2.0/7.0},
		cmd.ResultRow{Name: "switch_off", Count: 1, Distrib: 1.0/7.0, Proportion: 1.0/7.0},
		cmd.ResultRow{Name: "switch_on", Count: 1, Distrib: 1.0/7.0, Proportion: 1.0/7.0},
	}
	checkResultRowSliceEqual(t, expected, counter.GetIntentCounts())
	
	expected = []cmd.ResultRow{
		cmd.ResultRow{Name: "device", Count: 6, Distrib: 6.0/13.0, Proportion: 6.0/7.0},
		cmd.ResultRow{Name: "room", Count: 6, Distrib: 6.0/13.0, Proportion: 6.0/7.0},
		cmd.ResultRow{Name: "time", Count: 1, Distrib: 1.0/13.0, Proportion: 1.0/7.0},
	}
	checkResultRowSliceEqual(t, expected, counter.GetEntityTypeCounts())
}

func TestGetIntentEntityValueCounts(t *testing.T) {
	examples := []string{
		`*order [large](size) [coffee](coffee) with [cream](addition) and [cream](addition)`,
		`*order [small](size) [coffee](coffee) with [sugar](addition) and [milk](addition)
		*order i want a [americano](coffee) with [milk](addition) and [syrup](addition) please`,
		`*order order [latte](coffee) with [cream](addition) please`,
		`*order [small](size) [double](shot) [coffee](coffee) please
		*order order [small](size) [single](shot) [americano](coffee) with [syrup](addition)`,
		`*order may i have [triple](shot) [cafe latte](coffee)`,
		`*order [espresso](coffee) with [syrup](addition) please
		*order [medium](size) [double](shot) [espresso](coffee) please`,
		`*order [medium](size) [single](shot) [cafe latte](coffee) with [milk](addition) please`,
		`*order [double](shot) [latte](coffee) please`,
		`*order i'd like to have an [cappuccino](coffee) with [sugar](addition) please`,
		`*order [double](shot) [cafe latte](coffee) with [syrup](addition)
		*order [latte](coffee)`,
	}
	counter := cmd.CreateCounter(examples)

	var totalEnt, totalUtt float32
	totalEnt = 39.0
	totalUtt = 14.0
	names := []string{
		"order(addition=syrup)","order(shot=double)","order(addition=cream)","order(addition=milk)",
		"order(coffee=cafe latte)","order(coffee=coffee)","order(coffee=latte)","order(size=small)",
		"order(addition=sugar)","order(coffee=americano)","order(coffee=espresso)","order(shot=single)",
		"order(size=medium)","order(coffee=cappuccino)","order(shot=triple)","order(size=large)",
	}
	counts := []float32{4.0,4.0,3.0,3.0,3.0,3.0,3.0,3.0,2.0,2.0,2.0,2.0,2.0,1.0,1.0,1.0}
	expected := make([]cmd.ResultRow,0)
	for i := 0; i < len(names); i++ {
		cnt := counts[i]
		expected = append(expected, cmd.ResultRow{Name: names[i], Count: int(cnt), Distrib: cnt/totalEnt, Proportion: cnt/totalUtt})
	}
	checkResultRowSliceEqual(t, expected, counter.GetIntentEntityValueCounts())
}
