// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/getlantern/errors"
	"github.com/someanon/rkn-bypasser/proxy"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultAddr = "127.0.1.1:8000"

	serviceName        = "rkn-bypasser"
	serviceDisplayName = "RKN bypasser"
	serviceDescription = "Прокси сервер для обхода блокировок Роскомнадзора"
)

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func installService(c *cli.Context) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("fail to connect to service manager: %s", err.Error())
	}
	defer m.Disconnect()
	s, err := m.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", serviceName)
	}
	exePath, err := exePath()
	if err != nil {
		return fmt.Errorf("fail to get exe path: %s", err.Error())
	}
	args := []string{"service", "run", "-addr"}
	if c.IsSet("addr") {
		args = append(args, c.String("addr"))
	} else {
		args = append(args, defaultAddr)
	}
	if c.IsSet("tor") {
		args = append(args, c.String("tor"))
	}
	s, err = m.CreateService(serviceName, exePath, mgr.Config{
		DisplayName: serviceDisplayName,
		StartType:   mgr.StartAutomatic,
		Description: serviceDescription,
	}, args...)
	if err != nil {
		return fmt.Errorf("fail to create service: %s", err.Error())
	}
	if err := s.Start(); err != nil {
		return fmt.Errorf("fail to start service: %s", err.Error())
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("fail to setup eventlog: %s", err.Error())
	}
	return nil
}

func uninstallService(c *cli.Context) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("fail to connect to service manager: %s", err.Error())
	}
	defer m.Disconnect()
	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", serviceName)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return fmt.Errorf("fail to delete service: %s", err.Error())
	}
	err = eventlog.Remove(serviceName)
	if err != nil {
		return fmt.Errorf("fail to remove eventlog: %s", err.Error())
	}
	return nil
}

func startService(c *cli.Context) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("fail to connect to service manager: %s", err.Error())
	}
	defer m.Disconnect()
	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", serviceName)
	}
	defer s.Close()
	err = s.Start()
	if err != nil {
		return fmt.Errorf("fail to start service: %s", err.Error())
	}
	return nil
}

func runService(c *cli.Context) error {
	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		return errors.New("fail to determine if we are running in an interactive session: %s", err.Error())
	}
	if isIntSess {
		return errors.New("command not intended to be used in interactive session")
	}

	if !c.IsSet("addr") {
		return errors.New("addr is not set")
	}

	addr := c.String("addr")
	torAddr := c.String("tor")

	if torAddr == "" && c.Bool("with-tor") {
		torAddr = "tor-proxy:9150"
	}

	elog, err := eventlog.Open(serviceName)
	if err != nil {
		return fmt.Errorf("fail to open eventlog: %s", err.Error())
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", serviceName))

	err = svc.Run(serviceName, &service{addr: addr, torAddr: torAddr})
	if err != nil {
		errStr := fmt.Sprintf("fail to run service %s: %s", serviceName, err.Error())
		elog.Error(1, errStr)
		return fmt.Errorf("fail to run service %s: %s", serviceName, err.Error())
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", serviceName))
	return nil
}

func stopService(c *cli.Context) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("fail to connect to service manager: %s", err.Error())
	}
	defer m.Disconnect()
	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("fail to open service: %s", err.Error())
	}
	defer s.Close()
	status, err := s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %s", c, err.Error())
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", svc.Stopped)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %s", err.Error())
		}
	}
	return nil
}

type service struct {
	addr    string
	torAddr string
}

func (m *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	go proxy.Run(m.addr, m.torAddr)
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop:
				break loop
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
