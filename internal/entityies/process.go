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
	Os_pid         int
	ParametrsStart []string
	DataStart      time.Time
}

type ProcessStarted struct {
	Os_pid     int
	Id_command int
}
