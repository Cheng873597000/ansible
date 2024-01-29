package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apenella/go-ansible/pkg/playbook"
)

type Host struct {
	AnsibleShellType  string `json:"ansible_shell_type"`
	AnsibleConnection string `json:"ansible_connection"`
	AnsibleHost       string `json:"ansible_host"`
	AnsiblePort       string `json:"ansible_port"`
	AnsibleUser       string `json:"ansible_user"`
	AnsiblePassword   string `json:"ansible_password"`
}

func Generate(path_dir string) ([]byte, error) {

	appServerName := "23"
	ansibleHost := "192.168.28.23"
	AnsiblePort := "22"
	AnsibleUser := "administrator"
	AnsiblePassword := "Win.2012"
	host := &Host{
		AnsibleShellType:  "cmd",
		AnsibleConnection: "ssh",
		AnsibleHost:       ansibleHost,
		AnsiblePort:       AnsiblePort,
		AnsibleUser:       AnsibleUser,
		AnsiblePassword:   AnsiblePassword,
	}

	// 将 data 转换为 JSON 字符串并打印
	Marshalresult, err := json.MarshalIndent(host, "", "")
	if err != nil {
		fmt.Println("JSON encoding error:", err)
		return nil, err
	}

	Hostmsg := fmt.Sprintf(`{
		"all": {
			"hosts": {
				%q: %s
				
			}
		}
	}`, appServerName, string(Marshalresult))

	fmt.Println("Hostmsg:", Hostmsg)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(Hostmsg), &data); err != nil {
		return nil, err
	}

	convert, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	return convert, nil
}

// 将被用于生成inventory时将hosts信息写入文件
func WriteToFile(path string) (err error) {
	pathDir := filepath.Dir(path)
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		err := os.MkdirAll(pathDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fmt.Println("创建文件夹成功！")

	generateData, err := Generate(pathDir)
	if err != nil {
		return err
	}

	fmt.Println("==============Generate生成数据成功！\n", string(generateData))

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	fmt.Println("生成文件成功！")

	defer file.Close()

	// 将 JSON 字符串写入文件
	err = os.WriteFile("gotest.yml", generateData, 0700)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}
	fmt.Println("写入文件成功！")

	return nil
}

func main() {
	WriteToFile("/cn/test/ansible/gotest.yml")

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: "gotest.yml",
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks: []string{"playbook.yml"},
		Options:   ansiblePlaybookOptions,
	}

	err := playbook.Run(context.TODO())
	if err != nil {
		panic(err)
	}

}
