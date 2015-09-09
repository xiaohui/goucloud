package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/xiaohui/goucloud/ucloud"
	"github.com/xiaohui/goucloud/ucloud/auth"
	"github.com/xiaohui/goucloud/ucloud/request"
)

type Service struct {
	Config      *ucloud.Config
	ServiceName string
	APIVersion  string

	BaseUrl    string
	HttpClient *http.Client
}

func (s *Service) DoRequest(url string, params interface{}, response interface{}) error {
	requestURL, err := s.RequestURL(params)
	if err != nil {
		return fmt.Errorf("build request url failed, error: %s", err)
	}

	httpReq, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("new request url failed, error: %s", err)
	}

	httpResp, err := s.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request url failed, error: %s", err)
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return fmt.Errorf("do request url failed, error: %s", err)
	}

	statusCode := httpResp.StatusCode
	if statusCode >= 400 && statusCode <= 599 {

		//TODO: parse the error messages
		return fmt.Errorf("request error, status code:%s", statusCode)
	}

	err = json.Unmarshal(body, response)
	//log.Printf("%++v", response)

	if err != nil {
		return fmt.Errorf("unmarshal url failed, error: %s", err)
	}

	return nil
}

// RequestURL is fully url of api request
func (s *Service) RequestURL(params interface{}) (string, error) {
	if len(s.BaseUrl) == 0 {
		return "", errors.New("baseUrl is not set")
	}

	commonRequest := request.CommonRequest{
		PublicKey: s.Config.Credentials.PublicKey,
		ProjectId: s.Config.ProjectID,
	}

	values := url.Values{}
	convertParamsToValues(commonRequest, &values)
	convertParamsToValues(params, &values)

	url, err := urlWithSignature(&values, s.BaseUrl, s.Config.Credentials.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("convert params error: %s", err)
	}

	return url, nil
}

func convertParamsToValues(params interface{}, values *url.Values) {

	elem := reflect.ValueOf(params)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	elemType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		fieldName := elemType.Field(i).Name

		field := elem.Field(i)
		kind := field.Kind()
		if (kind == reflect.Ptr ||
			kind == reflect.Array ||
			kind == reflect.Slice ||
			kind == reflect.Map ||
			kind == reflect.Chan) && field.IsNil() {
			continue

		}

		if kind == reflect.Ptr {
			field = field.Elem()
			kind = field.Kind()
		}

		var v string
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = strconv.FormatInt(field.Int(), 10)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = strconv.FormatUint(field.Uint(), 10)

		case reflect.Float32:
			v = strconv.FormatFloat(field.Float(), 'f', 4, 32)

		case reflect.Float64:
			v = strconv.FormatFloat(field.Float(), 'f', 4, 64)

		case reflect.Bool:
			v = strconv.FormatBool(field.Bool())

		case reflect.String:
			v = field.String()
		}

		if v != "" {
			name := elemType.Field(i).Tag.Get("ArgName")
			if name == "" {
				name = fieldName
			}

			values.Set(name, v)
		}
	}
}

func urlWithSignature(values *url.Values, baseUrl, privateKey string) (string, error) {

	urlEncoded, err := url.QueryUnescape(values.Encode())
	if err != nil {
		return "", fmt.Errorf("unescape failed, error: %s", err)
	}

	// replace '&' and '=' in url
	urlEncoded = strings.Replace(urlEncoded, "=", "", -1)
	urlEncoded = strings.Replace(urlEncoded, "&", "", -1)

	signature, err := auth.GenerateSignature(urlEncoded, privateKey)
	if err != nil {
		return "", fmt.Errorf("generate signature error:%s", err)
	}

	return baseUrl + "?" + values.Encode() + "&Signature=" + url.QueryEscape(signature), nil
}
