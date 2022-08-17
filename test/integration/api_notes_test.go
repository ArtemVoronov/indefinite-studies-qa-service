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
	"github.com/stretchr/testify/assert"
)

var (
	ERROR_NOTE_TEXT_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Text\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_NOTE_TOPIC_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Topic\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_NOTE_TAG_ID_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"TagId\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_NOTE_USER_ID_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"UserId\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_NOTE_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_NOTE_ALL_ARE_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Text\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"Topic\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"TagId\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"UserId\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_NOTE_CREATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to create note. Wrong 'State' value. Possible values: %v", entities.GetPossibleNoteStates())
	ERROR_NOTE_UPDATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to update note. Wrong 'State' value. Possible values: %v", entities.GetPossibleNoteStates())
)

func TestApiNoteGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetNote("1")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		id := "1"
		text := utils.entityGenerators.GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, 1)
		topic := utils.entityGenerators.GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, 1)
		tagId := 1
		userId := 1
		state := entities.NOTE_STATE_NEW
		expectedBody := "{" +
			"\"Id\":" + id + "," +
			"\"Text\":\"" + text + "\"," +
			"\"Topic\":\"" + topic + "\"," +
			"\"TagId\":" + strconv.Itoa(tagId) + "," +
			"\"UserId\":" + strconv.Itoa(userId) + "," +
			"\"State\":\"" + state + "\"" +
			"}"

		httpStatusCode, body, _ := testHttpClient.CreateNote(text, topic, tagId, userId, state)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, id, body)

		httpStatusCode, body = testHttpClient.GetNote(id)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetNote("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetNote("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := testHttpClient.GetNote("")

		assert.Equal(t, http.StatusMovedPermanently, httpStatusCode)
		assert.Equal(t, "<a href=\"/notes\">Moved Permanently</a>.\n\n", body)
	})))
}

func TestApiNoteGetAll(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":10,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			text := utils.entityGenerators.GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, i)
			topic := utils.entityGenerators.GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, i)
			tagId := i
			userId := i
			state := entities.NOTE_STATE_NEW

			testHttpClient.CreateNote(text, topic, tagId, userId, state)
			expectedBody += "{" +
				"\"Id\":" + id + "," +
				"\"Text\":\"" + text + "\"," +
				"\"Topic\":\"" + topic + "\"," +
				"\"TagId\":" + strconv.Itoa(tagId) + "," +
				"\"UserId\":" + strconv.Itoa(userId) + "," +
				"\"State\":\"" + state + "\"" +
				"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		httpStatusCode, body, _ := testHttpClient.GetNotes(nil, nil)

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
		httpStatusCode, body, _ := testHttpClient.GetNotes(nil, nil)

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
			text := utils.entityGenerators.GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, i)
			topic := utils.entityGenerators.GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, i)
			tagId := i
			userId := i
			state := entities.NOTE_STATE_NEW
			expectedBody += "{" +
				"\"Id\":" + id + "," +
				"\"Text\":\"" + text + "\"," +
				"\"Topic\":\"" + topic + "\"," +
				"\"TagId\":" + strconv.Itoa(tagId) + "," +
				"\"UserId\":" + strconv.Itoa(userId) + "," +
				"\"State\":\"" + state + "\"" +
				"}"
			if i != 5 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		for i := 1; i <= 10; i++ {
			text := utils.entityGenerators.GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, i)
			topic := utils.entityGenerators.GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, i)
			tagId := i
			userId := i
			state := entities.NOTE_STATE_NEW

			testHttpClient.CreateNote(text, topic, tagId, userId, state)
		}

		httpStatusCode, body, _ := testHttpClient.GetNotes(5, 0)

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
			text := utils.entityGenerators.GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, i)
			topic := utils.entityGenerators.GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, i)
			tagId := i
			userId := i
			state := entities.NOTE_STATE_NEW
			expectedBody += "{" +
				"\"Id\":" + id + "," +
				"\"Text\":\"" + text + "\"," +
				"\"Topic\":\"" + topic + "\"," +
				"\"TagId\":" + strconv.Itoa(tagId) + "," +
				"\"UserId\":" + strconv.Itoa(userId) + "," +
				"\"State\":\"" + state + "\"" +
				"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		for i := 1; i <= 10; i++ {
			text := utils.entityGenerators.GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, i)
			topic := utils.entityGenerators.GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, i)
			tagId := i
			userId := i
			state := entities.NOTE_STATE_NEW

			testHttpClient.CreateNote(text, topic, tagId, userId, state)
		}

		httpStatusCode, body, _ := testHttpClient.GetNotes(50, 5)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
}

func TestApiNoteCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, "1", body)
	})))
	t.Run("WrongInput: Missed 'Text'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(nil, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TEXT_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Topic'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, nil, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TOPIC_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'TagId'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, nil, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TAG_ID_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'UserId'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, nil, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_USER_ID_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Text' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote("", TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TEXT_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Topic' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, "", TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TOPIC_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'TagId' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, "", TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'UserId' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, "", TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, "")
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Text' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Topic' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, 1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'TagId' is not an integer", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, "1", TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'UserId' is not an integer", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, "1", TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, 1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_NOTE_CREATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("DeletedCase: try to create as deleted", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, entities.NOTE_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
}

func TestApiNoteUpdate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		id := "1"
		text := TEST_NOTE_TEXT_2
		topic := TEST_NOTE_TOPIC_2
		tagId := TEST_NOTE_TAG_ID_2
		userId := TEST_NOTE_USER_ID_2
		state := TEST_NOTE_STATE_2
		expectedBody := "{" +
			"\"Id\":" + id + "," +
			"\"Text\":\"" + text + "\"," +
			"\"Topic\":\"" + topic + "\"," +
			"\"TagId\":" + strconv.Itoa(tagId) + "," +
			"\"UserId\":" + strconv.Itoa(userId) + "," +
			"\"State\":\"" + state + "\"" +
			"}"

		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, id, body)

		httpStatusCode, body, _ = testHttpClient.UpdateNote(id, text, topic, tagId, userId, state)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = testHttpClient.GetNote(id)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)

	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("text", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("2.15", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: Missed 'Text'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", nil, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TEXT_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Topic'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, nil, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TOPIC_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'TagId'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, nil, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TAG_ID_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'UserId'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, nil, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_USER_ID_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Text' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", "", TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TEXT_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Topic' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, "", TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_TOPIC_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'TagId' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, "", TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'UserId' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, "", TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, "")
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NOTE_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Text' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", 1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Topic' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, 1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'TagId' is not an integer", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, "1", TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'UserId' is not an integer", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, "1", TEST_NOTE_STATE_1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, 1)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_NOTE_UPDATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.UpdateNote("1", TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: find deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = testHttpClient.DeleteNote(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body, _ = testHttpClient.UpdateNote(expectedId, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, TEST_NOTE_TAG_ID_2, TEST_NOTE_USER_ID_2, TEST_NOTE_STATE_2)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to mark as deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = testHttpClient.UpdateNote(expectedId, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, TEST_NOTE_TAG_ID_2, TEST_NOTE_USER_ID_2, entities.NOTE_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
	t.Run("MultipleUpdateCase", RunWithRecreateDB((func(t *testing.T) {
		id := "1"
		expectedBody := "{" +
			"\"Id\":" + id + "," +
			"\"Text\":\"" + TEST_NOTE_TEXT_2 + "\"," +
			"\"Topic\":\"" + TEST_NOTE_TOPIC_2 + "\"," +
			"\"TagId\":" + strconv.Itoa(TEST_NOTE_TAG_ID_2) + "," +
			"\"UserId\":" + strconv.Itoa(TEST_NOTE_USER_ID_2) + "," +
			"\"State\":\"" + TEST_NOTE_STATE_2 + "\"" +
			"}"

		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, id, body)

		for i := 1; i <= 3; i++ {
			httpStatusCode, body, _ = testHttpClient.UpdateNote(id, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, TEST_NOTE_TAG_ID_2, TEST_NOTE_USER_ID_2, TEST_NOTE_STATE_2)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, "\""+api.DONE+"\"", body)

			httpStatusCode, body = testHttpClient.GetNote(id)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, expectedBody, body)
		}
	})))
}

func TestApiNoteDelete(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = testHttpClient.DeleteNote(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = testHttpClient.GetNote(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteNote("")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteNote("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := testHttpClient.DeleteNote("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("MultipleDeleteCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := testHttpClient.CreateNote(TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, TEST_NOTE_TAG_ID_1, TEST_NOTE_USER_ID_1, TEST_NOTE_STATE_1)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = testHttpClient.DeleteNote(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body, _ = testHttpClient.DeleteNote(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
}
