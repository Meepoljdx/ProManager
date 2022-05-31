package command

import (
	"ProManager/module/remote"
	log "ProManager/module/utils/logs"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (s *Service) OperateService(client *ssh.Client, operation string) bool {
	script := ""
	switch s.Name {
	case "prometheus":
		script = "sbin/prometheus.sh"
	case "node-exporter":
		script = "sbin/node-exporter.sh"
	case "datanode-exporter", "namenode-exporter":
		script = "sbin/hdfs_exporter.sh"
	case "nodemanager-exporter", "resourcemanager-exporter":
		script = "sbin/yarn_exporter.sh"
	default:
		script = "sbin/pro_manager.sh"
	}
	Cmd := fmt.Sprintf("%s/%s/%s %s", s.InstallPath, s.Name, script, operation) // 调用shell脚本的对应入口
	output, err := remote.CommandExec(client, Cmd)
	if err != nil {
		log.Logger.Errorf("%s\n", output)
		log.Logger.Errorf("Service %s %s failed, %s. Command: %s\n", s.Name, operation, err, Cmd)
		return false
	}
	log.Logger.Infof("Service %s %s success.\n%s", s.Name, operation, output)
	return true
}

func (s *Service) StartService(port string, user string, password string) bool { // 服务启动
	client, err := CreateClient(s.IP, port, user, password)
	if err != nil {
		log.Logger.Errorln(err)
		return false
	}
	isSuccess := s.OperateService(client, "start")
	defer client.Close()
	return isSuccess
}

func (s *Service) StopService(port string, user string, password string) bool { // 服务停止
	client, err := CreateClient(s.IP, port, user, password)
	if err != nil {
		log.Logger.Errorln(err)
		return false
	}
	isSuccess := s.OperateService(client, "stop")
	defer client.Close()
	return isSuccess
}

func (s *Service) RestartService(port string, user string, password string) bool { // 服务重启
	client, err := CreateClient(s.IP, port, user, password)
	if err != nil {
		log.Logger.Errorln(err)
		return false
	}
	isSuccess := s.OperateService(client, "restart")
	defer client.Close()
	return isSuccess
}

func (s *Service) ReloadService(port string, user string, password string) bool { // 服务重载
	client, err := CreateClient(s.IP, port, user, password)
	if err != nil {
		log.Logger.Errorln(err)
		return false
	}
	isSuccess := s.OperateService(client, "reload")
	defer client.Close()
	return isSuccess
}
