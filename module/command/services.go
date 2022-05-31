// @Title: install.go
// @Description: 用于服务安装的部分，设想中这部分需要给定对应的服务名称
// @Author: 李嘉栋
package command

import (
	"ProManager/module/remote"
	log "ProManager/module/utils/logs"
	"ProManager/setting"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
)

type Service struct {
	Name        string // 服务名称
	IP          string
	InstallPath string // 安装路径
	Version     string // 服务版本
	PackageName string // 软件包名称
	User        string // 运行用户
	Conf        string // 临时的配置文件名称
	// Conf        *interface{} // 配置，认为是一个接口
}

type PrometheusConf struct {
	LogDir       string `json:"prometheus_log_dir"`
	PidDir       string `json:"prometheus_pid_dir"`
	ServerPort   string `json:"prometheus_server_port"`
	DataDir      string `json:"prometheus_data_dir"`
	ConfPath     string `json:"prometheus_conf_path"`
	DataDay      string `json:"prometheus_lifecycle"`
	WebLifecycle string `json:"prometheus_web_lifecycle_enable"`
	WebAdminAPI  string `json:"prometheus_web_adminapi_enable"`
}

type NodeExporterConf struct {
	LogDir     string `json:"node_log_dir"`
	PidDir     string `json:"node_pid_dir"`
	ServerPort string `json:"node_server_port"`
}

func CreateInstallPath(client *ssh.Client, s *Service) bool {
	log.Logger.Infof("Begin create the install path %s.\n", s.InstallPath)
	cmd := fmt.Sprintf("ls %s", s.InstallPath)
	output, err := remote.CommandExec(client, cmd)
	if err != nil {
		if strings.Contains(output, "No such file or directory") {
			log.Logger.Infof("The install path %s is not exist, will exist it.\n", s.InstallPath) // 如果路径不存在，则对其进行创建
			// 创建一个目录
			cmd = fmt.Sprintf("mkdir -p %s", s.InstallPath)
			output, err = remote.CommandExec(client, cmd)
			if err != nil {
				log.Logger.Errorf("Create the install path %s failed, %s\n", s.InstallPath, err.Error())
				return false
			}
			log.Logger.Infof("Create the install path %s success.\n", s.InstallPath)
		} else if strings.Contains(output, "Permission denied") {
			log.Logger.Errorf("The install path %s is exist,but the user has no permission to access the install path.\n", cmd) // 若用户没有权限，则退出安装过程
			return false
		}
	} else { // 如果获取路径没有报错，判断文件数量，大于0则退出
		if len(output) > 0 {
			log.Logger.Warnf("The install path %s is not empty, Please check it.\n", s.InstallPath)
			return false
		}
	}
	return true
}

func (s *Service) InitInstallMsg() bool {
	log.Logger.Infof("Begin to init the install message.\n")
	s.User = "prometheus"
	switch s.Name {
	case "prometheus":
		s.InstallPath = "/opt/ProManager"
		s.PackageName = "prometheus.tar.gz"
		s.Version = "Linux-Prometheus-2.34.0"
	case "node-exporter":
		s.InstallPath = "/opt/ProManager/exporters/bases/"
		s.PackageName = "node-exporter.tar.gz"
		s.Version = "Linux-node_exporter-1.3.3"
	case "datanode-exporter":
		s.InstallPath = "/opt/ProManager/exporters/hadoop/"
		s.PackageName = "datanode_exporter.tar.gz"
		s.Version = "Linux-1.0.0"
	default:
		log.Logger.Errorf("Init the install message failed.\n")
		return false
	}
	log.Logger.Infof("Init the install message finished.\n")
	return true
}

func (s *Service) UploadInstallPackages(local string, remotes string, ssh remote.SSHConfig) bool {
	log.Logger.Infof("Begin to upload the install package %s.\n", s.PackageName)
	err := remote.CopyToRemote(local, remotes, &ssh)
	if err != nil {
		log.Logger.Errorf("Upload the install package failed. %s\n", err)
		return false
	}
	log.Logger.Infof("Upload the install package success.\n")
	return true
}

