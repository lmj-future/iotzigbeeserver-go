package httprequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

// /**
//  * Created by Administrator on 2018/1/13.
//  */
// var request = require("request");

//HTTPOptions HTTPOptions
type HTTPOptions struct {
	url     string
	method  string
	headers map[string]string
	body    []byte
}

// HTTPData1 HTTPData1
type HTTPData1 struct {
	ErrCode int                `json:"errCode"`
	Data    *config.DevMgrInfo `json:"data"`
}

// Data Data
type Data struct {
	Result string `json:"result"`
}

// HTTPData2 HTTPData2
type HTTPData2 struct {
	ErrCode int  `json:"errCode"`
	Data    Data `json:"data"`
}

// HTTPRequest HTTPRequest
func HTTPRequest(url string, flag string) (interface{}, error) {
	var options = HTTPOptions{}
	options.url = url
	options.method = "GET"
	options.headers = map[string]string{
		"content-type": "application/json",
		"encoding":     "UTF-8",
	}
	body, _, err := request(options)
	if err != nil {
		globallogger.Log.Errorln("problem with request:", err)
		return nil, err
	}
	if flag == "terminalInfoReq" {
		var resdata = HTTPData1{}
		json.Unmarshal(body, &resdata)
		return &resdata, nil
	} else if flag == "licenseCheck" {
		var resdata = HTTPData2{}
		json.Unmarshal(body, &resdata)
		return &resdata, nil
	}
	return nil, nil
}

func request(options HTTPOptions) ([]byte, []byte, error) {
	method := options.method
	url := options.url
	body := options.body
	header := options.headers
	res, err := SendURL(method, url, bytes.NewBuffer(body), header)
	if err != nil {
		return nil, nil, err
	}
	var resBody []byte
	if res != nil {
		defer res.Close()
		resBody, err = ioutil.ReadAll(res)
		if err != nil {
			return nil, nil, err
		}
	}
	return resBody, body, err
}

// SendURL sends http request by url
func SendURL(method, url string, body io.Reader, header map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	var cli = http.DefaultClient
	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		var resBody []byte
		if res.Body != nil {
			defer res.Body.Close()
			resBody, _ = ioutil.ReadAll(res.Body)
		}
		return nil, fmt.Errorf("[%d] %s", res.StatusCode, strings.TrimRight(string(resBody), ""))
	}
	return res.Body, nil
}
