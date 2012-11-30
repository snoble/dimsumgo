package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/big"
	. "os"
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
	inputLines := make(chan []byte)

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

	go readInput(input, inputLines)

	for line := range inputLines {

		if position < int64(sampleSize) {
			out[position] = line

		} else {
			r, err := rand.Int(rand.Reader, big.NewInt(position))

			if err != nil {
				panic(err)
			}

			if r.Int64() < int64(sampleSize) {
				out[r.Int64()] = line
			}
		}

		position = position + 1
	}

	for _, line := range out {
		fmt.Printf("%s", line)
	}
	// CLose input
}
