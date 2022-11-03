package rpc

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"
)

type Option struct {
	Sid string
	Tid string
}

func OptionWithSid(sid string) Option {
	return Option{Sid: sid}
}

func OptionWithTid(tid string) Option {
	return Option{Tid: tid}
}

func Rpc(request *http.Request, method string, addr string, uri string, query interface{}, body interface{}, opts ...Option) (string, error) {
	if addr == "" {
		return "", errors.New("没有找到指定的服务")
	}

	if addr != "" && strings.HasPrefix(addr, "/") {
		// 服务地址后统一去掉 "/"
		addr = addr[1:]
	}

	if uri != "" && !strings.HasPrefix(uri, "/") {
		uri += "/"
	}

	url := fmt.Sprintf("%s%s", addr, uri)
	userid := ""
	machineid := ""
	sid := ""
	tid := ""

	for _, opt := range opts {
		if opt.Sid != "" {
			sid = opt.Sid
		}

		if opt.Tid != "" {
			tid = opt.Tid
		}
	}

	if request != nil {
		userid = request.Header.Get("userid")
		machineid = request.Header.Get("machineid")
		if sid == "" {
			sid = request.Header.Get("sid")
		}
		if tid == "" {
			tid = request.Header.Get("tid")
		}
	}

	_, b, errs := gorequest.New().CustomMethod(method, url).
		Set("sid", sid).
		Set("tid", tid).
		Set("userid", userid).
		Set("machineid", machineid).
		Query(query).
		Send(body).
		End()
	if errs != nil {
		errtxt := fmt.Sprintf("调用服务出错，%s", errs[0].Error())
		return "", errors.New(errtxt)
	}

	return b, nil
}
