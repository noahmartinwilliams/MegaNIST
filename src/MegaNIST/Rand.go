package main

import "math/rand"

type Coord struct {
	X int
	Y int
}

func GetRands() chan int {
	retc := make(chan int, 1024)
	go func() {
		defer close(retc)
		retc <- rand.Int()
	} ()
	return retc
}

func GetRandCoords(inputc chan int, maxX int, maxY int) chan Coord {
	retc := make(chan Coord, 1024)
	go func() {
		defer close(retc)
		for {
			x, done := <- inputc
			if done {
				return
			}
			y, done2 := <- inputc

			if done2 {
				return
			}
			retc <- Coord{X:(x%maxX), Y:(y%maxY)}
		}
	} ()
	return retc
}

func GetRandAngles(inputc chan int, maxAngle float64) chan float64 {
	retc := make(chan float64, 1024)
	go func() {
		defer close(retc)
		for input := range inputc {
			input2 := input % (int(maxAngle * 10000))
			retc <- (float64(input2) / 10000)
		}
	} ()
	return retc
}
