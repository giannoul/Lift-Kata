package lift

import (
	"sort"
)

// Direction ..
type Direction int

// Directions ..
const (
	Up Direction = iota
	Down
)

// Call ..
type Call struct {
	Floor     int
	Direction Direction
}

// Lift ..
type Lift struct {
	ID        string
	Floor     int
	Requests  []int
	DoorsOpen bool
}

// System ..
type System struct {
	floors []int
	lifts  []Lift
	calls  []Call
}

// CallAssignment ..
type CallAssignment struct {
	call      *Call
	lift      *Lift
	fulfilled bool
}

var callAssignments []CallAssignment
var CA CallAssignment

func (ca CallAssignment) add(c *Call, l *Lift) {
	callAssignments = append(callAssignments, CallAssignment{c, l, false})
}

func removeCallAssignmentByIndex(i int) {
	callAssignments = append(callAssignments[0:i], callAssignments[i+1:]...)
}

// NewSystem ..
func NewSystem() *System {
	return &System{floors: []int{}, lifts: []Lift{}, calls: []Call{}}
}

// AddFloors ..
func (s *System) AddFloors(floors ...int) {
	s.floors = append(s.floors, floors...)
}

// AddLifts ..
func (s *System) AddLifts(lifts ...Lift) {
	s.lifts = append(s.lifts, lifts...)
}

// AddCalls ..
func (s *System) AddCalls(calls ...Call) {
	s.calls = append(s.calls, calls...)
}

// CallsFor ..
func (s System) CallsFor(floor int) (calls []Call) {
	calls = []Call{}
	for _, c := range s.calls {
		if c.Floor == floor {
			calls = append(calls, c)
		}
	}
	return calls
}

/*
If we need to print something while testing we can use:
fmt.Fprintf(os.Stdout, "%#v\n", s.lifts[i])
*/

// Tick ..
func (s *System) Tick() {
	s.UpdateState()
}

// UpdateState
func (s *System) UpdateState() {
	//1. distribute the calls according to a cost function
	s.distributeCalls()
	//2. check if an elevator has reached a floor whre a call came from
	s.checkCallsFulfillment()
	//3. move the lifts according to their existing requests queue
	s.updateLiftStatuses()
	s.cleanupCallsFulfillment()

}

/*
AddPanelReqToLift

This method will add a request coming from the pressing of the floor
button in the lift's panel
*/
func (s *System) AddPanelReqToLift(id string, floor int) {
	for i := 0; i < len(s.lifts); i++ {
		if s.lifts[i].ID == id {
			s.lifts[i].AddRequestToLift(floor)
		}
	}
}

/*
AddRequestToLift

This method will add the requested (from the panel inside the lift) floor to the lift's
Requests queue
*/
func (l *Lift) AddRequestToLift(reqFloor int) bool {

	if len(l.Requests) == 0 {
		l.Requests = append(l.Requests, reqFloor)
		return true
	}

	currentDirection := l.direction()

	newSlice := []int{}
	s := make([]int, len(l.Requests))
	copy(s, l.Requests)

	if currentDirection == Down {
		sort.Sort(sort.Reverse(sort.IntSlice(s)))
	}
	if currentDirection == Up {
		sort.Ints(s)
	}

	pos := 0
	for i, f := range s {
		if currentDirection == Down && l.Floor > f {
			pos = i
			break
		}

		if currentDirection == Up && l.Floor < f {
			pos = i
			break
		}
	}

	if currentDirection == Down {
		newSlice = s[pos:]
		tmp := s[0:pos]
		sort.Ints(tmp)
		newSlice = append(newSlice, tmp...)

	}

	if currentDirection == Up {
		newSlice = s[pos:]
		tmp := s[0:pos]
		sort.Sort(sort.Reverse(sort.IntSlice(tmp)))
		newSlice = append(newSlice, tmp...)
	}

	l.Requests = newSlice
	return true
}

/*
isOnRequestedFloor

This method justs checks if the lift is on a floor where there is a Request for it
*/
func (l *Lift) isOnRequestedFloor() bool {
	return len(l.Requests) > 0 && l.Floor == l.Requests[0]
}

/*
removeRequestAndOpenDoors

This method removes the fullfilled request from the queue and opens the lift's doors
*/
func (l *Lift) removeRequestAndOpenDoors() {
	l.Requests = l.Requests[1:]
	l.DoorsOpen = true
}

/*
isIdle

This method checks if there is no Request on the lift's queue and the doors are closed.
This means that the System should be able to move the lift.
*/
func (l Lift) isIdle() (idle bool) {
	return len(l.Requests) == 0 && l.DoorsOpen == false
}

/*
checkForCallPickup

This method checks if there is a callAssignment for the specific lift for the current floor
*/
func (l *Lift) checkForCallPickup() {
	for j := 0; j < len(callAssignments); j++ {
		if l.ID != callAssignments[j].lift.ID {
			continue
		}
		if len(l.Requests) > 0 && l.direction() != callAssignments[j].call.Direction {
			continue
		}
		if l.Floor == callAssignments[j].call.Floor {
			l.DoorsOpen = true
			break
		}
		if l.Floor == callAssignments[j].call.Floor && l.DoorsOpen == true {
			l.DoorsOpen = false
			break
		}
	}
}

/*
move

This method is intended to be used by the System
*/
func (l *Lift) move(direction Direction) {
	if direction == Down {
		l.Floor--
	}
	if direction == Up {
		l.Floor++
	}
}

