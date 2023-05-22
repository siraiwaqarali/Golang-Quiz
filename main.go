package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func main() {
	// 1. Input the name of the file
	fName := flag.String("f", "quiz.csv", "Path of csv file")
	// 2. Set the duration of the timer
	timer := flag.Int("t", 30, "Duration of the quiz in seconds")
	flag.Parse()
	// 3. Pull the problems from the file (calling our problem puller function)
	problems, err := problemPuller(*fName)
	// 4. Handle the error
	if err != nil {
		exit(fmt.Sprintf("Something Went Wrong: %s", err.Error()))
	}
	// 5. Create a variable to count our correct answers
	correctAns := 0
	// 6. Using the duration of the timer, we want to initialize a timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansChan := make(chan string)

	// 7. Loop through the problems, print the questions, we'll accept the answers
problemsLoop: // Label for the loop
	for i, problem := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, problem.question)

		go func() {
			fmt.Scanf("%s", &answer)
			ansChan <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemsLoop
		case iAns := <-ansChan:
			if iAns == problem.answer {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansChan)
			}
		}
	}

	// 8. We will calculate and print out the result
	fmt.Printf("You scored %d out of %d.\n", correctAns, len(problems))
	fmt.Print("Thank you for playing! Press Enter to exit.")
	<-ansChan
}

func parseProblems(lines [][]string) []Problem {
	// Go over the lines and parse them, with Problem struct
	problems := make([]Problem, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		problems[i] = Problem{question: line[0], answer: line[1]}
	}
	return problems
}

func problemPuller(fileName string) ([]Problem, error) {
	// Read all problems from the file quiz.csv
	// 1. Open the file and defer closing the file
	fileObj, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fileObj.Close()
	// 2. Create a new reader
	csvReader := csv.NewReader(fileObj)
	// 3. Read all the lines from the file
	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	// 4. Parse the lines into a slice of Problem struct by calling parseProblems
	return parseProblems(lines), nil
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
