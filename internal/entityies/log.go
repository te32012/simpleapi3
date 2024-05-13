package entityies

type Logs []LogMessages
type LogMessages struct {
	Process ProcessStarted `json:"Process"`
	Stream  string         `json:"Stream"`
	Message string         `json:"Message"`
}

type AnswerLog struct {
	Logs          Logs          `json:"Logs"`
	ProcessStatus ProcessStatus `json:"ProcessStatus"`
}
