package interfaces

import (
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type controller struct {
	mqtt            handler.MqttRepository
	logger          handler.LoggerRepository
	syslogMng       usecases.SyslogManager
	texMoneyHandler usecases.TexMoneyHandlerRepository
	noticeSend      NoticeSendRepository
}

// コールバック制御
func NewController(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	texMoneyHandler usecases.TexMoneyHandlerRepository) ControllerRepository {
	return &controller{
		mqtt:            mqtt,
		logger:          logger,
		syslogMng:       syslogMng,
		texMoneyHandler: texMoneyHandler,
		noticeSend:      NewNoticeSend(mqtt, logger)}
}

// 開始処理
func (c *controller) Start() {
	// ハンドラーへ各種状況変化コールバックを登録
	c.texMoneyHandler.RegisterCallbackNoticeIndata(c.callbackNoticeIndata)
	c.texMoneyHandler.RegisterCallbackNoticeOutdata(c.callbackNoticeOutdata)
	c.texMoneyHandler.RegisterCallbackNoticeCollectdata(c.callbackNoticeCollectdata)
	c.texMoneyHandler.RegisterCallbackNoticeAmountData(c.callbackNoticeAmountData)
	c.texMoneyHandler.RegisterCallbackNoticeStatusdata(c.callbackNoticeStatusdata)
	c.texMoneyHandler.RegisterCallbackNoticeReportStatusdata(c.callbackNoticeReportStatusdata)
	c.texMoneyHandler.RegisterCallbackNoticeExchangeStatusdata(c.callbackNoticeExchangeStatusdata)
	c.texMoneyHandler.Start()
}

// 停止処理
func (c *controller) Stop() {
	c.texMoneyHandler.RegisterCallbackNoticeIndata(nil)
	c.texMoneyHandler.RegisterCallbackNoticeOutdata(nil)
	c.texMoneyHandler.RegisterCallbackNoticeCollectdata(nil)
	c.texMoneyHandler.RegisterCallbackNoticeAmountData(nil)
	c.texMoneyHandler.RegisterCallbackNoticeStatusdata(nil)
	c.texMoneyHandler.RegisterCallbackNoticeReportStatusdata(nil)
	c.texMoneyHandler.RegisterCallbackNoticeExchangeStatusdata(nil)
	c.texMoneyHandler.Stop()
}

// 入出金管理:入金データ通知
func (c *controller) callbackNoticeIndata(texCon *domain.TexContext, noticeInfo *domain.StatusIndata) {
	topic := domain.TOPIC_TEXMONEY_BASE + "/notice_indata"

	err := c.noticeSend.Send(texCon, noticeInfo, topic)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_INDATA_FATAL, "", "入出金管理")
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_INDATA_SUCCESS, "", "入出金管理")
}

// 入出金管理:出金データ通知
func (c *controller) callbackNoticeOutdata(texCon *domain.TexContext, noticeInfo *domain.StatusOutdata) {
	topic := domain.TOPIC_TEXMONEY_BASE + "/notice_outdata"

	err := c.noticeSend.Send(texCon, noticeInfo, topic)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_OUTDATA_FATAL, "", "入出金管理")
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_OUTDATA_SUCCESS, "", "入出金管理")
}

// 入出金管理:回収データ通知
func (c *controller) callbackNoticeCollectdata(texCon *domain.TexContext, noticeInfo *domain.StatusCollectData) {
	topic := domain.TOPIC_TEXMONEY_BASE + "/notice_collectdata"

	err := c.noticeSend.Send(texCon, noticeInfo, topic)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_COLLECTDATA_FATAL, "", "入出金管理")
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_COLLECTDATA_SUCCESS, "", "入出金管理")
}

// 入出金管理:有高データ通知
func (c *controller) callbackNoticeAmountData(texCon *domain.TexContext, noticeInfo *domain.StatusAmount) {
	topic := domain.TOPIC_TEXMONEY_BASE + "/notice_amount"

	err := c.noticeSend.Send(texCon, noticeInfo, topic)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_AMOUNT_FATAL, "", "入出金管理")
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_AMOUNT_SUCCESS, "", "入出金管理")
}

// 入出金管理:現金入出金制御データ通知
func (c *controller) callbackNoticeStatusdata(texCon *domain.TexContext, noticeInfo *domain.StatusCash) {
	topic := domain.TOPIC_TEXMONEY_BASE + "/notice_status_cash"

	err := c.noticeSend.Send(texCon, noticeInfo, topic)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_STATUS_CASH_FATAL, "", "入出金管理")
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_STATUS_CASH_SUCCESS, "", "入出金管理")
}

// 入出金管理:入出金レポート印刷ステータス通知
func (c *controller) callbackNoticeReportStatusdata(texCon *domain.TexContext, noticeInfo *domain.StatusReport) {
	topic := domain.TOPIC_TEXMONEY_BASE + "/notice_report_status"

	err := c.noticeSend.Send(texCon, noticeInfo, topic)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_REPORT_STATUS_FATAL, "", "入出金管理")
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_REPORT_STATUS_SUCCESS, "", "入出金管理")
}

// 入出金管理:両替ステータス通知
func (c *controller) callbackNoticeExchangeStatusdata(texCon *domain.TexContext, noticeInfo *domain.StatusExchange) {
	topic := domain.TOPIC_TEXMONEY_BASE + "/notice_status_exchange"

	err := c.noticeSend.Send(texCon, noticeInfo, topic)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_STATUS_EXCHANGE_FATAL, "", "入出金管理")
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_STATUS_EXCHANGE_SUCCESS, "", "入出金管理")
}
