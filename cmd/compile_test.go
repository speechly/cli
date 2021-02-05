package cmd_test

	
import (
	"testing"
	"github.com/speechly/cli/cmd"
)

func TestGetIntentAndEntityStats(t *testing.T) {
	examples := []string{
		"*turn_on turn off the [lights](device) in the [kitchen](room)",
		`*turn_off turn off the [air conditioner](device) in the [living room](room) 
		and *turn_on turn on the [heating](device) in the [garage](room)`,
	}
	intents, entityTypes, entityValues := cmd.GetIntentAndEntityStats(examples)
	for _, expectedEntityVal := range []string{"lights","kitchen","living room"} {
		if _,ok := entityValues[expectedEntityVal]; !ok {
			t.Errorf("Should contain %s", expectedEntityVal)
		}
	}
	for _, expectedEntityType := range []string{"device","room"} {
		if _,ok := entityTypes[expectedEntityType]; !ok {
			t.Errorf("Should contain %s", expectedEntityType)
		}
	}
	for _, intent := range []string{"turn_on","turn_off"} {
		if _,ok := intents[intent]; !ok {
			t.Errorf("Should contain %s", intent)
		} 
	}
	if val, _ := entityTypes["room"]; int(val) != 3 {
		t.Errorf("The count should be 3")
	}
}

