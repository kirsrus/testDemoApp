package main

import (
	"context"
	"flag"
	"fmt"
	stdLog "log"
	"os"
	"runtime"

	"TestDemoApp/pkg/logger"
	"TestDemoApp/service/ssh"
	"TestDemoApp/service/web"

	"github.com/juju/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	webPort  = 80
	logLevel = logrus.DebugLevel
)

//go:generate go-bindata -o ../bindata/bindata.go -pkg bindata -fs -prefix "../" ../assets/...
func main() {
	err := run()
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
		fmt.Println("StackTrace:")
		stdLog.Fatal(errors.ErrorStack(err))
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())

	// region Входные параметры аргументов

	webPortArg := flag.Int("port", webPort, "порт WEB-сервера")
	logLevelArg := flag.String("level", logLevel.String(), "уровень логирования (debug|info|warn|error)")
	gottyArgs := flag.String("gottyargs", "-w bash", "параметры для gotty")
	helpArg := flag.Bool("help", false, "помощь")

	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.Parse()
	if *helpArg {
		flag.Usage()
		return nil
	}

	// endregion
	// region Настройка логирования

	log := logrus.New()
	pLogLevel, err := logrus.ParseLevel(*logLevelArg)
	if err != nil {
		pLogLevel = logLevel
	}
	log.Level = pLogLevel
	log.Formatter = &logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006.01.02 15:04:05",
		ForceColors:     true,
	}
	log.Out = os.Stdout
	log.AddHook(logger.LogrusContextHook{})

	tmpLevel := log.Level
	log.Level = logrus.InfoLevel
	log.Printf("config: logLevel = %s", pLogLevel)
	log.Println("runtime.GOARCH:", runtime.GOARCH)       // amd64
	log.Println("runtime.GOOS:", runtime.GOOS)           // windows
	log.Println("runtime.Version():", runtime.Version()) // go1.15.6
	log.Println("runtime.GOROOT():", runtime.GOROOT())   // c:\Program Files\Go
	log.Println("os.TempDir():", os.TempDir())
	log.Level = tmpLevel

	// endregion
	// region Настройка WEB-сервера

	webService, err := web.NewWeb(ctx, web.Config{
		Log:  log,
		Port: *webPortArg,
	})
	if err != nil {
		return errors.Trace(err)
	}

	// todo: подумать как реализовать извне
	//webService.Login("/login")
	//webService.Static("/")
	//webService.Finish()

	// endregion
	// region Настройка SSH

	sshService, err := ssh.NewSSH(ctx, ssh.Config{
		Log:       log,
		GottyArgs: *gottyArgs,
	})

	// endregion
	// region Запуск сервисов

	g, gctx := errgroup.WithContext(ctx)

	// WEB-сервер
	g.Go(func() error {
		end := make(chan error)

		go func() {
			end <- webService.Run()
		}()

		select {
		case err := <-end:
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		case <-gctx.Done():
			err := webService.Stop()
			if err != nil {
				log.Error(err)
			}
			return gctx.Err()
		}
	})

	// SSH-сервер
	g.Go(func() error {
		end := make(chan error)

		go func() {
			end <- sshService.Run()
		}()

		select {
		case err := <-end:
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		case <-gctx.Done():
			return gctx.Err()
		}

	})

	// endregion

	if err = g.Wait(); err != nil {
		cancel()
		if err == context.Canceled {
			return nil
		}
		return errors.Trace(err)
	}

	return nil
}
