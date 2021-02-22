package ssh

import (
	"TestDemoApp/service"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/juju/errors"
	"github.com/sirupsen/logrus"
)

const (
	// Порт трасляции SSH данные в WEB-интерфейс
	port = 8080
)

// SSH инициализируется через NewSSH
type SSH struct {
	ctx       context.Context
	port      int
	log       *logrus.Entry
	gottyArgs string
}

// Config конфигурация SSH
type Config struct {
	Port      int
	Log       *logrus.Logger
	GottyArgs string
}

// NewSSH конструктор структуры SSH. В byte получаем бинарник c gotty.
func NewSSH(ctx context.Context, cfg Config) (service.SSHSvc, error) {
	log := logrus.New()
	log.Out = ioutil.Discard
	if cfg.Log != nil {
		log = cfg.Log
	}

	ssh := &SSH{
		ctx:  ctx,
		port: port,
		log: cfg.Log.WithFields(map[string]interface{}{
			"module": "ssh",
			"scope":  "service",
		}),
		gottyArgs: strings.TrimSpace(cfg.GottyArgs),
	}

	if cfg.Port != 0 {
		ssh.port = cfg.Port
	}

	// Сохраняем бинарный файл во временную директорию

	return ssh, nil
}

// Run
func (m SSH) Run() error {
	var pathInGobindata string
	var outFilePath = filepath.Join(os.TempDir(), "gotty")

	// Определяние запускного файла в gobindata
	if runtime.GOOS == "windows" {
		m.log.Warn("обноружена ОС Windows; SSH модуль не загружен")
		//pathInGobindata = "assets/gotty/linux_amd64/gotty"
	} else if runtime.GOOS == "linux" {
		switch runtime.GOARCH {
		case "amd64":
			pathInGobindata = "assets/gotty/linux_amd64/gotty"
		case "arm":
			pathInGobindata = "assets/gotty/linux_arm/gotty"
		}
		m.log.Debug("pathInGobindata: ", pathInGobindata)
	} else {
		m.log.Errorf("неподдерживаемя ОС: %s; SSH модуль не загружен", runtime.GOOS)
	}

	// Распаковываем файл во временную директорию
	if pathInGobindata != "" {
		fileBin, err := Asset(pathInGobindata)
		if err != nil {
			return errors.Trace(err)
		}
		if _, err := os.Stat(outFilePath); os.IsNotExist(err) {
			err = ioutil.WriteFile(outFilePath, fileBin, os.ModePerm)
			if err != nil {
				m.log.Error(err)
				return errors.Trace(err)
			}
			m.log.Info("записан файл: ", outFilePath)
		} else {
			m.log.Infof("файл %s уже сущесвтует в системе", outFilePath)
		}
	}

	// Запуск программы на исполнение
	if pathInGobindata != "" {
		param := []string{"-w", "bash"}
		if m.gottyArgs != "" {
			param = strings.Split(m.gottyArgs, " ")
		}
		m.log.Info("запуск программы: ", outFilePath, " ", strings.Join(param, " "))
		subProcess := exec.Command(outFilePath, param...)

		subProcess.Stdout = os.Stdout
		subProcess.Stderr = os.Stderr

		if err := subProcess.Start(); err != nil {
			m.log.Error(err.Error())
			return errors.Trace(err)
		}

		err := subProcess.Wait()
		if err != nil {
			m.log.Error(err.Error())
			return errors.Trace(err)
		}

	}

	// Ожидаем сигнала завершения работы
	select {
	case <-m.ctx.Done():
		return m.ctx.Err()
	}

}
