package base

import (
	"pgpro2024/internal/entityies"
	"time"
)

type BaseInterface interface {
	GetAvailibleCommandById(id int) (entityies.Command, error)
	GetListAvailibleCommands() (entityies.Commands, error)
	CreateCommand(command entityies.Command) (int, error)
	StartCommand(start entityies.ProcessStart) (int, error)
	GetLogsProcess(start entityies.ProcessStarted) (entityies.Logs, error)
	StopProcess(start entityies.ProcessStarted, data time.Time, code int) error
	AdddLog(msg entityies.LogMessages) error
	GetStatusProcess(start entityies.ProcessStarted) (entityies.ProcessStatus, error)
}
