package avest

import (
	"testing"
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"encoding/json"
	"bytes"
)

func TestGetRequest(t *testing.T, url string, expect string) {
	res, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	if fmt.Sprintf("%s", body) != expect {
		t.Errorf("Response body not equal %s != %s", body, expect)
	}
}
func TestCookies(t *testing.T, url string, expect map[string]string) {
	res, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	res.Body.Close()

	cookies := res.Cookies()
	success := true
	for name, val := range expect {
		for _, cookie := range cookies {
			if cookie.Name == name && cookie.Value == val {
				break
			}
		}

		t.Errorf("cookie with name %s and value %v not found", name, val)
		success = false
	}
	if !success {
		fmt.Println(fmt.Sprintf("%v", cookies))
	}
}
func TestSession(t *testing.T, first string, second string, expect map[string]string)  {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	res, err := client.Get(first)
	if err != nil {
		t.Error(err)
	}
	res.Body.Close()

	cookies := res.Cookies()

	req, err := http.NewRequest("GET", second, nil)
	if err != nil {
		t.Error(err)
		return
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	response, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	var out map[string]interface{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		t.Error(err)
	}

	for key, val := range expect {
		if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", out[key]) {
			t.Errorf("Parameter %s is not as expected: %v != %v [%v]", key, val, out[key], out)
		}
	}
}
func TestHTMLRequest(t *testing.T, url string, see string) {
	res, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(fmt.Sprintf("%s", body), see) {
		t.Errorf("Response body not equal to expected: %s != %s", body, see)
	}
}
func TestGetJSONRequest(t *testing.T, url string, expect map[string]interface{}) {
	res, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	var out map[string]interface{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		t.Error(err)
	}

	for key, val := range expect {
		if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", out[key]) {
			t.Errorf("Parameter %s is not as expected: %v != %v [%v]", key, val, out[key], out)
		}
	}
}
func TestPostJSONRequest(t *testing.T, url string, params map[string]interface{}, expect map[string]interface{}) {
	js, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(js))
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	var out map[string]interface{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		t.Error(err)
	}

	out, ok := out["data"].(map[string]interface{})
	if ok {
		for key, val := range expect {
			if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", out[key]) {
				t.Errorf("Parameter %s is not as expected: %v != %v [%v]", key, val, out[key], out)
			}
		}
	}
}
func TestPutJSONRequest(t *testing.T, url string, params map[string]interface{}, expect map[string]interface{}) {
	js, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(js))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	var out map[string]interface{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		t.Error(err)
	}

	out, ok := out["data"].(map[string]interface{})
	if ok {
		for key, val := range expect {
			if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", out[key]) {
				t.Errorf("Parameter %s is not as expected: %v != %v [%v]", key, val, out[key], out)
			}
		}
	}
}
func TestPutJSONRequestString(t *testing.T, url string, params map[string]interface{}, expect string) {
	js, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(js))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}


	if fmt.Sprintf("%s", body) != expect {
		t.Errorf("Response body not equal to expected: %s != %s", body, expect)
	}
}
func TestDeleteRequestString(t *testing.T, url string, expect string) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Error(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}


	if fmt.Sprintf("%s", body) != expect {
		t.Errorf("Response body not equal to expected: %s != %s", body, expect)
	}
}

