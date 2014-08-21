// Copyright 2014 EPICPaaS authors
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
// Maintain by https://github.com/EPICPaaS

package apps

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego/httplib"

	"github.com/EPICPaaS/social-auth"
)

type Weixin struct {
	BaseProvider
}

func (p *Weixin) GetType() social.SocialType {
	return social.SocialWeixin
}

func (p *Weixin) GetName() string {
	return "Weixin"
}

func (p *Weixin) GetPath() string {
	return "weixin"
}

func (p *Weixin) GetIndentify(tok *social.Token) (string, error) {
	uri := "https://api.weixin.qq.com/sns/oauth2/access_token?grant_type=authorization_code&code=" + url.QueryEscape(tok.AccessToken)
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.String()
	if err != nil {
		return "", err
	}

	vals, err := url.ParseQuery(body)
	if err != nil {
		return "", err
	}

	if vals.Get("code") != "" {
		return "", fmt.Errorf("code: %s, msg: %s", vals.Get("code"), vals.Get("msg"))
	}

	return vals.Get("openid"), nil
}

var _ social.Provider = new(Weixin)

func NewWeixin(clientId, secret string) *Weixin {
	p := new(Weixin)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "snsapi_login"
	p.AuthURL = "https://open.weixin.qq.com/connect/qrconnect"
	p.TokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	p.RedirectURL = social.DefaultAppUrl + "login/weixin/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
