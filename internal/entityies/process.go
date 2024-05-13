package entityies

import "time"

type ProcessStatus struct {
	Pid            int      `json:"Pid"`
	ExitCode       int      `json:"ExitCode"`
	Exited         bool     `json:"Exited"`
	Id_logs        int      `json:"Id_logs"`
	ParametrsStart []string `json:"ParametrsStart"`
}

type ProcessStart struct {
	IdCommand      int `json:"IdCommand"`
	Os_pid         int
	ParametrsStart []string `json:"ParametrsStart"`
	DataStart      time.Time
	InputStream    string `json:"InputStream"`
}

type ProcessStarted struct {
	Os_pid     int `json:"Os_pid"`
	Id_command int `json:"Id_command"`
}
