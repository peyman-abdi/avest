package avest

import (
	"time"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"errors"
	"testing"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
)

var orderInt = 1
type TestRouteModel struct {
	ID int64
	MyTest *string
	MyInt *int
	CreatedAt *time.Time
}
func (t *TestRouteModel) PrimaryKey() string {
	return "id"
}
func (t *TestRouteModel) TableName() string {
	return "tests"
}
type TestRouteMigratable struct {
}
func (t *TestRouteMigratable) Up(migrator core.Migrator) error {
	var err error
	if err = migrator.AutoMigrate(&TestRouteModel{}); err != nil {
		return err
	}
	return nil
}
func (t *TestRouteMigratable) Down(migrator core.Migrator) error {
	var err error
	if err = migrator.DropTableIfExists(&TestRouteModel{}); err != nil {
		return err
	}
	return nil
}
type TestRouteModule struct {
	S core.Services
}
var _ core.Module = (*TestRouteModule)(nil)
func (t *TestRouteModule) Title() string       { return "TestMigrationModule" }
func (t *TestRouteModule) Description() string { return "Test module" }
func (t *TestRouteModule) Version() string     { return "1.01" }
func (t *TestRouteModule) Activated() bool     { return true }
func (t *TestRouteModule) Installed() bool     { return true }
func (t *TestRouteModule) Deactivated()        { }
func (t *TestRouteModule) Purged()             { }
func (t *TestRouteModule) Migrations() []core.Migratable {
	return []core.Migratable {
		new(TestRouteMigratable),
	}
}
func (t *TestRouteModule) Routes() []*core.Route {
	return []*core.Route {
		{
			Group: "/api/tests",
			MiddleWares: []string {
				"oauth",
			},
			Methods: core.ANY,
			Urls: []string{"/id/<id:\\d+>/str/<name>"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				request.SetValue("route:test", orderInt); orderInt++
				values := request.GetAll("route:test", "middleware:auth", "group:api", "group:tests", "id", "name")
				response.SuccessJSON(values)
				return nil
			},
		},
		{
			Group: "/api",
			MiddleWares: []string {
				"oauth",
			},
			Methods: core.ANY,
			Urls: []string{"/", ""},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				response.SuccessString("hello api")
				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Methods: core.ANY,
			Urls: []string{"/", ""},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				response.SuccessString("hello world")
				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Group: "/api",
			Methods: core.GET,
			Urls: []string{"/models"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				var models []*TestRouteModel
				t.S.Repository().Query(&TestRouteModel{}).Get(&models)
				response.SuccessJSON(map[string]interface{} {
					"data": models,
					"count": len(models),
				})
				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Group: "/api",
			Methods: core.GET,
			Urls: []string{"/models/<id:\\d+>"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				var model *TestRouteModel
				t.S.Repository().Query(&TestRouteModel{}).Where("id = ?", request.GetValue("id")).Get(&model)
				response.SuccessJSON(map[string]interface{} {
					"data": model,
				})
				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Group: "/api",
			Methods: core.PUT,
			Urls: []string{"/models/<id:\\d+>"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				var model *TestRouteModel
				t.S.Repository().Query(&TestRouteModel{}).Where("id = ?", request.GetValue("id")).GetFirst(&model)
				if model == nil {
					return errors.New("object not found")
				}

				vals := request.GetAll("text", "int")
				t.S.Repository().Query(&TestRouteModel{}).Update(model, map[string]interface{} {
					"my_test": StringRefOrNil(vals["text"]),
					"my_int": IntRefOrNil(vals["int"]),
				})

				response.SuccessJSON(map[string]interface{} {
					"data": model,
				})
				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Group: "/api",
			Methods: core.POST,
			Urls: []string{"/models"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				vals := request.GetAll("text", "int")

				var model = &TestRouteModel{
					MyTest: StringRefOrNil(vals["text"]),
					MyInt: IntRefOrNil(vals["int"]),
				}

				t.S.Repository().Insert(model)

				response.SuccessJSON(map[string]interface{} {
					"data": model,
				})
				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Group: "/api",
			Methods: core.DELETE,
			Urls: []string{"/models/<id:\\d+>"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				var model *TestRouteModel
				t.S.Repository().Query(&TestRouteModel{}).Where("id = ?", request.GetValue("id")).GetFirst(&model)
				if model == nil {
					return errors.New("object not found")
				}
				t.S.Repository().DeleteEntity(model)
				response.SuccessString("deleted")

				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Methods: core.GET,
			Urls: []string{"/home"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				response.View("home", nil)
				return nil
			},
		},
		{
			MiddleWares: []string {
				"oauth",
			},
			Methods: core.GET,
			Urls: []string{"/error"},
			Verify: nil,
			Handle: func(request core.Request, response core.Response) error {
				response.View("error", nil)
				return nil
			},
		},
	}
}
func (t *TestRouteModule) MiddleWares() []*core.MiddleWare {
	return []*core.MiddleWare {
		{
			Name: "oauth",
			Handler: func(request core.Request, response core.Response) error {
				request.SetValue("middleware:auth", orderInt); orderInt++
				return nil
			},
		},
	}
}
func (t *TestRouteModule) GroupsHandlers() []*core.RouteGroup {
	return[]*core.RouteGroup {
		{
			Name: "tests",
			Handler: func(request core.Request, response core.Response) error {
				request.SetValue("group:tests", orderInt); orderInt++
				return nil
			},
		},
		{
			Name: "api",
			Handler: func(request core.Request, response core.Response) error {
				request.SetValue("group:api", orderInt); orderInt++
				return nil
			},
		},
	}
}
func (t *TestRouteModule) Templates() []*core.Template {
	return []*core.Template {
		{
			Name:"home",
			Path:"home.jet",
		},
		{
			Name:"error",
			Path:"error.jet",
		},
	}
}

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
		t.Errorf("Response body not equal to expected: %s != %s", body, expect)
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
