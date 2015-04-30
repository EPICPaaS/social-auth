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
	"encoding/json"
	"fmt"
	"github.com/EPICPaaS/social-auth"
	"github.com/astaxie/beego/httplib"
)

type Weibo struct {
	BaseProvider
}

func (p *Weibo) GetType() social.SocialType {
	return social.SocialWeibo
}

func (p *Weibo) GetName() string {
	return "Weibo"
}

func (p *Weibo) GetPath() string {
	return "weibo"
}

func (p *Weibo) GetIndentify(tok *social.Token) (string, error) {
	return tok.GetExtra("uid"), nil
}

//TODO 待完善
func (p *Weibo) GetUserInfo(tok *social.Token) (string, error) {

	uri := "https://api.weibo.com/2/users/show.json?source=" + p.ClientId + "&access_token=" + tok.AccessToken + "&uid=" + tok.GetExtra("uid")
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.Bytes()
	if err != nil {
		return "", err
	}

	var ret = map[string]interface{}{}
	if err := json.Unmarshal(body, &ret); err != nil {
		return "", err
	}

	userName, ok := ret["screen_name"].(string)
	if !ok {
		if userName, ok = ret["name"].(string); !ok {
			return "", fmt.Errorf("error_code:%v , [ERROR]: %v", ret["error_code"], ret["error"])
		}
	}
	return "新浪微博_" + userName, nil
}

var _ social.Provider = new(Weibo)

func NewWeibo(clientId, secret string) *Weibo {
	p := new(Weibo)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "email"
	p.AuthURL = "https://api.weibo.com/oauth2/authorize"
	p.TokenURL = "https://api.weibo.com/oauth2/access_token"
	p.RedirectURL = social.DefaultAppUrl + "login/weibo/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
