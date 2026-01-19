package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"monolith/internal/service/auth"

	"github.com/labstack/echo/v4"
)

type TestContext struct {
	Echo     *echo.Echo
	Context  echo.Context
	Request  *http.Request
	Recorder *httptest.ResponseRecorder
}

func NewTestContext(method, path string, body any) *TestContext {
	e := echo.New()

	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return &TestContext{
		Echo:     e,
		Context:  c,
		Request:  req,
		Recorder: rec,
	}
}

func (tc *TestContext) SetUser(user *auth.AuthUser) {
	tc.Context.Set("user", user)
}

func (tc *TestContext) SetCookie(name, value string) {
	tc.Request.AddCookie(&http.Cookie{
		Name:  name,
		Value: value,
	})
}

func (tc *TestContext) SetPathParams(params map[string]string) {
	names := make([]string, 0, len(params))
	values := make([]string, 0, len(params))
	for name, value := range params {
		names = append(names, name)
		values = append(values, value)
	}
	tc.Context.SetParamNames(names...)
	tc.Context.SetParamValues(values...)
}

func (tc *TestContext) ResponseBody() map[string]any {
	var result map[string]any
	_ = json.Unmarshal(tc.Recorder.Body.Bytes(), &result)
	return result
}

func (tc *TestContext) ResponseBodyArray() []map[string]any {
	var result []map[string]any
	_ = json.Unmarshal(tc.Recorder.Body.Bytes(), &result)
	return result
}

func (tc *TestContext) GetCookie(name string) *http.Cookie {
	for _, cookie := range tc.Recorder.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

type MockValidator struct{}

func (v *MockValidator) Validate(i any) error {
	return nil
}

func NewEchoWithValidator() *echo.Echo {
	e := echo.New()
	e.Validator = &MockValidator{}
	return e
}

func NewTestContextWithValidator(method, path string, body any) *TestContext {
	tc := NewTestContext(method, path, body)
	tc.Echo.Validator = &MockValidator{}
	return tc
}
