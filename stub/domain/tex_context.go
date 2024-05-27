package domain

import (
	"fmt"
	"strings"
	"sync"
)

// 使い方
/*
詰めたい固有情報を、「RegisterTexContext」にセットして
NewTexContextを呼び出し利用する。
作成した情報は不変な物として扱う為、TexContextの中身の情報は
Getを使って取得する
今後、伝搬したい内容が増えた場合に、変更しやすいように
登録用と伝搬する情報を分離してみた。

	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceiveRequestInfo: reqInfo.RequestInfo,
	})
c.logger.Trace("%v", texCon.GetUniqueKey())

(texCon *domain.TexContext,)
texCon.GetUniqueKey()
*/

type (
	TexContext struct {
		uniqueKey string
		//		receiveRequestInfo RequestInfo
		receivingTopicName string
	}

	RegisterTexContext struct {
		//		ReceiveRequestInfo RequestInfo
		ReceivingTopicName string
	}
)

var fitAContextMutex = &sync.Mutex{}
var fitAContextUniqueKeyCount int

func NewTexContext(req RegisterTexContext) *TexContext {
	fitAContextMutex.Lock()
	defer fitAContextMutex.Unlock()

	fitAContextUniqueKeyCount++
	u := fitAContextUniqueKeyCount
	if u == 1000000 { // 最大6桁を超えた場合には1から再スタート
		u = 1
	}
	key := fmt.Sprintf("texMoney-%06d", u)
	if req.ReceivingTopicName != "" {
		key = fmt.Sprintf("%s-%06d", shorteningTopicName(req.ReceivingTopicName), u)
	}

	return &TexContext{
		uniqueKey:          key,
		receivingTopicName: req.ReceivingTopicName,
	}
}

func (tc *TexContext) GetUniqueKey() string {
	result := tc.uniqueKey
	return result
}

func shorteningTopicName(s string) string {
	topic := strings.Split(s, "_")
	if len(topic) == 0 {
		return s
	}

	var res string

	switch topic[0] {
	case "result":
		res += "res"
		for i, v := range topic {
			if i == 0 {
				continue
			}
			res += "_" + v
		}
		return res

	case "request":
		res += "req"
		for i, v := range topic {
			if i == 0 {
				continue
			}
			res += "_" + v
		}
		return res

	case "notice":
		res += "not"
		for i, v := range topic {
			if i == 0 {
				continue
			}
			res += "_" + v
		}
		return res
	}
	return s
}
