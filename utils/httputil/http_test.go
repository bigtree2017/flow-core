package httputil

import (
	"context"
	"encoding/json"
	"github.com/bigtree8/flow-core/http/ctxkit"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

type responseData struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

var client Client
var c *gin.Context

func init() {
	client = NewClient(time.Second * 5)
	c = &gin.Context{}
}

/*
*
<?php

	$data = [
	  "code" => 200,
	  "message" => "ok",
	  "data" => [
	      "type" => $_SERVER['CONTENT_TYPE'],
	      "post" => $_POST,
	      "get" => $_GET,
	      "input" => file_get_contents("php://input")
	  ]
	];
	echo json_encode($data);
*/
func TestDoGet(t *testing.T) {
	url := "http://localhost:8080/hello"
	req, _ := NewGetRequest(url, nil)
	ctxkit.GenerateTraceId(c)
	response, err := client.Do(c, req)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := DealResponse(response)

	resp := new(responseData)
	json.Unmarshal(result, resp)
	if resp.Code != 200 {
		t.Error("get result is not ok")
		return
	}
}

func TestPost(t *testing.T) {
	url := "http://localhost:8080/testPost"
	// 参数为空
	req, err := NewFormPostRequest(url, nil)
	response, err := client.Do(context.TODO(), req)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := DealResponse(response)
	resp := new(responseData)
	json.Unmarshal(result, resp)

	if resp.Code != 200 {
		t.Error("post result is not ok")
		return
	} else if resp.Data["type"] != ContentTypeForm {
		t.Error("post content-type is not equal " + ContentTypeForm)
		return
	}

	//参数为空map
	req, err = NewFormPostRequest(url, make(map[string]interface{}))
	response, err = client.Do(context.TODO(), req)
	if err != nil {
		t.Error(err)
		return
	}
	result, err = DealResponse(response)
	resp = new(responseData)
	json.Unmarshal(result, resp)
	if resp.Code != 200 {
		t.Error("post result is not ok")
		return
	}

	//参数非空map
	params := map[string]interface{}{
		"name": "hts",
	}
	req, err = NewFormPostRequest(url, params)
	response, err = client.Do(context.TODO(), req)
	if err != nil {
		t.Error(err)
		return
	}
	result, err = DealResponse(response)
	resp = new(responseData)
	json.Unmarshal(result, resp)
	if resp.Code != 200 {
		t.Error("post result is not ok")
		return
	}
}

func TestPostJsonData(t *testing.T) {
	url := "http://localhost:8080/test"

	//参数为nil
	req, err := NewJsonPostRequest(url, nil)
	response, err := client.Do(context.TODO(), req)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := DealResponse(response)
	resp := new(responseData)
	json.Unmarshal(result, resp)
	if resp.Code != 400 {
		t.Error("postJsonData result is not ok")
		return
	}

	//参数为空map
	req, err = NewJsonPostRequest(url, make(map[string]interface{}))
	response, err = client.Do(context.TODO(), req)
	if err != nil {
		t.Error(err)
		return
	}
	result, err = DealResponse(response)
	resp = new(responseData)
	json.Unmarshal(result, resp)
	if resp.Code != 200 {
		t.Error("postJsonData result is not ok")
		return
	}

	//参数非空map
	params := map[string]interface{}{
		"name": "hts",
	}
	req, err = NewJsonPostRequest(url, params)
	response, err = client.Do(context.TODO(), req)
	if err != nil {
		t.Error(err)
		return
	}
	resp = new(responseData)
	json.Unmarshal(result, resp)
	if resp.Code != 200 {
		t.Error("postJsonData result is not ok")
		return
	}
}

func TestStringListToMap(t *testing.T) {
	m := StringListToMap([]string{"hts:11 ", "name:", "key", "v: 1:2"})
	_, ok := m["key"]
	if ok {
		t.Error("not right filter")
		return
	}

	val, ok := m["hts"]
	if val != "11" {
		t.Error("not right trim")
		return
	}

	val, ok = m["v"]
	if val != "1:2" {
		t.Error("not right split")
		return
	}
}
