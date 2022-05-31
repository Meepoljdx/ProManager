package command

import (
	"ProManager/module/local"
	"ProManager/module/remote"
	log "ProManager/module/utils/logs"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type ScriptMsg struct {
	Name  string   // 脚本名称
	Path  string   // 脚本本地路径
	Type  string   // 脚本类型
	Param []string // 脚本参数
}

func InitRemoteScript() { // 初始化远程脚本信息

}

func (s *ScriptMsg) UploadRemoteScript(ssh remote.SSHConfig) bool {
	log.Logger.Infof("Begin to upload the %s script. %s.\n", s.Type, s.Path)
	remotepath := fmt.Sprintf("/tmp/%s", s.Name)
	err := remote.CopyToRemote(s.Path, remotepath, &ssh)
	if err != nil {
		log.Logger.Errorf("Uploead script %s to %s failed.\n", s.Name, ssh.IP)
		return false
	}
	log.Logger.Infof("Uploead script %s to %s successed.\n", s.Name, ssh.IP)
	return true
}

func (s *ScriptMsg) RunRemoteScript(client *ssh.Client) (success bool) { // 运行远程脚本
	remotepath := fmt.Sprintf("/tmp/%s", s.Name)
	cmd := ""
	switch s.Type {
	case "python":
		cmd = fmt.Sprintf("python %s", remotepath)
	case "shell":
		cmd = fmt.Sprintf("sh %s", remotepath)
	default:
		log.Logger.Warnf("Unadapted script type\n")
		return false
	}
	for _, v := range s.Param {
		cmd = fmt.Sprintf("%s %s", cmd, v) // 拼接脚本参数
	}
	output, err := remote.CommandExec(client, cmd)
	log.Logger.Infof("-----Run the remote %s Script log begin-----\n%s", s.Type, output)
	if err != nil {
		log.Logger.Warnf("Run remote script %s return is not 0. command is %s, %s\n", remotepath, cmd, err)
		success = false
	}
	log.Logger.Infof("-----Run the remote %s Script log end-----\n", s.Type)
	success = true
	return success
}

func (s *ScriptMsg) RunLocalScript() (success bool) { // 本地运行脚本
	cmd := ""
	switch s.Type {
	case "python":
		cmd = fmt.Sprintf("python %s", s.Path)
	case "shell":
		cmd = fmt.Sprintf("sh %s", s.Path)
	default:
		log.Logger.Warnf("Unadapted script type\n")
		return false
	}
	for _, v := range s.Param {
		cmd = fmt.Sprintf("%s %s", cmd, v) // 拼接脚本参数
	}
	output, err := local.CommandExec(cmd)
	log.Logger.Infof("-----Run the local %s Script log begin-----\n%s", s.Type, output)
	if err != nil {
		log.Logger.Warnf("Run local script %s return is not 0. command is %s, %s\n", s.Path, cmd, err)
		success = false
	}
	log.Logger.Infof("-----Run the local %s Script log end-----\n%s", s.Type, output)
	success = true
	return success
}

func RemoveTmpFile(filename string) (err error) {
	log.Logger.Infof("Begin to remove temp config file %s.\n", filename)
	if err = os.Remove(filename); err != nil {
		log.Logger.Warnf("Remove temp config file %s failed.\n", filename)
		return err

	}
	log.Logger.Infof("Remove temp config file %s success.\n", filename)
	return nil
}
