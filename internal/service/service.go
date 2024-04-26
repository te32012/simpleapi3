package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"pgpro2024/internal/base"
	"pgpro2024/internal/entityies"
	"time"
)

type Service struct {
	Base   base.BaseInterface
	Proces map[entityies.ProcessStarted]*exec.Cmd
}

func (s *Service) GetAvailibleCommandById(id int) ([]byte, entityies.Error) {
	cmd, err := s.Base.GetAvailibleCommandById(id)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	data, err := json.Marshal(&cmd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	return data, entityies.Error{}
}

func (s *Service) GetListAvailibleCommands() ([]byte, entityies.Error) {
	cmd, err := s.Base.GetListAvailibleCommands()
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	data, err := json.Marshal(&cmd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	return data, entityies.Error{}
}

func (s *Service) CreateCommand(data []byte) ([]byte, entityies.Error) {
	var cmd entityies.Command
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	id, err := s.Base.CreateCommand(cmd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	var cmdc entityies.CommandCreated
	cmdc.Id = id
	ans, err := json.Marshal(&cmdc)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	return ans, entityies.Error{}
}

func (s *Service) StartCommand(id int, data []byte) ([]byte, entityies.Error) {
	cmd, err := s.Base.GetAvailibleCommandById(id)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	var pst entityies.ProcessStart
	err = json.Unmarshal(data, &pst)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	pst.DataStart = time.Now()
	log_id, err := s.Base.StartCommand(pst)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	file, err := os.CreateTemp("", "command")
	info, _ := file.Stat()
	p := info.Name()
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	file.WriteString(cmd.Script)
	file.Close()
	var ok chan bool
	c := exec.Command(p, pst.ParametrsStart...)
	go func() {
		ok <- c.ProcessState.Exited()
		close(ok)
	}()
	writer1 := bufio.NewWriter(nil)
	writer2 := bufio.NewWriter(nil)
	reader1 := bufio.NewReader(nil)
	c.Stdout = writer1
	c.Stderr = writer2
	c.Stdin = reader1
	go s.running(reader1, writer1, writer2, log_id, c, ok)
	err = c.Start()
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	var pstd entityies.ProcessStarted
	pstd.Log_id = log_id
	pstd.Pid = c.Process.Pid
	ans, err := json.Marshal(&pstd)
	return ans, entityies.Error{}
}

func (s *Service) GetStatusProcess(data []byte) ([]byte, entityies.Error) {
	var pstd entityies.ProcessStarted
	err := json.Unmarshal(data, &pstd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	var answerLog entityies.AnswerLog
	answerLog.Stdin, answerLog.Stderr, answerLog.Stdout, err = s.Base.GetLogsProcess(pstd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	_, ok := s.Proces[pstd]
	if ok {
		answerLog.Status = "running"
	} else {
		answerLog.Status = "exited"
	}
	ans, err := json.Marshal(&answerLog)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	return ans, entityies.Error{}
}

func (s *Service) StopProcess(data []byte) entityies.Error {
	var pstd entityies.ProcessStarted
	err := json.Unmarshal(data, &pstd)
	if err != nil {
		return entityies.Error{E: err, Err: []byte(err.Error())}
	}
	c, ok := s.Proces[pstd]
	if ok {
		err = c.Cancel()
		if err != nil {
			return entityies.Error{E: err, Err: []byte(err.Error())}
		}
		return entityies.Error{}
	}
	return entityies.Error{E: errors.New("процесс не найден"), Err: []byte("процесс не найден")}
}

func (s *Service) running(stdin *bufio.Reader, stdout *bufio.Writer, stderr *bufio.Writer, log_id int, c *exec.Cmd, ok chan bool) {
	t2 := time.NewTicker(1 * time.Minute)
	t3 := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-t2.C:
			r := bufio.NewReader(nil)
			stdout.ReadFrom(r)
			var lg entityies.LogMessages
			lg.Data = time.Now()
			lg.Stream = "stdout"
			lg.Process = entityies.ProcessStarted{Log_id: log_id, Pid: c.Process.Pid}
			data, _ := io.ReadAll(r)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
		case <-t3.C:
			r := bufio.NewReader(nil)
			stderr.ReadFrom(r)
			var lg entityies.LogMessages
			lg.Data = time.Now()
			lg.Stream = "stderr"
			lg.Process = entityies.ProcessStarted{Log_id: log_id, Pid: c.Process.Pid}
			data, _ := io.ReadAll(r)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
		case <-ok:
			r := bufio.NewReader(nil)
			stdout.ReadFrom(r)
			var lg entityies.LogMessages
			lg.Data = time.Now()
			lg.Stream = "stdout"
			lg.Process = entityies.ProcessStarted{Log_id: log_id, Pid: c.Process.Pid}
			data, _ := io.ReadAll(r)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
			r = bufio.NewReader(nil)
			stderr.ReadFrom(r)
			lg = entityies.LogMessages{}
			lg.Data = time.Now()
			lg.Stream = "stderr"
			lg.Process = entityies.ProcessStarted{Log_id: log_id, Pid: c.Process.Pid}
			data, _ = io.ReadAll(r)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
			delete(s.Proces, entityies.ProcessStarted{Log_id: log_id, Pid: c.Process.Pid})
			break
		}
	}
}
