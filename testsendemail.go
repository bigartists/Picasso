package main

import (
	"davinci/config"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func main() {
	e := email.NewEmail()
	//设置发送方的邮箱
	//e.From = "rh <154529156@qq.com>"
	e.From = config.SysYamlconfig.SMTP.From
	fmt.Println("e.From:=", e.From)
	// 设置接收方的邮箱
	e.To = []string{"2856197796@qq.com"}

	//设置主题
	e.Subject = "这是主题"
	//设置文件发送的内容
	e.HTML = []byte(`
    <h1><a href="javascript:;">from davinci by yaml hahaha</a></h1>    
    `)
	//设置服务器相关的配置
	domain := fmt.Sprintf("%s:%s", config.SysYamlconfig.SMTP.Host, config.SysYamlconfig.SMTP.Port)
	fmt.Println("domain:=", domain)
	fmt.Println("config.SysYamlconfig.SMTP.Username:=", config.SysYamlconfig.SMTP.Username,
		"config.SysYamlconfig.SMTP.Password:=", config.SysYamlconfig.SMTP.Password,
		"config.SysYamlconfig.SMTP.Host:=", config.SysYamlconfig.SMTP.Host)
	err := e.Send(domain, smtp.PlainAuth("", config.SysYamlconfig.SMTP.Username, config.SysYamlconfig.SMTP.Password, config.SysYamlconfig.SMTP.Host))
	if err != nil {
		log.Fatal(err)
	}
}
