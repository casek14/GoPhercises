package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	DEFAULTFILEPATH = "/home/casek/go/src/github.com/casek14/GoPhercises/QuizGame/problem.csv"
	NOTIMEMODE      = "NOTIMEMODE"
	TIMEMODE        = "TIMEMODE"
	ENDOFCOUNTDOWN  = "END"
	SHUFFLE_REVERT =  "REVERT"
)

var filePath string
var gameMode string
var countdown int
var shuffleType string

// Question struct represents a question, which consists of question and answer.
type Question struct {
	Question string
	Answer   string
}

// Quiz struct represents quiz, which holds list of all questions.
type Quiz struct {
	Questions []Question
}

// ReadQuestions function reads file,
// parse this file and return Quiz struct.
func ReadQuestions(file string) Quiz {

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var quiz Quiz
	for _, line := range lines {
		quiz.Questions = append(quiz.Questions, Question{Question: line[0], Answer: line[1]})
	}

	return quiz
}

// init flags
func Init() {
	flag.StringVar(&filePath, "questions", DEFAULTFILEPATH, "Specify path to the questions file.")
	flag.StringVar(&gameMode, "mode", NOTIMEMODE, "For game without time use NOTIMEMODE, for time game use TIMEMODE.")
	flag.IntVar(&countdown, "time", 10, "Set time for whole test.")
	flag.StringVar(&shuffleType, "shuffle", SHUFFLE_REVERT,"Shuffle questions. Possible options are REVERT, RANDOM")
	flag.Parse()
}

func countdownCounter(t int, c chan string){
	log.Printf("Starting %d countdown !!\n",t)
	time.Sleep(time.Duration(t) * time.Second)
	c <- ENDOFCOUNTDOWN
}

func shuffleQuiz(quiz Quiz) []Question{
	rand.Seed(time.Now().UnixNano())
	pole := make([]Question,len(quiz.Questions),len(quiz.Questions))

	for i := 0; i < len(quiz.Questions);i++{
		newNumber:
		pos := rand.Intn(len(quiz.Questions))
		if quiz.Questions[pos].Question != ""{
			pole[i] = quiz.Questions[pos]
			quiz.Questions[pos] = Question{Question:"", Answer:""}
		}else {
			goto newNumber
		}
	}
	return pole
}

func main() {
	Init()

	quiz := ReadQuestions(filePath)
	quiz.Questions = shuffleQuiz(quiz)
	var znak string
	correctAnswers := 0
	question := 0
	fmt.Println("***** Zacina test *****")
	if gameMode == NOTIMEMODE {

		for c, otazka := range quiz.Questions {
			fmt.Println("====	")
			fmt.Printf("Question number %d) - %s\n", c+1, otazka.Question)
			n, err := fmt.Scanln(&znak)
			if err != nil {
				log.Fatal("CANNOT READ FROM CONSOLE ", n)
			}
			if znak == otazka.Answer {
				fmt.Println("Correct answer.")
				correctAnswers++

			} else {
				fmt.Printf("Wrong answer. Should be %s.\n", otazka.Answer)
			}
			question++
		}
	} else if gameMode == TIMEMODE && countdown > 0 {
		timer := time.NewTimer(time.Duration(countdown) * time.Second)
		kanal := make(chan string)
	problemloop:
		for c, otazka := range quiz.Questions {
			fmt.Println("====	")
			fmt.Printf("Question number %d) - %s\n", c+1, otazka.Question)
			go func() {
				var odpoved string
				fmt.Scanf("%s\n",&odpoved)
				kanal <-odpoved
			}()
				select {
				case <- timer.C:
					fmt.Println("COUNTDOWN !!!")
					break problemloop
				case answer := <-kanal:
				if answer == otazka.Answer {
					fmt.Println("Correct answer.")
					correctAnswers++
				} else {
					fmt.Printf("Wrong answer. Should be %s.\n", otazka.Answer)
				}
				}
question++
		}

	} else {
		log.Fatalf("Game mode: %s is not supported please choose %s or %s game mode. "+
			"Or time must be greated than 0.",
			gameMode, NOTIMEMODE, TIMEMODE)
	}
		var percentage float32
		fmt.Printf("Number of correct answers: %d\n", correctAnswers)
		percentage = (float32(correctAnswers) / float32(len(quiz.Questions)) * 100)

		fmt.Printf("Test success percentage: %d%.\n", int(percentage))


}
