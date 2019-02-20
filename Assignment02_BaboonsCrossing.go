// Started 14.02.19
// Ended   20.02.19

package main

import (
	s "CADP/semaphore"
	"fmt"
	"math/rand"
	"runtime"
	t "time"
)

// Set this to true to display diagnostic info
var debug = true

// Split binary semaphore
var signalEastMale = s.Make(0)
var signalEastFemale = s.Make(0)
var signalWestMale = s.Make(0)
var signalWestFemale = s.Make(0)
var weightLock = s.Make(0)
var e = s.Make(1)

// Change this variable as needed
// Controls the amount of baboons spawned
const baboonCount = 100

// Rope control variables
const ropeCapacity = 50

var ropeWeight = 0

// Baboon weights
const maleWeight = 20
const femaleWeight = 10

// Fairness controller
var westHeadingLast = false

// Counters
var eastHeadingCount = 0
var delayedEastHeading = 0
var delayedEastHeadingMales = 0
var delayedEastHeadingFemales = 0

var westHeadingCount = 0
var delayedWestHeading = 0
var delayedWestHeadingMales = 0
var delayedWestHeadingFemales = 0

func main() {
	// Spawns baboonCount baboons, divided equally for each direction.
	// As a baboon has crossed the canyon, it will "respawn".
	for i := 0; i < baboonCount; i++ {
		gender := rand.Intn(2)

		if i%2 == 0 {
			go eastHeadingBaboon(gender)
		} else if i%2 == 1 {
			go westHeadingBaboon(gender)
		}
	}

	for {
		// Prevent Go scheduler from getting stuck in this loop.
		runtime.Gosched()
	}
}

func eastHeadingBaboon(gender int) {
	for {
		// East side of canyon

		e.Acquire(1)
		if westHeadingCount != 0 || (gender == 0 && ropeWeight > 30) || (gender == 1 && ropeWeight > 40) {
			if gender == 0 {
				delayedEastHeadingMales++
			} else if gender == 1 {
				delayedEastHeadingFemales++
			}
			delayedEastHeading = delayedEastHeadingMales + delayedEastHeadingFemales
			e.Release(1)

			if gender == 0 {
				signalEastMale.Acquire(1)
			} else if gender == 1 {
				signalEastFemale.Acquire(1)
			}
		}

		// Entering the rope
		eastHeadingCount++
		if gender == 0 {
			ropeWeight += maleWeight
		} else if gender == 1 {
			ropeWeight += femaleWeight
		}

		// Climbing over
		if gender == 0 {
			fmt.Printf("[East][Male][westHeaded:%d][eastHeaded:%d][delayedWest:%d][delayedEast:%d][ropeWeight:%d]\n", westHeadingCount, eastHeadingCount, delayedWestHeading, delayedEastHeading, ropeWeight)
		} else if gender == 1 {
			fmt.Printf("[East][Fema][westHeaded:%d][eastHeaded:%d][delayedWest:%d][delayedEast:%d][ropeWeight:%d]\n", westHeadingCount, eastHeadingCount, delayedWestHeading, delayedEastHeading, ropeWeight)
		}

		SIGNAL()

		if debug {
			if westHeadingCount > 0 {
				fmt.Println("EAST: Collision!")
				t.Sleep(3000 * t.Millisecond)
			}
			if ropeWeight > ropeCapacity {
				fmt.Println("EAST: Rope broke!")
				t.Sleep(3000 * t.Millisecond)
			}
		}

		e.Acquire(1)

		// Left the rope
		if gender == 0 {
			ropeWeight -= maleWeight
		} else if gender == 1 {
			ropeWeight -= femaleWeight
		}

		eastHeadingCount--
		westHeadingLast = false

		SIGNAL()
	}
}

