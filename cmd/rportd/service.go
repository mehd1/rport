package main

import (
	"log"
	"path/filepath"

	"github.com/kardianos/service"

	chserver "github.com/cloudradar-monitoring/rport/server"
	chshare "github.com/cloudradar-monitoring/rport/share"
)

var svcConfig = &service.Config{
	Name:        "rportd",
	DisplayName: "CloudRadar Rport Server",
	Description: "Create reverse tunnels with ease.",
}

func handleSvcCommand(svcCommand string, configPath string) error {
	svc, err := getService(nil, configPath)
	if err != nil {
		return err
	}

	return chshare.HandleServiceCommand(svc, svcCommand)
}

func runAsService(s *chserver.Server, configPath string) error {
	svc, err := getService(s, configPath)
	if err != nil {
		return err
	}

	return svc.Run()
}

func getService(s *chserver.Server, configPath string) (service.Service, error) {
	if configPath != "" {
		absConfigPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, err
		}
		svcConfig.Arguments = []string{"-c", absConfigPath}
	}
	return service.New(&serviceWrapper{s}, svcConfig)
}

type serviceWrapper struct {
	*chserver.Server
}

func (w *serviceWrapper) Start(service.Service) error {
	if w.Server == nil {
		return nil
	}
	go func() {
		if err := w.Server.Run(); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (w *serviceWrapper) Stop(service.Service) error {
	return w.Server.Close()
}
