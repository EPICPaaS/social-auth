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

type Renren struct {
	BaseProvider
}

func (p *Renren) GetType() social.SocialType {
	return social.SocialRenren
}

func (p *Renren) GetName() string {
	return "Renren"
}

func (p *Renren) GetPath() string {
	return "renren"
}

func (p *Renren) GetIndentify(tok *social.Token) (string, error) {
	fmt.Println(tok)
	uri := "https://graph.renren.com/oauth/token?grant_type=authorization_code&code=" + url.QueryEscape(tok.AccessToken)
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

var _ social.Provider = new(Renren)

func NewRenren(clientId, secret string) *Renren {
	p := new(Renren)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = ""
	p.AuthURL = "https://graph.renren.com/oauth/authorize"
	p.TokenURL = "https://graph.renren.com/oauth/token"
	p.RedirectURL = social.DefaultAppUrl + "login/renren/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
