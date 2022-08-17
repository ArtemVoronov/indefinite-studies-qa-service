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
	ERROR_TAG_NAME_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TAG_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TAG_NAME_AND_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TAG_CREATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to create tag. Wrong 'State' value. Possible values: %v", entities.GetPossibleTagStates())
	ERROR_TAG_UPDATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to update tag. Wrong 'State' value. Possible values: %v", entities.GetPossibleTagStates())
)

func TestApiTagGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTag("1")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Tag 1"
		expectedState := entities.TAG_STATE_NEW
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		httpStatusCode, body, _ := testHttpClient.CreateTag(expectedName, expectedState)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body = testHttpClient.GetTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTag("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTag("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetTag("")

		assert.Equal(t, http.StatusMovedPermanently, httpStatusCode)
		assert.Equal(t, "<a href=\"/tags\">Moved Permanently</a>.\n\n", body)
	})))
}

func TestApiTagGetAll(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":10,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW

			testHttpClient.CreateTag(name, state)
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		httpStatusCode, body, _ := testHttpClient.GetTags(nil, nil)

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
		httpStatusCode, body, _ := testHttpClient.GetTags(nil, nil)

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
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
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
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
			testHttpClient.CreateTag(name, state)
		}

		httpStatusCode, body, _ := testHttpClient.GetTags(5, 0)

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
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
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
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
			testHttpClient.CreateTag(name, state)
		}

		httpStatusCode, body, _ := testHttpClient.GetTags(50, 5)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
}

func TestApiTagCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, "1", body)
	})))
	t.Run("WrongInput: Missed 'Name'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag(nil, entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag("Test Tag 1", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag(nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag(1, entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag("Test Tag 1", 1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Name' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag("", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag("Test Tag 1", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag("Test Tag 1", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_TAG_CREATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to create as deleted", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
}

func TestApiTagUpdate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Tag 2"
		expectedState := entities.TAG_STATE_BLOCKED
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"
		testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = testHttpClient.GetTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)

	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("text", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("2.15", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: Missed 'Name'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", nil, entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", "Test Tag 2", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", 10000, entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", "Test Tag 2", 10000)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Name' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", "", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", "Test Tag 2", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", "Test Tag 2", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_TAG_UPDATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateTag("1", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: find deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)
		testHttpClient.DeleteTag(expectedId)

		httpStatusCode, body, _ := testHttpClient.UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to mark as deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)
		httpStatusCode, body, _ := testHttpClient.UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "2"

		testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)
		testHttpClient.CreateTag("Test Tag 2", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.UpdateTag(expectedId, "Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("MultipleUpdateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Tag 2"
		expectedState := entities.TAG_STATE_BLOCKED
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		for i := 1; i <= 3; i++ {
			httpStatusCode, body, _ := testHttpClient.UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_BLOCKED)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, "\""+api.DONE+"\"", body)

			httpStatusCode, body = testHttpClient.GetTag(expectedId)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, expectedBody, body)
		}
	})))
}

func TestApiTagDelete(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.DeleteTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = testHttpClient.GetTag(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteTag("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteTag("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteTag("")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("MultipleDeleteCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		testHttpClient.CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := testHttpClient.DeleteTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body, _ = testHttpClient.DeleteTag(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
}
