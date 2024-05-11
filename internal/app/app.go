package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"pgpro2024/internal/base"
	"pgpro2024/internal/middleware"
	"pgpro2024/internal/service"
	"syscall"
)

func Run() {
	b, err := base.NewBase(fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("POSTGRESQL_USER"), os.Getenv("POSTGRESQL_PASSWORD"), os.Getenv("POSTGRESQL_HOST"), os.Getenv("POSTGRESQL_PORT"), os.Getenv("POSTGRESQL_BASE")))
	if err != nil {
		slog.Error("ошибка соединения с базой %s", err.Error())
	}
	s := service.NewService(b)
	r := middleware.NewMyRouter("pgpro", "2024", s)

	go func() {
		r.ListenAndServe()
	}()

	slog.Info("success starting application")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	slog.Info("application has been shut down")

}