func (s *Service) DecompressInstallPackages(remotes string, client *ssh.Client) bool { // 解压安装包
	log.Logger.Infof("Begin to decompress the install package.\n")
	cmd := fmt.Sprintf("tar -zxvf %s -C %s", remotes, s.InstallPath)
	log.Logger.Infof("CommandExec: %s\n", cmd)
	_, err := remote.CommandExec(client, cmd)
	if err != nil {
		log.Logger.Errorf("Decompress the install package failed. %s.\n", err.Error())
		return false
	}
	log.Logger.Infof("Decompress the install package success.\n")
	return true
}

func (s *Service) DeleteService(port string, user string, password string) bool {
	if !s.InitInstallMsg() {
		return false
	}
	// 1、对端服务器是否能够ssh登录
	ssh := remote.SSHConfig{
		IP:       s.IP,
		PORT:     port,
		USER:     user,
		PASSWORD: password,
	}
	client, err := remote.SSHConnect(&ssh)
	if err != nil {
		log.Logger.Errorln(err)
		return false
	}
	defer client.Close() // 用完关闭好习惯
	// 2、对端是否存在安装目录，如果存在则对其进行bak
	log.Logger.Infof("Begin delete the installed prometheus. Node: %s", s.IP)
	cmd := fmt.Sprintf("ls %s", s.InstallPath)
	output, err := remote.CommandExec(client, cmd)
	if err != nil {
		if strings.Contains(output, "No such file or directory") {
			log.Logger.Infof("The install path is not exist, no operation is required.\n")
			return true
		} else if strings.Contains(output, "Permission denied") {
			log.Logger.Errorf("The install path is exist. But the user %s has no permission to access the install path.\n", user)
			return false
		}
	} else {
		bakPath := fmt.Sprintf("%s_%v", s.InstallPath, time.Now().Unix())
		cmd := fmt.Sprintf("mv %s %s", s.InstallPath, bakPath)
		_, err := remote.CommandExec(client, cmd)
		if err != nil {
			return false
		}
	}
	return true
}

func (s *Service) CreateUser(client *ssh.Client, ssh remote.SSHConfig, user string, group string) bool {
	log.Logger.Infof("Create the service user %s start.\n", s.User)
	sMsg := ScriptMsg{
		Name:  "user_create.py",
		Path:  "scripts/python/user_create.py",
		Type:  "python",
		Param: []string{group, user},
	}
	if !sMsg.UploadRemoteScript(ssh) {
		return false
	}
	if !sMsg.RunRemoteScript(client) {
		return false
	}
	log.Logger.Infof("Create the service user %s successed.\n", s.User)
	return true
}

func (s *Service) InitEnv(client *ssh.Client, ssh remote.SSHConfig) bool {
	log.Logger.Infof("Init the service %s env start.\n", s.Name)
	sMsg := ScriptMsg{
		Name:  "init_env.py",
		Path:  "scripts/python/init_env.py",
		Type:  "python",
		Param: []string{s.Name, s.InstallPath},
	}
	if !sMsg.UploadRemoteScript(ssh) {
		return false
	}
	if !sMsg.RunRemoteScript(client) {
		return false
	}
	log.Logger.Infof("Init the service %s env successed.\n", s.Name)
	return true
}

