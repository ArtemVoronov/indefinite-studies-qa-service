//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/api"
	"github.com/stretchr/testify/assert"
)

var (
	ERROR_TASK_NAME_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TASK_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TASK_NAME_AND_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TASK_CREATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to create task. Wrong 'State' value. Possible values: %v", entities.GetPossibleTaskStates())
	ERROR_TASK_UPDATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to update task. Wrong 'State' value. Possible values: %v", entities.GetPossibleTaskStates())
)

func TestApiTaskGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTask("1")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Task 1"
		expectedState := entities.TASK_STATE_NEW
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		httpStatusCode, body, _ := testHttpClient.CreateTask(expectedName, expectedState)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body = testHttpClient.GetTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTask("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTask("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTask("")

		assert.Equal(t, http.StatusMovedPermanently, httpStatusCode)
		assert.Equal(t, "<a href=\"/tasks\">Moved Permanently</a>.\n\n", body)
	})))
}

func TestApiTaskGetAll(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":10,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
			state := entities.TASK_STATE_NEW

			testHttpClient.CreateTask(name, state)
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		httpStatusCode, body, _ := testHttpClient.GetTasks(nil, nil)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("EmptyResult", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":0,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":[]"
		expectedBody += "}"
		httpStatusCode, body, _ := testHttpClient.GetTasks(nil, nil)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("LimitCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":5,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":5,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 5; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
			state := entities.TASK_STATE_NEW
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 5 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
			state := entities.TASK_STATE_NEW
			testHttpClient.CreateTask(name, state)
		}

		httpStatusCode, body, _ := testHttpClient.GetTasks(5, 0)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("OffsetCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":5,"
		expectedBody += "\"Offset\":5,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 6; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
			state := entities.TASK_STATE_NEW
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
			state := entities.TASK_STATE_NEW
			testHttpClient.CreateTask(name, state)
		}

		httpStatusCode, body, _ := testHttpClient.GetTasks(50, 5)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
}

func TestApiTaskCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, "1", body)
	})))
	t.Run("WrongInput: Missed 'Name'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask(nil, entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask("Test Task 1", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask(nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask(1, entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask("Test Task 1", 1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Name' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask("", entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask("Test Task 1", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask("Test Task 1", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_TASK_CREATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to create as deleted", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
}

func TestApiTaskUpdate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Task 2"
		expectedState := entities.TASK_STATE_DONE
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"
		testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = testHttpClient.GetTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)

	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("text", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("2.15", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: Missed 'Name'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", nil, entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", "Test Task 2", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", 10000, entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", "Test Task 2", 10000)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Name' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", "", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", "Test Task 2", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TASK_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", "Test Task 2", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_TASK_UPDATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTask("1", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: find deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)
		testHttpClient.DeleteTask(expectedId)

		httpStatusCode, body, _ := testHttpClient.UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to mark as deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)
		httpStatusCode, body, _ := testHttpClient.UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "2"

		testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)
		testHttpClient.CreateTask("Test Task 2", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.UpdateTask(expectedId, "Test Task 1", entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("MultipleUpdateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Task 2"
		expectedState := entities.TASK_STATE_DONE
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		for i := 1; i <= 3; i++ {
			httpStatusCode, body, _ := testHttpClient.UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DONE)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, "\""+api.DONE+"\"", body)

			httpStatusCode, body = testHttpClient.GetTask(expectedId)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, expectedBody, body)
		}
	})))
}

func TestApiTaskDelete(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.DeleteTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = testHttpClient.GetTask(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteTask("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteTask("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteTask("")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("MultipleDeleteCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.DeleteTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body, _ = testHttpClient.DeleteTask(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
}
