package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Get(url string, params url.Values) (map[string]interface{}, error) {
	var (
		err     error
		resp    *http.Response
		body    []byte
		jsonMap map[string]interface{}
	)

	if resp, err = http.Get(url + "?" + params.Encode()); err != nil {
		return jsonMap, err
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return jsonMap, err
	}
	defer resp.Body.Close()

	jsonMap = make(map[string]interface{})
	err = json.Unmarshal(body, &jsonMap)

	return jsonMap, err
}

func PostForm(url string, params url.Values) (map[string]interface{}, error) {
	var (
		err     error
		resp    *http.Response
		body    []byte
		jsonMap map[string]interface{}
	)

	if resp, err = http.PostForm(url, params); err != nil {
		return jsonMap, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return jsonMap, err
	}
	defer resp.Body.Close()

	jsonMap = make(map[string]interface{})
	err = json.Unmarshal(body, &jsonMap)

	return jsonMap, err
}

func PostJson(url string, params map[string]interface{}) (jsonMap map[string]interface{}, err error) {
	var (
		req     *http.Request
		jsonStr []byte
		client  *http.Client
		resp    *http.Response
		body    []byte
	)

	if jsonStr, err = json.Marshal(params); err != nil {
		return
	}
	if req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr)); err != nil {

	}
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	jsonMap = make(map[string]interface{})
	err = json.Unmarshal(body, &jsonMap)

	return
}
