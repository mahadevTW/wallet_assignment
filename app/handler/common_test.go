package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type args struct {
	writer  *httptest.ResponseRecorder
	payload interface{}
}
type TestCase struct {
	name string
	args args
}
type StudentPayload struct {
	Name   string
	RollNo int32
}

func (ts *TestCase) assertMockVars(t *testing.T) {
	resp := ts.args.writer.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	payload := StudentPayload{}
	json.Unmarshal(body, &payload)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, payload, StudentPayload{"mahadev", 23})
}
func Test_respondSuccess(t *testing.T) {
	tests := []*TestCase{
		{
			name: "Write Json response to Response Writer",
			args: args{
				writer:  httptest.NewRecorder(),
				payload: StudentPayload{"mahadev", 23},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respondSuccess(tt.args.writer, tt.args.payload)
		})
		tt.assertMockVars(t)
	}
}

func Test_responseError(t *testing.T) {
	writer := httptest.NewRecorder()
	message := "something went wrong"
	code := http.StatusInternalServerError
	errorResponse := make(map[string]string)

	respondError(writer, code, message)

	resp := writer.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &errorResponse)

	assert.Equal(t, errorResponse["error"], message)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, 500, resp.StatusCode)
}
