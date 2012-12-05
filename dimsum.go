package main

import (
	"bufio"
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/rand"
	. "os"
	"os/signal"
	"syscall"
)

func readInput(f *File, c chan []byte) {
	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		c <- line
	}

	close(c)
}

func main() {
	var input *File
	var err error
	var inputName string
	var sampleSize int
	var position int64
	position = 0
	inputLines := make(chan []byte, 100)

	seed := make([]byte, 1)

	signals := make(chan Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	io.ReadFull(cryptorand.Reader, seed)
	rand.Seed(int64(seed[0]))

	flag.StringVar(&inputName, "input", "", "input file name")
	flag.IntVar(&sampleSize, "samplesize", 1, "how many lines to sample")
	flag.Parse()

	out := make([][]byte, sampleSize)

	if inputName == "" {
		input = Stdin
	} else {
		input, err = Open(inputName)

		if err != nil {
			panic(err)
		}
	}

	defer input.Close()

	go readInput(input, inputLines)

	for line := range inputLines {

		bail := false

		select {
		case foo := <-signals:
			bail = true
			fmt.Printf("%v", foo)
		default:
		}

		if position < int64(sampleSize) {
			out[position] = line

		} else {
			r := rand.Int63n(position + 1)

			if err != nil {
				panic(err)
			}

			if r < int64(sampleSize) {
				out[r] = line
			}
		}

		if bail {
			break
		}

		position = position + 1
	}

	for _, line := range out {
		fmt.Printf("%s", line)
	}
}
