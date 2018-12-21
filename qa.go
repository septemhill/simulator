package simulator

import "time"

type AnswerSet struct {
	Answer string
	Count  int
	Time   time.Time
}

type Question struct {
	Ques    string
	Answers []AnswerSet
}

func PickRandomOne() {

}

func PickRandomOneWithoutAnswer() {

}
