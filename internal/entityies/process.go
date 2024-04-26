package entityies

import "time"

type ProcessStatus struct {
	Pid            int
	ExitCode       int
	Exited         bool
	Id_logs        int
	ParametrsStart []string
}

type ProcessStart struct {
	IdCommand      int
	ParametrsStart []string
	DataStart      time.Time
}

type ProcessStarted struct {
	Pid    int
	Log_id int
}
