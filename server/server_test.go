package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/smart-cs/scheduler-backend/database"
	"github.com/smart-cs/scheduler-backend/server"
	"github.com/stretchr/testify/assert"
)

func TestSchedulesHandlerEmptyBody(t *testing.T) {
	t.Log("hitting schedules endpoint with a course without possible schedules should return an empty body")
	assert := assert.New(t)
	database.LoadLocalDatabase("../database/test-coursedb.json")

	s := server.NewServer()
	assert.NotNil(s, "a new server shouldn't be nil")

	req, err := http.NewRequest("GET", "/schedules", strings.NewReader(url.Values{"courses": {"APSC 210"}}.Encode()))
	assert.Nil(err, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.SchedulesHandler)
	handler.ServeHTTP(rr, req)

	var actual server.StandardResponse
	json.Unmarshal(rr.Body.Bytes(), &actual)
	expected := server.StandardResponse{
		OK:     true,
		Status: http.StatusOK,
		Body:   []interface{}{},
	}
	assert.EqualValues(expected, actual)
}
