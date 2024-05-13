package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"pgpro2024/internal/base"
	"pgpro2024/internal/entityies"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	Base   base.BaseInterface
	Proces *sync.Map
}

func NewService(base base.BaseInterface) *Service {
	return &Service{Base: base, Proces: new(sync.Map)}
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
	for i := 0; i < len(cmd); i++ {
		cmd[i].Script = ""
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

func (s *Service) StartCommand(data []byte) ([]byte, entityies.Error) {
	var pst entityies.ProcessStart
	err := json.Unmarshal(data, &pst)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	pst.DataStart = time.Now()
	cmd, err := s.Base.GetAvailibleCommandById(pst.IdCommand)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	file, err := os.CreateTemp("/tmp", strconv.Itoa(pst.IdCommand)+"*")
	info, _ := file.Stat()
	p := info.Name()
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	file.WriteString(cmd.Script)
	file.Close()
	ok := make(chan bool)
	var args []string
	args = append(args, "/tmp/"+p)
	args = append(args, pst.ParametrsStart...)
	c := exec.Command("bash", args...)
	writer1 := bytes.NewBuffer([]byte{})
	writer2 := bytes.NewBuffer([]byte{})
	reader1 := bytes.NewBuffer([]byte{})
	reader1.WriteString(pst.InputStream)
	c.Stdout = writer1
	c.Stderr = writer2
	c.Stdin = reader1
	log_id := make(chan int)
	go s.running(writer1, writer2, pst, c, ok, log_id)
	err = c.Start()
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	go func(ok1 chan bool) {
		fmt.Println("cheking")
		c.Wait()
		fmt.Println("cheked")
		ok1 <- true
		time.Sleep(5 * time.Second)
		close(ok1)
	}(ok)
	var pstd entityies.ProcessStarted
	pstd.Id_logs = <-log_id
	pstd.Os_pid = c.Process.Pid
	s.Proces.LoadOrStore(pstd, c)
	ans, err := json.Marshal(&pstd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}

	}
	return ans, entityies.Error{}
}

func (s *Service) GetStatusProcess(data []byte) ([]byte, entityies.Error) {
	var pstd entityies.ProcessStarted
	err := json.Unmarshal(data, &pstd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	var answerLog entityies.AnswerLog
	answerLog.Logs, err = s.Base.GetLogsProcess(pstd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
	}
	answerLog.ProcessStatus, err = s.Base.GetStatusProcess(pstd)
	if err != nil {
		return nil, entityies.Error{E: err, Err: []byte(err.Error())}
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
	c, ok := s.Proces.Load(pstd)
	if !ok {
		return entityies.Error{E: errors.New("404"), Err: []byte("процесс не найден")}
	}
	if cmd, ok2 := c.(*exec.Cmd); ok2 && ok {
		err = cmd.Process.Kill()
		if err != nil {
			return entityies.Error{E: err, Err: []byte(err.Error())}
		}
		return entityies.Error{}
	}
	return entityies.Error{E: errors.New("какая-то ошибка"), Err: []byte("какая-то ошибка")}
}

func (s *Service) running(stdout *bytes.Buffer, stderr *bytes.Buffer, pst entityies.ProcessStart, c *exec.Cmd, ok chan bool, lid chan int) {
	t2 := time.NewTicker(1 * time.Minute)
	t3 := time.NewTicker(1 * time.Minute)
	for c.Process == nil {
		time.Sleep(1 * time.Microsecond)
	}
	fmt.Println(c.Process.Pid)
	pst.Os_pid = c.Process.Pid
	log_id, _ := s.Base.StartCommand(pst)
	lid <- log_id
	close(lid)
	for {
		select {
		case <-t2.C:
			var lg entityies.LogMessages
			lg.Stream = "stdout"
			lg.Process = entityies.ProcessStarted{Id_logs: log_id, Os_pid: c.Process.Pid}
			data, _ := io.ReadAll(stdout)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
		case <-t3.C:
			var lg entityies.LogMessages
			lg.Stream = "stderr"
			lg.Process = entityies.ProcessStarted{Id_logs: log_id, Os_pid: c.Process.Pid}
			data, _ := io.ReadAll(stderr)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
		case <-ok:
			var lg entityies.LogMessages
			lg.Stream = "stdout"
			lg.Process = entityies.ProcessStarted{Id_logs: log_id, Os_pid: c.Process.Pid}
			data, _ := io.ReadAll(stdout)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
			lg = entityies.LogMessages{}
			lg.Stream = "stderr"
			lg.Process = entityies.ProcessStarted{Id_logs: log_id, Os_pid: c.Process.Pid}
			data, _ = io.ReadAll(stderr)
			if len(data) > 0 {
				lg.Message = string(data[:])
				s.Base.AdddLog(lg)
			}
			s.Proces.Delete(entityies.ProcessStarted{Id_logs: log_id, Os_pid: c.Process.Pid})
			s.Base.StopProcess(entityies.ProcessStarted{Id_logs: log_id, Os_pid: c.Process.Pid}, time.Now(), c.ProcessState.ExitCode())
			return
		}
	}
}