/*
computecostToFloor

This method computes the score of the lift in order to get to the floor where a call exists
*/
func (l *Lift) computecostToFloor(floor int) (cost int) {
	cost = 0
	if l.DoorsOpen == true {
		cost++
	}
	cost += Abs(l.Floor - floor)
	return cost
}

/*
moveLift

This method moves the lift taking into account its requests queue
*/
func (l *Lift) moveLift() {
	if l.DoorsOpen == true {
		return
	}
	if len(l.Requests) > 0 && l.Requests[0] == l.Floor {
		return
	}
	ldirection := l.direction()
	if len(l.Requests) > 0 && ldirection == Down {
		l.Floor--
	}
	if len(l.Requests) > 0 && ldirection == Up {
		l.Floor++
	}
}

/*
direction

Determine current direction of elevator
*/
func (l Lift) direction() Direction {
	var currentDirection Direction

	if len(l.Requests) > 0 && l.Floor < l.Requests[0] {
		currentDirection = Up
	}

	if len(l.Requests) > 0 && l.Floor > l.Requests[0] {
		currentDirection = Down
	}
	return currentDirection
}

/*
updateLiftStatuses

This is the main logic from the system's perspective
*/
func (s *System) updateLiftStatuses() {
	for i := 0; i < len(s.lifts); i++ {
		if s.lifts[i].DoorsOpen == true {
			s.lifts[i].DoorsOpen = false
			continue
		}

		// lift is on the correct floor regarding a request -> open doors
		if s.lifts[i].isOnRequestedFloor() && s.lifts[i].DoorsOpen == false {
			s.lifts[i].removeRequestAndOpenDoors()
			continue
		}
		// lift is on the requested floor with doors open -> close doors
		if s.lifts[i].isOnRequestedFloor() && s.lifts[i].DoorsOpen == true {
			s.lifts[i].DoorsOpen = false
			continue
		}
		// check if a call exists for the current floor
		s.lifts[i].checkForCallPickup()

		// lift is idle (no requests in its queue) and there is a pending call for it -> move it
		if s.lifts[i].isIdle() {
			var direction Direction
			for j := 0; j < len(callAssignments); j++ {
				if s.lifts[i].ID != callAssignments[j].lift.ID {
					continue
				}
				if callAssignments[j].lift.Floor < callAssignments[j].call.Floor {
					direction = Up
				} else {
					direction = Down
				}
				if callAssignments[j].fulfilled == false && callAssignments[j].lift.DoorsOpen == false {
					callAssignments[j].lift.move(direction)
				}

			}
			continue
		}
		// lift is idle and there is a pending call for it -> move it
		s.lifts[i].moveLift()
	}
}

/*
distributeCalls

This method distributes the calls to lifts according to a cost function
*/
func (s *System) distributeCalls() {
	undistributed := Abs(len(callAssignments) - len(s.calls))
	if undistributed > 0 {
		for i := 0; i < undistributed; i++ {
			s.assignCallToLift(&s.calls[i])
		}
	}
}

/*
assignCallToLift

This method determines the cost for a call for each lift. It will assign
the call to the lift with the minimum cost
*/
func (s *System) assignCallToLift(c *Call) {
	var (
		idleIndex, closestIndex, sameDirectionindex = -1, -1, -1
	)
	min := len(s.floors)
	for i := 0; i < len(s.lifts); i++ {
		// worst option: since it's idle just use it
		if s.lifts[i].isIdle() {
			idleIndex = i
		}
		// normal option: it's going in the same direction as the caller
		if s.lifts[i].direction() == c.Direction {
			sameDirectionindex = i
		}
		// best option: closer lift going to the same direction as the caller
		if s.lifts[i].computecostToFloor(c.Floor) < min && s.lifts[i].direction() == c.Direction {
			closestIndex = i
		}
	}

	if closestIndex > -1 {
		CA.add(c, &s.lifts[closestIndex])
	} else if sameDirectionindex > -1 {
		CA.add(c, &s.lifts[sameDirectionindex])
	} else if idleIndex > -1 {
		CA.add(c, &s.lifts[idleIndex])
	}

}

/*
checkCallsFulfillment

This is mainly for presentation purposes. The call sign (^ or v) should be removed
once the call is fullfilled.
*/
func (s *System) checkCallsFulfillment() {
	for i := 0; i < len(callAssignments); i++ {
		if callAssignments[i].fulfilled == false {
			if callAssignments[i].lift.Floor == callAssignments[i].call.Floor && callAssignments[i].lift.DoorsOpen == true {
				callAssignments[i].fulfilled = true
			}
		}
	}
}

/*
removeCallByIndex

This method will just remove the call from the queue (for display reasons)
*/
func (s *System) removeCallByIndex(i int) {
	s.calls = append(s.calls[:i], s.calls[i+1:]...)
}

/*
cleanupCallsFulfillment

This method has the logic to find a call for removal
*/
func (s *System) cleanupCallsFulfillment() {
	indexesToRemove := []int{}
	for i := 0; i < len(callAssignments); i++ {
		if callAssignments[i].fulfilled {
			indexesToRemove = append(indexesToRemove, i)
		}
	}

	for i := 0; i < len(indexesToRemove); i++ {
		for j := 0; j < len(s.calls); j++ {
			if s.calls[j].Floor == callAssignments[indexesToRemove[i]].call.Floor {
				s.removeCallByIndex(j)
			}
		}
		removeCallAssignmentByIndex(indexesToRemove[i])
	}

}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
