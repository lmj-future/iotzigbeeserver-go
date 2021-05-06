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

// Data3 Data3
type Data3 struct {
	ErrCode  int    `json:"errCode"`
	Reason   string `json:"reason"`
	TmnDevSN string `json:"tmnDevSN"`
}

// HTTPData3 HTTPData3
type HTTPData3 struct {
	ErrCode int     `json:"errCode"`
	Reason  string  `json:"reason"`
	Data    []Data3 `json:"data"`
}

// Data4 Data4
type Data4 struct {
	TenantID          string   `json:"tenantID"`     //用户名
	TerminalType      string   `json:"terminalType"` //型号名
	Alias             string   `json:"alias"`        //型号别名
	FirmTopic         string   `json:"firmTopic"`    //厂商topic
	FirmName          string   `json:"firmName"`     //厂商名
	AccessType        string   `json:"accessType"`   //接入方式
	ProductKey        string   `json:"productKey"`   //型号productKey
	SubDevice         string   `json:"subDevice"`    //设备类型：device-终端，gateway-网关
	DataFormat        string   `json:"dataFormat"`   //数据解析格式
	NewAccessType     []string `json:"newAccessType"`
	EdgeSpecification string   `json:"edgeSpecification"` //边缘规格
	SysType           string   `json:"sysType"`           //边缘配置-系统类型
	SysArchitecture   string   `json:"sysArchitecture"`   //边缘配置-系统架构
}

// HTTPData4 HTTPData4
type HTTPData4 struct {
	ErrCode int    `json:"errCode"`
	Reason  string `json:"reason"`
	Data    Data4  `json:"data"`
}

// HTTPRequest HTTPRequest
func HTTPRequest(url string, method string, flag string, body []byte) (interface{}, error) {
	response, _, err := request(HTTPOptions{
		url:    url,
		method: method,
		headers: map[string]string{
			"content-type": "application/json",
			"encoding":     "UTF-8",
		},
		body: body,
	})
	if err != nil {
		globallogger.Log.Errorln("problem with request:", err)
		return nil, err
	}
	if flag == "terminalInfoReq" {
		var resdata = HTTPData1{}
		json.Unmarshal(response, &resdata)
		return &resdata, nil
	} else if flag == "licenseCheck" {
		var resdata = HTTPData2{}
		json.Unmarshal(response, &resdata)
		return &resdata, nil
	} else if flag == "addTerminal" {
		var resdata = HTTPData3{}
		json.Unmarshal(response, &resdata)
		return &resdata, nil
	} else if flag == "terminalTypeByAlias" {
		var resdata = HTTPData4{}
		json.Unmarshal(response, &resdata)
		return &resdata, nil
	}
	return nil, nil
}

func request(options HTTPOptions) ([]byte, []byte, error) {
	res, err := SendURL(options.method, options.url, bytes.NewBuffer(options.body), options.headers)
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
	return resBody, options.body, err
}

// SendURL sends http request by url
func SendURL(method, url string, body io.Reader, header map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{}
	if len(header) > 0 {
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
