package main

import (
	"os"
	"stub/domain"
	"stub/infrastructure"
	"stub/pkg/session"

	"github.com/kardianos/service"
)

type program struct{}

// サービス開始
func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

// サービス停止
func (p *program) Stop(s service.Service) error {
	Exit()
	return nil
}

// サービス起動
func (p *program) run() error {
	serviceMain()
	return nil
}

// メイン
func main() {
	svcConfig := &service.Config{
		Name:        domain.SrvName,
		DisplayName: domain.DspName,
		Description: domain.Description,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		println("Cannot create the service:" + err.Error())
	}
	//起動時パラメータにてインストール･アンインストールを指定
	//指定のない場合は、そのまま実行。
	if len(os.Args) > 1 {
		// パラメータは、install or uninstall
		err = service.Control(s, os.Args[1])
		if err != nil {
			return
		}
		return
	}

	if session.GetProcessSessionID() == 0 {
		// SessionID = 0 は SYSTEM(サービス)
		_ = s.Run()
	} else {
		// SessionID != 0 は ユーザーとして起動
		serviceMain()
	}
}

// サービスメイン処理
func serviceMain() {
	//開始処理
	infrastructure.Router()
	//os.Exit(0)
}

// 終了検出
func Exit() {
	infrastructure.RouterStop()
}
