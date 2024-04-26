package entityies

import "time"

type Stderr []LogMessages
type Stdin []LogMessages
type Stdout []LogMessages

type LogMessages struct {
	Process ProcessStarted
	Stream  string
	Message string
	Data    time.Time
}

type AnswerLog struct {
	Stdin  Stdin
	Stdout Stdout
	Stderr Stderr
	Status string
}
