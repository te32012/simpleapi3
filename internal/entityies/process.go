package entityies

import "time"

type ProcessStatus struct {
	Pid        int        `json:"Pid"`
	ExitCode   *int       `json:"ExitCode,omitempty"`
	Id_logs    int        `json:"Id_logs"`
	IdCommand  int        `json:"IdCommand"`
	DataStart  time.Time  `json:"DataStart"`
	DataFinish *time.Time `json:"DataFinish,omitempty"`
}

type ProcessStart struct {
	IdCommand      int `json:"IdCommand"`
	Os_pid         int
	ParametrsStart []string `json:"ParametrsStart"`
	DataStart      time.Time
	InputStream    string `json:"InputStream"`
}

type ProcessStarted struct {
	Os_pid  int `json:"Os_pid"`
	Id_logs int `json:"Id_logs"`
}
