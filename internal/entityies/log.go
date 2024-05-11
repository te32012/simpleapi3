package entityies

type Logs []LogMessages
type LogMessages struct {
	Process ProcessStarted
	Stream  string
	Message string
}

type AnswerLog struct {
	Logs   Logs
	Status string
}