func westHeadingBaboon(gender int) {
	for {
		// West side of the canyon

		e.Acquire(1)
		if eastHeadingCount != 0 || (gender == 0 && ropeWeight > 30) || (gender == 1 && ropeWeight > 40) {
			if gender == 0 {
				delayedWestHeadingMales++
			} else if gender == 1 {
				delayedWestHeadingFemales++
			}
			delayedWestHeading = delayedWestHeadingMales + delayedWestHeadingFemales
			e.Release(1)

			if gender == 0 {
				signalWestMale.Acquire(1)
			} else if gender == 1 {
				signalWestFemale.Acquire(1)
			}
		}

		// Entering the rope
		westHeadingCount++
		if gender == 0 {
			ropeWeight += maleWeight
		} else if gender == 1 {
			ropeWeight += femaleWeight
		}

		// Climbing over
		if gender == 0 {
			fmt.Printf("[West][Male][westHeaded:%d][eastHeaded:%d][delayedWest:%d][delayedEast:%d][ropeWeight:%d]\n", westHeadingCount, eastHeadingCount, delayedWestHeading, delayedEastHeading, ropeWeight)
		} else if gender == 1 {
			fmt.Printf("[West][Fema][westHeaded:%d][eastHeaded:%d][delayedWest:%d][delayedEast:%d][ropeWeight:%d]\n", westHeadingCount, eastHeadingCount, delayedWestHeading, delayedEastHeading, ropeWeight)
		}

		SIGNAL()

		if debug {
			if eastHeadingCount > 0 {
				fmt.Println("WEST: Collision!")
				t.Sleep(3000 * t.Millisecond)
			}
			if ropeWeight > ropeCapacity {
				fmt.Println("WEST: Rope broke!")
				t.Sleep(3000 * t.Millisecond)
			}
		}

		e.Acquire(1)

		// Left the rope
		if gender == 0 {
			ropeWeight -= maleWeight
		} else if gender == 1 {
			ropeWeight -= femaleWeight
		}

		westHeadingCount--
		westHeadingLast = true

		SIGNAL()
	}
}

// SIGNAL is a baboon traffic controller
func SIGNAL() {
	randGender := rand.Intn(2)

	if westHeadingCount == 0 && delayedEastHeading > 0 && (delayedWestHeading == 0 || westHeadingLast) {
		if ropeWeight == 40 && delayedEastHeadingFemales > 0 { // If ropeWeight is 40 => female or none
			delayedEastHeadingFemales--
			delayedEastHeading--
			signalEastFemale.Release(1)
			// fmt.Println("Case01A")
		} else if ropeWeight < 40 { // If ropeWeight is less, either male or female
			if randGender == 0 {
				if delayedEastHeadingMales > 0 {
					delayedEastHeadingMales--
					delayedEastHeading--
					signalEastMale.Release(1)
					// fmt.Println("Case01Bi")
				} else {
					delayedEastHeadingFemales--
					delayedEastHeading--
					signalEastFemale.Release(1)
					// fmt.Println("Case01Bii")
				}
			} else if randGender == 1 {
				if delayedEastHeadingFemales > 0 {
					delayedEastHeadingFemales--
					delayedEastHeading--
					signalEastFemale.Release(1)
					// fmt.Println("Case01Biii")
				} else {
					delayedEastHeadingMales--
					delayedEastHeading--
					signalEastMale.Release(1)
					// fmt.Println("Case01Biv")
				}
			}
		} else {
			e.Release(1)
			// fmt.Println("Case01C")
		}
	} else if eastHeadingCount == 0 && delayedWestHeading > 0 && (delayedEastHeading == 0 || !westHeadingLast) {
		if ropeWeight == 40 && delayedWestHeadingFemales > 0 { // If ropeWeight is 40 => female or none
			delayedWestHeadingFemales--
			delayedWestHeading--
			signalWestFemale.Release(1)
			// fmt.Println("Case02A")
		} else if ropeWeight < 40 { // If ropeWeight is less, either male or female
			if randGender == 0 {
				if delayedWestHeadingMales > 0 {
					delayedWestHeadingMales--
					delayedWestHeading--
					signalWestMale.Release(1)
					// fmt.Println("Case02Bi")
				} else {
					delayedWestHeadingFemales--
					delayedWestHeading--
					signalWestFemale.Release(1)
					// fmt.Println("Case02Bii")
				}
			} else if randGender == 1 {
				if delayedWestHeadingFemales > 0 {
					delayedWestHeadingFemales--
					delayedWestHeading--
					signalWestFemale.Release(1)
					// fmt.Println("Case02Biii")
				} else {
					delayedWestHeadingMales--
					delayedWestHeading--
					signalWestMale.Release(1)
					// fmt.Println("Case02Biv")
				}
			}
		} else {
			e.Release(1)
			// fmt.Println("Case02C")
		}
	} else {
		e.Release(1)
		// fmt.Println("Case03")
	}
}

// func SIGNAL_OLD() {
// 	if westHeadingCount == 0 && delayedEastHeading > 0 {
// 		delayedEastHeading--
// 		signalEast.Release(1)
// 	} else if eastHeadingCount == 0 && delayedWestHeading > 0 {
// 		delayedWestHeading--
// 		signalWest.Release(1)
// 	} else {
// 		e.Release(1)
// 	}
// }