// 初始化配置，这里允许传入一个接口
func InitConf(confs interface{}) (bool, string) {
	var filename string
	switch confs.(type) {
	case PrometheusConf:
		log.Logger.Infof("Init the service Prometheus configure start.\n")
		// 若配置类型为prometheus的配置，则进行prometheus配置初始化
		// 需要生成对应的config.ini文件
		c, ok := confs.(PrometheusConf)
		if !ok {
			log.Logger.Errorf("Configure type is error.\n")
		}
		c.DataDay = "30"
		c.ConfPath = "conf/prometheus.yml"
		uuid := uuid.NewV4().String()
		filename = fmt.Sprintf("config.ini_%s", uuid)
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Logger.Errorf("Open the %s failed.\n", filename)
			return false, ""
		}
		defer f.Close()
		t := reflect.TypeOf(c)
		v := reflect.ValueOf(c)
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			item := fmt.Sprintf("%s=%s\n", strings.ToUpper(sf.Tag.Get("json")), v.Field(i).Interface())
			_, err := f.WriteString(item)
			if err != nil {
				log.Logger.Printf(err.Error())
				return false, ""
			}
		}
		log.Logger.Infof("Init the service Prometheus configure end.\n")

	case NodeExporterConf:
		log.Logger.Infof("Init the service Node Exporter configure start.\n")
		c, ok := confs.(NodeExporterConf)
		if !ok {
			log.Logger.Errorf("Configure type is error.\n")
		}
		uuid := uuid.NewV4().String()
		filename = fmt.Sprintf("config.ini_%s", uuid)
		f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return false, ""
		}
		defer f.Close()
		t := reflect.TypeOf(c)
		v := reflect.ValueOf(c)
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			item := fmt.Sprintf("%s=%s\n", strings.ToUpper(sf.Tag.Get("json")), v.Field(i).Interface())
			_, err := f.WriteString(item)
			if err != nil {
				log.Logger.Printf(err.Error())
				return false, ""
			}
		}
		log.Logger.Infof("Init the service Node Exporter configure end.\n")
	}
	return true, filename
}

func (s *Service) DistributeConf(ssh remote.SSHConfig) bool {
	log.Logger.Infof("Begin to distribute the service conf file.\n")
	tmp := setting.Conf.AppBaseConfig.TmpDir
	local := fmt.Sprintf("%s/%s", tmp, s.Conf)
	remotes := fmt.Sprintf("%s/%s/conf/%s", s.InstallPath, s.Name, "config.ini")
	// 本地位置放在/tmp下，所以就是tmp+文件名
	err := remote.CopyToRemote(local, remotes, &ssh)
	if err != nil {
		log.Logger.Errorf("Upload the config.ini failed. May be you should check it. %s\n", err)
		return false
	}
	log.Logger.Infof("Upload the config.ini successed.\n")
	return true
}

func (s *Service) InstallService(port string, user string, password string) bool {
	if !s.InitInstallMsg() {
		return false
	}
	// 组件安装步骤：
	// 1、对端服务器是否能够ssh登录
	ssh := remote.SSHConfig{
		IP:       s.IP,
		PORT:     port,
		USER:     user,
		PASSWORD: password,
	}
	client, err := remote.SSHConnect(&ssh)
	if err != nil {
		log.Logger.Errorf("%s\n", err)
		return false
	}
	defer client.Close() // 用完关闭好习惯
	// 2、对端的安装目录是否已存在，不存在则创建
	if !CreateInstallPath(client, s) {
		return false
	}
	// 3、安装包上传
	localpath := fmt.Sprintf("./parcels/%s", s.PackageName)
	remotepath := fmt.Sprintf("/tmp/%s", s.PackageName)
	if !s.UploadInstallPackages(localpath, remotepath, ssh) {
		return false
	}
	// 4、解压安装包，进行安装
	if !s.DecompressInstallPackages(remotepath, client) {
		return false
	}
	// 5、目录权限处理
	s.CreateUser(client, ssh, "prometheus", "prometheus")
	// 6、分发启动配置文件
	s.DistributeConf(ssh)
	// 7、基础环境初始化依赖于配置文件的参数值，那么应该先初始化配置文件
	if !s.InitEnv(client, ssh) {
		return false
	}
	log.Logger.Infof("The %s service install successed.\n", s.Name)
	return true
}

func CreateClient(ip string, port string, user string, password string) (client *ssh.Client, err error) {
	ssh := remote.SSHConfig{
		IP:       ip,
		PORT:     port,
		USER:     user,
		PASSWORD: password,
	}
	client, err = remote.SSHConnect(&ssh)
	if err != nil {
		log.Logger.Errorln(err)
		return nil, err
	}
	return client, err
}
