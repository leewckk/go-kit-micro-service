package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/resty.v1"
)

func RESTPost(url string, body interface{}, debug bool) (string, error) {

	js, _ := json.Marshal(body)
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(js)).
		Post(url)

	if debug == true || err != nil {
		fmt.Println("payload: ", string(js))
		RESTDebugResp(url, resp, err)
	}
	return resp.String(), err
}

func RESTPostWithHeader(url string, body interface{}, header map[string]string, debug bool) (string, error) {
	js, _ := json.Marshal(body)
	header["Content-Type"] = "application/json"

	resp, err := resty.R().
		SetHeaders(header).
		SetBody(string(js)).
		Post(url)

	if debug == true || err != nil {
		RESTDebugResp(url, resp, err)
	}
	return resp.String(), err
}

func RESTPost2(url string, body interface{}, debug bool) (string, error) {

	js, _ := json.Marshal(body)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
	if nil != err {
		fmt.Println(err)
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(respBody), err
}

func RESTGet(url string, debug bool, params ...string) (string, error) {
	var resp *resty.Response
	var err error

	if 0 == len(params) {
		resp, err = resty.R().Get(url)
	} else {
		resp, err = resty.R().SetQueryString(params[0]).Get(url)
	}

	if debug == true || err != nil {
		RESTDebugResp(url, resp, err)
	}
	return resp.String(), err
}

func RESTPut(url string, body string, debug bool) (string, error) {

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(url)

	if debug == true || err != nil {
		RESTDebugResp(url, resp, err)
	}
	return resp.String(), err
}

func RESTPut2(url string, body interface{}, debug bool) (string, error) {
	js, err := json.Marshal(body)
	if nil != err {
		return err.Error(), err
	}
	return RESTPut(url, string(js), debug)
}

func RESTDelete(url string, body string, debug bool) (string, error) {

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Delete(url)

	if debug == true || err != nil {
		RESTDebugResp(url, resp, err)
	}
	return resp.String(), err
}

func RESTDebugResp(requesturl string, resp *resty.Response, err error) {
	if nil != err {
		log.Println(err)
		return
	}
	fmt.Printf("\nRequest: %s", requesturl)
	fmt.Printf("\nMethod: %s", resp.Request.Method)
	fmt.Printf("\nHead: %v", resp.Request.Header)
	fmt.Printf("\nBody: %v", resp.Request.Body)

	fmt.Printf("\nERROR: %v", err)
	fmt.Printf("\nResponse status code: %v", resp.StatusCode())
	fmt.Printf("\nResponse status: %v", resp.Status())
	fmt.Printf("\nResponse time: %v", resp.Time())
	fmt.Printf("\nResponse received at: %v", resp.ReceivedAt())
	fmt.Printf("\nResponse body: %v\n", resp)
}
