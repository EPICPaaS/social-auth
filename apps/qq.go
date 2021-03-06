// Copyright 2014 beego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//
// Maintain by https://github.com/slene

package apps

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego/httplib"

	"encoding/json"
	"github.com/EPICPaaS/social-auth"
)

type QQ struct {
	BaseProvider
}

func (p *QQ) GetType() social.SocialType {
	return social.SocialQQ
}

func (p *QQ) GetName() string {
	return "QQ"
}

func (p *QQ) GetPath() string {
	return "qq"
}

func (p *QQ) GetIndentify(tok *social.Token) (string, error) {
	uri := "https://graph.qq.com/oauth2.0/me?access_token=" + url.QueryEscape(tok.AccessToken)

	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.String()
	if err != nil {
		return "", err
	}
	bodyLen := len(body)
	bodySub := body[10 : bodyLen-4]
	dataMap := make(map[string]string)
	err = json.Unmarshal([]byte(bodySub), &dataMap)
	if err != nil {
		return "", err
	}

	if len(dataMap["code"]) != 0 {
		return "", fmt.Errorf("code: %s, msg: %s", dataMap["code"], dataMap["msg"])
	}
	return dataMap["openid"], nil
}

func (p *QQ) GetUserInfo(tok *social.Token) (string, error) {

	openid, err := p.GetIndentify(tok)
	if err != nil {
		return "", err
	}

	uri := "https://graph.qq.com/user/get_user_info?access_token=" + url.QueryEscape(tok.AccessToken) + "&oauth_consumer_key=" + p.ClientId + "&openid=" + openid

	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.Bytes()
	if err != nil {
		return "", err
	}
	var temp = map[string]interface{}{}
	if err = json.Unmarshal(body, &temp); err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	//目前暂时取昵称
	nickname, ok := temp["nickname"].(string)
	if !ok {
		return "", fmt.Errorf("code: %v, msg: %v", temp["ret"], temp["msg"])
	}

	return "QQ_" + nickname, nil
}

var _ social.Provider = new(QQ)

func NewQQ(clientId, secret string) *QQ {
	p := new(QQ)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "get_user_info"
	p.AuthURL = "https://graph.qq.com/oauth2.0/authorize"
	p.TokenURL = "https://graph.qq.com/oauth2.0/token"
	p.RedirectURL = social.DefaultAppUrl + "login/qq/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
