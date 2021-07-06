package lift_test

import (
	"fmt"
	"os"
	"testing"

	approvaltests "github.com/approvals/go-approval-tests"
	"github.com/lift-kata/lift"
)

type Lift struct {
	ID        string
	Floor     int
	Requests  []int
	DoorsOpen bool
}

func PrintStateToStdout(x string) {
	if true {
		fmt.Fprintln(os.Stdout, x)
	}
}

func TestSingleLiftFiveFloorsOneCallOnFloor3(t *testing.T) {
	liftSystem := lift.NewSystem()
	liftSystem.AddLifts(
		lift.Lift{"A", 0, []int{}, false})
	liftSystem.AddFloors(0, 1, 2, 3, 4)
	// Initial state
	x := lift.PrintLifts(liftSystem, lift.NewPrinter())
	PrintStateToStdout(x)
	liftSystem.AddCalls(lift.Call{3, lift.Down})

	for i := 0; i < 6; i++ {
		liftSystem.Tick()
		x = lift.PrintLifts(liftSystem, lift.NewPrinter())
		PrintStateToStdout(x)
	}

	approvaltests.VerifyString(t, x)
}
func TestSingleLiftEightFloorsOneCallOnFloorTwo(t *testing.T) {
	liftSystem := lift.NewSystem()
	liftSystem.AddLifts(
		lift.Lift{"A", 4, []int{3, 7}, false})
	liftSystem.AddFloors(0, 1, 2, 3, 4, 5, 6, 7, 8)
	// Initial state
	x := lift.PrintLifts(liftSystem, lift.NewPrinter())
	PrintStateToStdout(x)

	for i := 0; i < 4; i++ {
		liftSystem.Tick()
		x = lift.PrintLifts(liftSystem, lift.NewPrinter())
		PrintStateToStdout(x)
	}

	liftSystem.AddCalls(lift.Call{2, lift.Down})
	for i := 0; i < 12; i++ {
		liftSystem.Tick()
		x = lift.PrintLifts(liftSystem, lift.NewPrinter())
		PrintStateToStdout(x)
	}
	approvaltests.VerifyString(t, x)
}

func TestTwoLiftsNineFloorsOneCallOnFloorFiveOneCallOnFloorTwo(t *testing.T) {
	liftSystem := lift.NewSystem()
	liftSystem.AddLifts(
		//		  ID ,	Floor	Requests	DoorsOpen
		lift.Lift{"A", 4, []int{3, 7}, false},
		lift.Lift{"B", 4, []int{4, 8}, false})
	liftSystem.AddFloors(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	// Initial state
	x := lift.PrintLifts(liftSystem, lift.NewPrinter())
	PrintStateToStdout(x)

	// do some Ticks
	liftSystem.AddCalls(lift.Call{5, lift.Up})
	for i := 0; i < 8; i++ {
		liftSystem.Tick()
		x = lift.PrintLifts(liftSystem, lift.NewPrinter())
		PrintStateToStdout(x)
	}

	// Add a call to floor 2
	liftSystem.AddCalls(lift.Call{2, lift.Down})
	// do some Ticks until all the calls are fulfilled
	for i := 0; i < 8; i++ {
		liftSystem.Tick()
		x = lift.PrintLifts(liftSystem, lift.NewPrinter())
		PrintStateToStdout(x)
	}
	// The person that made  the call on floor 2 is now onboarded
	// and on the lift's panel clicked the floor 5
	liftSystem.AddPanelReqToLift("A", 5)
	// do some Ticks until the request for floor 5 is fulfilled
	for i := 0; i < 5; i++ {
		liftSystem.Tick()
		x = lift.PrintLifts(liftSystem, lift.NewPrinter())
		PrintStateToStdout(x)
	}
	// let's do another Tick in order to verify that the status is stable
	liftSystem.Tick()
	x = lift.PrintLifts(liftSystem, lift.NewPrinter())
	PrintStateToStdout(x)

	approvaltests.VerifyString(t, x)
}
