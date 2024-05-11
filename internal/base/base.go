package base

import (
	"context"
	"errors"
	"pgpro2024/internal/entityies"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	q1 = "select id_command, script, description_command from commands where id_command=$1;"
	q2 = "select id_command, script, description_command from commands;"
	q3 = "insert into commands(script, description_command) values ($1, $2) returns id_command;"
	q4 = "insert into log_pids(id_command, os_pid) values ($1, $2) returns id_pid;"
	q5 = "insert into data_pids(data_start) values returns id_pid;"
	q6 = "insert into log_command(id_pid, data_logs, type_log) values ($1, $2, $3) returns id;"
	q7 = "update data_pids set data_finish = $1 where id_pid=$2;"
	q8 = "select data_logs, type_log from log_command where id_pid = $1;"
	q9 = "select id_pid from log_pids where id_command = $1 and os_pid = $2;"
)

type Base struct {
	Pool *pgxpool.Pool
}

func NewBase(uri string) (*Base, error) {
	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, err
	}
	return &Base{Pool: pool}, nil
}

func (b *Base) GetAvailibleCommandById(id int) (entityies.Command, error) {
	rows, err := b.Pool.Query(context.Background(), q1, id)
	if err != nil {
		return entityies.Command{}, err
	}
	defer rows.Close()
	if rows.Next() {
		var cmd entityies.Command
		rows.Scan(&cmd.Id, &cmd.Script, &cmd.Description)
		return cmd, nil
	} else {
		return entityies.Command{}, errors.New("404")
	}
}

func (b *Base) GetListAvailibleCommands() (entityies.Commands, error) {
	rows, err := b.Pool.Query(context.Background(), q2)
	if err != nil {
		return entityies.Commands{}, err
	}
	defer rows.Close()
	var cmds entityies.Commands
	for rows.Next() {
		var cmd entityies.Command
		rows.Scan(&cmd.Id, &cmd.Script, &cmd.Description)
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

func (b *Base) CreateCommand(command entityies.Command) (int, error) {
	rows, err := b.Pool.Query(context.Background(), q3, command.Script, command.Description)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&command.Id)
	}
	return command.Id, nil
}

func (b *Base) StartCommand(start entityies.ProcessStart) (int, error) {
	tx, err := b.Pool.Begin(context.Background())
	if err != nil {
		return -1, err
	}
	defer tx.Commit(context.Background())

	tx.Begin(context.Background())
	rows, err := tx.Query(context.Background(), q4, start.IdCommand, start.Os_pid)
	if err != nil {
		tx.Rollback(context.Background())
		return -1, err
	}
	var id_pid int
	if rows.Next() {
		rows.Scan(&id_pid)
	}
	var builder strings.Builder
	for i := 0; i < len(start.ParametrsStart); i++ {
		builder.WriteString(start.ParametrsStart[i] + " ")
	}
	_, err = tx.Query(context.Background(), q6, id_pid, builder.String(), "INIT")
	if err != nil {
		tx.Rollback(context.Background())
		return -1, err
	}
	_, err = tx.Query(context.Background(), q5, start.DataStart)
	if err != nil {
		tx.Rollback(context.Background())
		return -1, err
	}
	return id_pid, nil
}

func (b *Base) GetLogsProcess(start entityies.ProcessStarted) (entityies.Logs, error) {
	rows, err := b.Pool.Query(context.Background(), q9, start.Id_command, start.Os_pid)
	if err != nil {
		return entityies.Logs{}, err
	}
	var id_pid int
	if rows.Next() {
		rows.Scan(&id_pid)
	}
	rows.Close()
	rows, err = b.Pool.Query(context.Background(), q8, id_pid)
	if err != nil {
		return entityies.Logs{}, err
	}
	defer rows.Close()
	var ans entityies.Logs
	for rows.Next() {
		var lg entityies.LogMessages
		lg.Process = start
		rows.Scan(lg.Message, lg.Stream)
		ans = append(ans, lg)
	}
	return ans, nil
}

func (b *Base) StopProcess(start entityies.ProcessStarted, data time.Time) error {
	rows, err := b.Pool.Query(context.Background(), q9, start.Id_command, start.Os_pid)
	if err != nil {
		return err
	}
	var id_pid int
	if rows.Next() {
		rows.Scan(&id_pid)
	}
	rows.Close()
	rows, err = b.Pool.Query(context.Background(), q7, id_pid)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (b *Base) AdddLog(msg entityies.LogMessages) error {
	rows, err := b.Pool.Query(context.Background(), q9, msg.Process.Id_command, msg.Process.Os_pid)
	if err != nil {
		return err
	}
	var id_pid int
	if rows.Next() {
		rows.Scan(&id_pid)
	}
	rows.Close()
	_, err = b.Pool.Query(context.Background(), q6, id_pid, msg.Message, msg.Stream)
	if err != nil {
		return err
	}
	return nil
}
