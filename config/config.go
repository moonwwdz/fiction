package config

import (
	"bytes"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ConfFile struct {
	Title string
	Conf  []Fconf
	Type  string //发送消息方式：weixin/telegram
	Wx    weixin
	Tg    telegram
}

type Fconf struct {
	Name      string
	Url       string
	TableName string
	ExecTime  string
}

type weixin struct {
}

type telegram struct {
	Token  string
	ChatId string
}

var confName = "conf.toml"

func init() {
	//配置文件不存在时，创建一个模板文件
	_, err := os.Stat(confName)
	if err == nil {
		return
	}

	var confs []Fconf
	confs = append(confs, Fconf{
		Name:      "圣虚",
		Url:       "/",
		TableName: "shenxu",
		ExecTime:  "0 8 * * *",
	})
	var initConf = ConfFile{
		Title: "小说订阅",
		Conf:  confs,
	}

	var firstBuffer bytes.Buffer
	e := toml.NewEncoder(&firstBuffer)
	err = e.Encode(initConf)
	if err != nil {
		return
	}

	if _, err := os.Create(confName); err != nil {
		log.Fatalln("Failed to create new config")
		return
	}

	f, err := os.OpenFile(confName, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln("Failed to write config file")
		return
	}

	f.Write(firstBuffer.Bytes())
	f.Close()
}

// GetConf 获取配置文件
func GetConf() (*ConfFile, error) {
	var confInfo ConfFile

	_, err := toml.DecodeFile(confName, &confInfo)
	if err != nil {
		log.Fatalf("%v", err)
		log.Println("Failed read config")
		return nil, err
	}
	return &confInfo, nil
}
