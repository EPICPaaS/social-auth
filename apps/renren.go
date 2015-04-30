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
	//"fmt"

	"encoding/json"
	"github.com/EPICPaaS/social-auth"
	"github.com/astaxie/beego/httplib"
	"net/url"
	"strconv"
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
	//fmt.Println(tok.GetExtra("id"))
	uri := "https://api.renren.com/v2/user/login/get?access_token=" + url.QueryEscape(tok.AccessToken)
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.String()
	if err != nil {
		return "", err
	}
	//fmt.Println(body)

	var rd map[string]interface{}
	err = json.Unmarshal([]byte(body), &rd)

	if err == nil {
		user := rd["response"].(map[string]interface{})

		uid := user["id"].(float64)
		//fmt.Println(uid)
		ruid := strconv.FormatFloat(uid, 'f', -1, 64)
		//fmt.Println(ruid)
		return ruid, nil
	}
	return "", err
}

//TODO 待完善
func (p *Renren) GetUserNname(tok *social.Token) (string, error) {
	return "", nil
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
