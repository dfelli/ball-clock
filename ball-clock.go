package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const MIN_SIZE int = 5
const FIVEMIN_SIZE int = 12
const HOUR_SIZE int = 12

const MIN_BALL_COUNT = 27
const MAX_BALL_COUNT = 127

type Ball int
type Sequences struct {
	Main    []Ball
	Min     []Ball
	FiveMin []Ball
	Hour    []Ball
}

// fills an array sequencially with balls numbered 1 to size and returns it
func fillSequencialArray(size int) []Ball {
	array := make([]Ball, size)
	for index := 1; index <= size; index++ {
		array[index-1]=Ball(index)
	}
	return array
}

//checks and returns true if an array ordered sequencially
func sequenceIsOrdered(array []Ball) bool {
	for index := 0; index < len(array)-1; index++ {
		if array[index] >= array[index+1] {
			return false
		}
	}
	return true
}

//prints the sequences struct as formatted json
func printSequences(sequences *Sequences) {
	jsonOutput, err := json.Marshal(*sequences)
	if err != nil {
		fmt.Printf("error converting to json\n")
		return
	}
	fmt.Printf("%v\n", string(jsonOutput))
}

// simulates one minute of time passing for a ball clock
func simulateOneMinute(sequences *Sequences) {
	// store the ball we are taking from the main sequence
	ballToMove := sequences.Main[0]
	// remove the first ball of the main sequence
	sequences.Main = sequences.Main[1:]
	// add the ball to the minute sequence
	sequences.Min = append(sequences.Min, ballToMove)

	// if a sequence is full, dump it
	// a dump is only needed if a ball fills the track, if lower "value" track doesn't dump
	// other of the higher tracks will fill and require a dump
	if len(sequences.Min) == MIN_SIZE {
		dump(&sequences.Min, &sequences.FiveMin, &sequences.Main)

		if len(sequences.FiveMin) == FIVEMIN_SIZE {
			dump(&sequences.FiveMin, &sequences.Hour, &sequences.Main)

			if len(sequences.Hour) == HOUR_SIZE {
				dump(&sequences.Hour, &sequences.Main, &sequences.Main)
			}
		}
	}
}

// performs a ball clock dump when a holding track is full.
// balls up to the last one return to the main pool in reverse order.
// then the last ball returns to the next higher ordered track
// in the case where the hour sequence is dumps (it goes to the main pool)
// the last ball on the hour track returns to the main pool last
func dump(dumpSequence *[]Ball, nextSequence *[]Ball, Main *[]Ball) {
	lastBall := (*dumpSequence)[len(*dumpSequence)-1]
	for index := len(*dumpSequence) - 2; index >= 0; index-- {
		*Main = append(*Main, (*dumpSequence)[index])
	}
	*nextSequence = append(*nextSequence, lastBall)
	*dumpSequence = (*dumpSequence)[:0]
}

// checks to see if input is a valid integer in between two values min and max
func validateIntBetween(input string, min int, max int, text string) (number int, err error) {
	number, err = strconv.Atoi(input)
	if err != nil {
		err = fmt.Errorf("Error: non integer provided: %s : %v. Exiting", input, err)
		return number, err
	}
	if number < min || number > max {
		err = fmt.Errorf("Error: the value for %s must be between %d and %d inclusive:"+
			"You provided %d. Exiting", text, min, max, number)
	}
	return number, err
}

// expects the user to provide the num of balls an integer between 27 and 127 inclusive
// and additionally a second parameter integer specifying the number of minutes to simulate
// sets up mode baseed on number of parameters. 1 parameter = mode 1. 2 parameters is mode 2.
// otherwise throws and error
func handleUserInput(args []string) (numBalls, minutesToSimulate, mode int, err error) {
	if len(args) == 3 {
		mode = 2
		minutesToSimulate, err = validateIntBetween(args[2], 0, math.MaxInt64, "minutes to simulate")
	} else if len(args) == 2 {
		mode = 1
		minutesToSimulate = math.MaxInt64
	} else {
		fmt.Printf("Please Choose one of the following:" +
			"\n\tMode 1: by entering in the number of balls (int) to simulate %d to %d" +
			"\n or"+
			"\n\tMode 2: by entering in the number of balls (int) to simulate %d to %d and" +
			"\n\tthe number of minutes (int) to run the simulation separated by a space" +
			"\n",MIN_BALL_COUNT, MAX_BALL_COUNT, MIN_BALL_COUNT, MAX_BALL_COUNT)
		err = fmt.Errorf("improper number of arguments provided. expects 1 to 2 arguments. Exiting")
	}
	if err != nil {
		return numBalls, minutesToSimulate, mode, err
	}
	numBalls, err = validateIntBetween(args[1], MIN_BALL_COUNT, MAX_BALL_COUNT, "number of balls")

	return numBalls, minutesToSimulate, mode, err
}

func main() {
	// handle and validate what the user gives as input
	numBalls, minutesToSimulate, mode, err := handleUserInput(os.Args)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	//creates sequences
	sequences := Sequences{}
	// fills up the main sequence with balls
	sequences.Main = fillSequencialArray(numBalls)

	finished := false
	minutesSimulated := 0

	// note the start time
	startTime := time.Now()

	// run the simulation until it finished
	for !finished {
		simulateOneMinute(&sequences)
		// increment the minutes of the simulation
		minutesSimulated += 1
		//check for the end conditions depending on the mode
		if mode == 2 && minutesSimulated >= minutesToSimulate {
			finished = true
		} else if mode == 1 && len(sequences.Main) == numBalls && sequenceIsOrdered(sequences.Main) {
			finished = true
		}
	}

	// calculate the time difference or time it took to run the simulation in seconds and milliseconds
	timeDiff := time.Now().Sub(startTime)
	// rounds milliseconds to the nearest whole number. to match the formatting seconds receives with setting
	// precision wit fmt.Printf("%3f")
	milliseconds := int(timeDiff.Seconds()*1000+.5)
	seconds := timeDiff.Seconds()

	//output the results depending on the mode
	if mode == 2 {
		printSequences(&sequences)
	} else if mode == 1 {
		fmt.Printf("%v balls cycle after %v days.\n", numBalls, int(minutesSimulated/60.0/24.0))
	}

	// print how much time it took.
	fmt.Printf("Completed in %d milliseconds (%.3f seconds)\n", milliseconds, seconds)
}
