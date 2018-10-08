package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// (ex) input.data
/*
5, , , ,1, ,8, ,
 , , , , , ,9, ,
 , , , ,4,7, , ,
 ,9, , , , , ,4,
2, , ,1,5, , ,9,
 , ,6, , , , ,2,8
 , , , ,7,9,6, ,
 , ,9, , , , , ,
 , ,8,3, ,4,7, ,
*/
const inputFileName = "input.data"
const oneSideLength = 9
const multipliedSideLength = oneSideLength * oneSideLength
const tempInputFile = "temp_input.cnf"   // TODO: dynamically
const tempOutputFIle = "temp_output.cnf" // TODO: dynamically
var squareRootSideLength = int(math.Sqrt(oneSideLength))

func main() {
	preProcess()
	writeConstraints()
	solveConstraints()
	outputResults()
	postProcess()
}

func preProcess() {
	intSquareRootSideLength := int(math.Sqrt(oneSideLength))
	// math.Sqrt(oneSideLength) must be int
	if float64(intSquareRootSideLength) != math.Sqrt(oneSideLength) {
		fmt.Println("invalid parameter!")
		os.Exit(1)
	}
}

func writeConstraints() {
	inputFile, err := os.OpenFile(tempInputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer inputFile.Close()
	writeInputConstraints(inputFile)
	writeNumberConstraints(inputFile)
	writeHorizontalConstraints(inputFile)
	writeVerticalConstraints(inputFile)
	writeSquareConstraints(inputFile)
}

func writeInputConstraints(inputFile *os.File) {
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	nowLine := 0
	for scanner.Scan() {
		writer := csv.NewReader(strings.NewReader(scanner.Text()))
		writer.Comma = ','
		writer.Comment = '#'
		line, err := writer.Read()
		if len(line) != oneSideLength {
			err = errors.New("input file is invalid")
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		nowColumn := 0
		for _, v := range line {
			if v == " " || v == "_" {
				nowColumn++
				continue
			}
			val, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			constVar := nowLine*multipliedSideLength + nowColumn*oneSideLength + val
			fmt.Fprintln(inputFile, fmt.Sprintf("%d 0", constVar))
			nowColumn++
		}
		nowLine++
	}
}

func writeNumberConstraints(inputFile *os.File) {
	for _, line := range arrayRange(1, oneSideLength) {
		for _, column := range arrayRange(1, oneSideLength) {
			constraints := make([]string, oneSideLength)
			for _, number := range arrayRange(1, oneSideLength) {
				constraints[number-1] = strconv.Itoa(number + (column-1)*oneSideLength + (line-1)*multipliedSideLength)
			}
			joinedNums := strings.Join(constraints, " ")
			fmt.Fprintln(inputFile, fmt.Sprintf("%s 0", joinedNums))
			writeExclusiveConstraints(inputFile, constraints)
		}
	}
}

func writeHorizontalConstraints(inputFile *os.File) {
	for _, line := range arrayRange(1, oneSideLength) {
		for _, number := range arrayRange(1, oneSideLength) {
			constraints := make([]string, oneSideLength)
			for column := 1; column <= oneSideLength; column++ {
				constraints[column-1] = strconv.Itoa(number + (column-1)*oneSideLength + (line-1)*multipliedSideLength)
			}
			joinedNums := strings.Join(constraints, " ")
			fmt.Fprintln(inputFile, fmt.Sprintf("%s 0", joinedNums))
			writeExclusiveConstraints(inputFile, constraints)
		}
	}
}

func writeVerticalConstraints(inputFile *os.File) {
	for _, column := range arrayRange(1, oneSideLength) {
		for _, number := range arrayRange(1, oneSideLength) {
			constraints := make([]string, oneSideLength)
			for line := 1; line <= oneSideLength; line++ {
				constraints[line-1] = strconv.Itoa(number + (column-1)*oneSideLength + (line-1)*multipliedSideLength)
			}
			joinedNums := strings.Join(constraints, " ")
			fmt.Fprintln(inputFile, fmt.Sprintf("%s 0", joinedNums))
			writeExclusiveConstraints(inputFile, constraints)
		}
	}
}

func writeSquareConstraints(inputFile *os.File) {
	for line := 1; line <= 1+(squareRootSideLength-1)*squareRootSideLength; line += squareRootSideLength {
		for column := 1; column <= 1+(squareRootSideLength-1)*squareRootSideLength; column += squareRootSideLength {
			for _, number := range arrayRange(1, oneSideLength) {
				constraints := make([]string, oneSideLength)
				for _, sqLine := range arrayRange(1, squareRootSideLength) {
					for _, sqColumn := range arrayRange(1, squareRootSideLength) {
						constraints[(sqLine-1)*squareRootSideLength+sqColumn-1] = strconv.Itoa(number + (column-1+sqColumn-1)*oneSideLength + (line-1+sqLine-1)*multipliedSideLength)
					}
				}
				joinedNums := strings.Join(constraints, " ")
				fmt.Fprintln(inputFile, fmt.Sprintf("%s 0", joinedNums))
				writeExclusiveConstraints(inputFile, constraints)
			}
		}
	}
}

func writeExclusiveConstraints(inputFile *os.File, arr []string) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			fmt.Fprintln(inputFile, fmt.Sprintf("-%s -%s 0", arr[i], arr[j]))
		}
	}
}

// Returns an array of elements from start to end, inclusive.
// like range(start, end) @ PHP, [from..to] @ haskell...
func arrayRange(start int, end int) []int {
	if start > end {
		return nil
	}
	arr := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		arr[i-start] = i
	}
	return arr
}

func solveConstraints() {
	// difficult to use...
	exec.Command("minisat", tempInputFile, tempOutputFIle).Run()
}

func outputResults() {
	file, err := os.Open(tempOutputFIle)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "SAT" {
			// If it is SAT, answer exists @ line 2 (see https://dwheeler.com/essays/minisat-user-guide.html)
			scanner.Scan()
			answerCnf := scanner.Text()
			answerCnfSplitted := strings.Split(answerCnf, " ")
			answer := make([]string, oneSideLength)
			for _, line := range arrayRange(1, oneSideLength) {
				answerLine := make([]string, oneSideLength)
				for _, column := range arrayRange(1, oneSideLength) {
					for _, index := range arrayRange(1, oneSideLength) {
						ansNum, err := strconv.Atoi(answerCnfSplitted[(line-1)*multipliedSideLength+(column-1)*oneSideLength+index-1])
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						if ansNum > 0 {
							if ansNum%oneSideLength == 0 {
								answerLine[column-1] = strconv.Itoa(oneSideLength)
							} else {
								answerLine[column-1] = strconv.Itoa(ansNum % oneSideLength)
							}
							break
						}
					}
				}
				answer[line-1] = strings.Join(answerLine, ",") + "\n"
			}
			fmt.Println(answer)
		} else if line == "UNSAT" {
			fmt.Println("This Problem can not be solved")
		} else {
			fmt.Println("There is something wrong")
		}
	}
}

func postProcess() {
	if err := os.Remove(tempInputFile); err != nil {
		fmt.Println(err)
	}
	if err := os.Remove(tempOutputFIle); err != nil {
		fmt.Println(err)
	}
}
