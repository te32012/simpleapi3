package service

import "pgpro2024/internal/entityies"

type ServiceInterface interface {
	GetAvailibleCommandById(id int) ([]byte, entityies.Error)
	GetListAvailibleCommands() ([]byte, entityies.Error)
	CreateCommand(data []byte) ([]byte, entityies.Error)
	StartCommand(data []byte) ([]byte, entityies.Error)
	GetStatusProcess(data []byte) ([]byte, entityies.Error)
	StopProcess(data []byte) entityies.Error
}
