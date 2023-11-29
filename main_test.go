package main_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	dp "github.com/miraddo/dp"
)

func TestProcessRequest(t *testing.T) {

	processor := dp.InitializeUserDataProcessor()
	go processor.ResetRequestCounts()

	testCases := []struct {
		desc     string
		method   string
		body     string
		wantCode int
	}{
		{
			desc:     "Test with GET request method",
			method:   "GET",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Test with invalid JSON",
			method:   "POST",
			body:     `{"id": "1","userID": "user1"`,
			wantCode: http.StatusBadRequest,
		},
		{
			desc:     "Test with proper JSON",
			method:   "POST",
			body:     `{"id": "1","userID": "user1","data": "Hello, World!"}`,
			wantCode: http.StatusOK,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(tC.method, "/data", bytes.NewBuffer([]byte(tC.body)))
			rr := httptest.NewRecorder()

			// Use ProcessRequest to handle request
			processor.ProcessRequest(rr, req)

			if status := rr.Code; status != tC.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tC.wantCode)
			}
		})
	}
}
