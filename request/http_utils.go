package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whimthen/kits/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var usePool = false

type h struct {
	set bool
	k   string
	v   string
}

type HttpRequest struct {
	Url      string
	request  *http.Request
	response *http.Response
	doErr    error
	header   []*h
	cookies  []*http.Cookie

	queryString string
	body        string
}

type Resp struct {
	response *http.Response
}

type ReqHeader map[string]string

func Get(url string, params url.Values) *HttpRequest {
	request := &HttpRequest{Url: url}
	return request.DoGet(params)
}

func GetWithHeader(url string, params url.Values, headers ReqHeader) *HttpRequest {
	request := &HttpRequest{Url: url}
	if headers != nil {
		for k, v := range headers {
			request.AddHeader(k, v)
		}
	}
	return request.DoGet(params)
}

func (r *HttpRequest) ContentType(cType string) *HttpRequest {
	r.AddHeader("Content-Type", cType)
	return r
}

func (r *HttpRequest) JsonContentType() *HttpRequest {
	r.ContentType("application/json;charset=utf8")
	return r
}

func (r *HttpRequest) AddHeader(k, v string) *HttpRequest {
	if k == "" || v == "" {
		return r
	}
	r.header = append(r.header, &h{
		set: false,
		k:   k,
		v:   v,
	})
	return r
}

func (r *HttpRequest) SetHeader(k, v string) *HttpRequest {
	if r.doErr != nil {
		return r
	}
	if k == "" || v == "" {
		return r
	}
	r.header = append(r.header, &h{
		set: true,
		k:   k,
		v:   v,
	})
	return r
}

func (r *HttpRequest) AddCookie(cookie *http.Cookie) *HttpRequest {
	if r.doErr != nil {
		return r
	}
	if r.cookies == nil {
		r.cookies = []*http.Cookie{}
	}
	r.cookies = append(r.cookies, cookie)
	return r
}

func (r *HttpRequest) Set(username, password string) *HttpRequest {
	if r.doErr != nil {
		return r
	}
	r.request.SetBasicAuth(username, password)
	return r
}

func (r *HttpRequest) DoPost(body interface{}) *HttpRequest {
	if r.doErr != nil {
		return r
	}

	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			r.doErr = err
			return r
		}

		r.body = fmt.Sprintf("%s", bytes)
	}

	r.r(http.MethodPost)

	return r
}

func (r *HttpRequest) DoFromPost(params url.Values) *HttpRequest {
	if r.doErr != nil {
		return r
	}
	if params != nil {
		r.queryString = params.Encode()
	}
	r.ContentType("application/x-www-form-urlencoded")
	r.r(http.MethodPost)
	return r
}

func (r *HttpRequest) DoGet(params url.Values) *HttpRequest {
	if r.doErr != nil {
		return r
	}
	if params != nil {
		r.queryString = params.Encode()
	}
	r.r(http.MethodGet)
	return r
}

func (r *HttpRequest) r(method string) {
	err := r.newRequest(method)
	if err != nil {
		r.doErr = err
		return
	}
	if method == http.MethodGet || r.queryString != "" {
		logger.Info("Method: %s, URL: %s?%s", method, r.request.URL.String(), r.queryString)
	} else {
		logger.Info("Method: %s, URL: %s, Params: %s", method, r.request.URL.String(), r.body)
	}
	r.doRequest()
}

func (r *HttpRequest) Scan(model interface{}) (err error) {
	if r.doErr != nil {
		return r.doErr
	}
	if r.response == nil {
		return r.doErr
	}
	defer r.response.Body.Close()
	bytes, err := ioutil.ReadAll(r.response.Body)
	if err != nil {
		return r.doErr
	}
	if r.response.StatusCode == http.StatusOK {
		logger.Info("URL: %s, Response: %s\n", r.response.Request.URL.String(), string(bytes))
		return json.Unmarshal(bytes, &model)
	}
	return errors.New(string(bytes))
}

func (r *HttpRequest) ScanToMap() (result map[string]interface{}, err error) {
	if r.doErr != nil {
		return nil, r.doErr
	}
	err = r.Scan(&result)
	return
}

func (r *HttpRequest) Response() *http.Response {
	return r.response
}

func (r *HttpRequest) ResponseString() (string, error) {
	bytes, err := r.ResponseBytes()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bytes), nil
}

func (r *HttpRequest) ResponseBytes() ([]byte, error) {
	if r.doErr != nil {
		return nil, r.doErr
	}
	defer r.response.Body.Close()

	bytes, err := ioutil.ReadAll(r.response.Body)
	if err != nil {
		r.doErr = err
		return nil, err
	}
	return bytes, nil
}

func (r *HttpRequest) newRequest(method string) error {
	if r.doErr != nil {
		return r.doErr
	}
	var p io.Reader
	u := r.Url
	if method == http.MethodGet || r.queryString != "" {
		u += "?" + r.queryString
	} else {
		p = strings.NewReader(r.body)
	}
	request, err := http.NewRequest(method, u, p)
	if err != nil {
		logger.Error("Can't create http request, URL: [%s %s], err: %+v", method, r.Url, err)
		return err
	}
	r.request = request
	return nil
}

func (r *HttpRequest) doRequest() {
	if r.doErr != nil {
		return
	}
	client := http.DefaultClient
	if r.cookies != nil && len(r.cookies) > 0 {
		for _, cookie := range r.cookies {
			r.request.AddCookie(cookie)
		}
	}

	for _, v := range r.header {
		if v.set {
			r.request.Header.Set(v.k, v.v)
		} else {
			r.request.Header.Add(v.k, v.v)
		}
	}
	res, err := client.Do(r.request)

	if err != nil {
		r.doErr = err
	}
	r.response = res
}
