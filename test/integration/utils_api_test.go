//go:build integration
// +build integration

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var testHttpClient TestHttpClient = TestHttpClient{}

type TagsApi interface {
	CreateTag(name any, state any) (int, string, error)
	GetTag(id string) (int, string)
	GetTags(limit any, offset any) (int, string, error)
	UpdateTag(id any, name any, state any) (int, string, error)
	DeleteTag(id any) (int, string, error)
}

type TasksApi interface {
	CreateTask(name any, state any) (int, string, error)
	GetTask(id string) (int, string)
	GetTasks(limit any, offset any) (int, string, error)
	UpdateTask(id any, name any, state any) (int, string, error)
	DeleteTask(id any) (int, string, error)
}

type UsersApi interface {
	CreateUser(login any, email any, password any, role any, state any) (int, string, error)
	GetUser(id string) (int, string)
	GetUsers(limit any, offset any) (int, string, error)
	UpdateUser(id any, login any, email any, password any, role any, state any) (int, string, error)
	DeleteUser(id any) (int, string, error)
}

type NotesApi interface {
	CreateNote(text any, topic any, tagId any, userId any, state any) (int, string, error)
	GetNote(id string) (int, string)
	GetNotes(limit any, offset any) (int, string, error)
	UpdateNote(id any, text any, topic any, tagId any, userId any, state any) (int, string, error)
	DeleteNote(id any) (int, string, error)
}

type AuthApi interface {
	Authenicate(email any, password any) (int, string, error)
	RefreshToken(refreshToken any) (int, string, error)
}

type PingApi interface {
	Ping() (int, string, error)
	SafePing() (int, string, error)
}

type TestApi interface {
	TagsApi
	TasksApi
	UsersApi
	NotesApi
	AuthApi
	PingApi
}

type TestHttpClient struct {
}

func (p *TestHttpClient) CreateTask(name any, state any) (int, string, error) {
	body, err := CreateTaskPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) GetTask(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func (p *TestHttpClient) GetTasks(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks"+queryParams, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) UpdateTask(id any, name any, state any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}
	body, err := CreateTaskPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tasks"+idParam, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) DeleteTask(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tasks"+idParam, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) CreateTag(name any, state any) (int, string, error) {
	body, err := CreateTagPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tags", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) GetTag(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tags/"+id, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func (p *TestHttpClient) GetTags(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tags"+queryParams, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) UpdateTag(id any, name any, state any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}
	body, err := CreateTagPutOrPostBody(name, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tags"+idParam, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) DeleteTag(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tags"+idParam, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) CreateUser(login any, email any, password any, role any, state any) (int, string, error) {
	body, err := CreateUserPutOrPostBody(login, email, password, role, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) GetUser(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+id, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func (p *TestHttpClient) GetUsers(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users"+queryParams, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) UpdateUser(id any, login any, email any, password any, role any, state any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}
	body, err := CreateUserPutOrPostBody(login, email, password, role, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users"+idParam, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) DeleteUser(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users"+idParam, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) CreateNote(text any, topic any, tagId any, userId any, state any) (int, string, error) {
	body, err := CreateNotePutOrPostBody(text, topic, tagId, userId, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) GetNote(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/notes/"+id, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func (p *TestHttpClient) GetNotes(limit any, offset any) (int, string, error) {
	queryParams, err := CreateLimitAndOffsetQueryParams(limit, offset)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/notes"+queryParams, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) UpdateNote(id any, text any, topic any, tagId any, userId any, state any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}
	body, err := CreateNotePutOrPostBody(text, topic, tagId, userId, state)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/notes"+idParam, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) DeleteNote(id any) (int, string, error) {
	idParam, err := ParseForPathParam("id", id)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/notes"+idParam, nil)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) Authenicate(email any, password any) (int, string, error) {
	body, err := CreateAuthenicateBody(email, password)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) RefreshToken(token any) (int, string, error) {
	body, err := CreateRefreshTokenBody(token)
	if err != nil {
		return -1, "", err
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/refresh-token", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) Ping() (int, string, error) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Content-Type", "application/json")
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func (p *TestHttpClient) SafePing(accessToken string) (int, string, error) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/safe-ping", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	TestRouter.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func ParseForJsonBody(paramName string, paramValue any) (string, error) {
	result := ""
	switch paramType := paramValue.(type) {
	case int:
		result = "\"" + paramName + "\": " + strconv.Itoa(paramValue.(int))
	case string:
		result = "\"" + paramName + "\": \"" + paramValue.(string) + "\""
	case nil:
		result = ""
	default:
		return "", fmt.Errorf("unkown type for '%s': %v", paramName, paramType)
	}
	return result, nil
}

func ParseForPathParam(paramName string, paramValue any) (string, error) {
	result := ""
	switch paramType := paramValue.(type) {
	case int:
		result = "/" + strconv.Itoa(paramValue.(int))
	case string:
		result = "/" + paramValue.(string)
	case nil:
		result = ""
	default:
		return "", fmt.Errorf("unkown type for '%s': %v", paramName, paramType)
	}
	return result, nil
}

func ParseForQueryParam(paramName string, paramValue any) (string, error) {
	result := ""
	switch paramType := paramValue.(type) {
	case int:
		result = paramName + "=" + strconv.Itoa(paramValue.(int))
	case nil:
		result = ""
	default:
		return "", fmt.Errorf("unkown type for '%s': %v", paramName, paramType)
	}
	return result, nil
}

func CreateLimitAndOffsetQueryParams(limit any, offset any) (string, error) {
	limitQueryParam, err := ParseForQueryParam("limit", limit)
	if err != nil {
		return "", err
	}
	offsetQueryParam, err := ParseForQueryParam("offset", offset)
	if err != nil {
		return "", err
	}

	queryParams := ""
	if limitQueryParam != "" && offsetQueryParam != "" {
		queryParams += "?" + limitQueryParam + "&" + offsetQueryParam
	} else if limitQueryParam != "" {
		queryParams += "?" + limitQueryParam
	} else if offsetQueryParam != "" {
		queryParams += "?" + offsetQueryParam
	}

	return queryParams, nil
}

func CreateTaskPutOrPostBody(name any, state any) (string, error) {
	nameField, err := ParseForJsonBody("Name", name)
	if err != nil {
		return "", err
	}
	stateField, err := ParseForJsonBody("State", state)
	if err != nil {
		return "", err
	}

	result := "{"
	if nameField != "" {
		result += nameField + ","
	}
	if stateField != "" {
		result += stateField + ","
	}
	if len(result) != 1 {
		result = result[:len(result)-1]
	}
	result += "}"
	return result, nil
}

func CreateTagPutOrPostBody(name any, state any) (string, error) {
	nameField, err := ParseForJsonBody("Name", name)
	if err != nil {
		return "", err
	}
	stateField, err := ParseForJsonBody("State", state)
	if err != nil {
		return "", err
	}

	result := "{"
	if nameField != "" {
		result += nameField + ","
	}
	if stateField != "" {
		result += stateField + ","
	}
	if len(result) != 1 {
		result = result[:len(result)-1]
	}
	result += "}"
	return result, nil
}

func CreateUserPutOrPostBody(login any, email any, password any, role any, state any) (string, error) {
	loginField, err := ParseForJsonBody("Login", login)
	if err != nil {
		return "", err
	}
	emailField, err := ParseForJsonBody("Email", email)
	if err != nil {
		return "", err
	}
	passwordField, err := ParseForJsonBody("Password", password)
	if err != nil {
		return "", err
	}
	roleField, err := ParseForJsonBody("Role", role)
	if err != nil {
		return "", err
	}
	stateField, err := ParseForJsonBody("State", state)
	if err != nil {
		return "", err
	}

	result := "{"
	if loginField != "" {
		result += loginField + ","
	}
	if emailField != "" {
		result += emailField + ","
	}
	if passwordField != "" {
		result += passwordField + ","
	}
	if roleField != "" {
		result += roleField + ","
	}
	if stateField != "" {
		result += stateField + ","
	}
	if len(result) != 1 {
		result = result[:len(result)-1]
	}
	result += "}"
	return result, nil
}

func CreateNotePutOrPostBody(text any, topic any, tagId any, userId any, state any) (string, error) {
	textField, err := ParseForJsonBody("Text", text)
	if err != nil {
		return "", err
	}
	topicField, err := ParseForJsonBody("Topic", topic)
	if err != nil {
		return "", err
	}
	tagIdField, err := ParseForJsonBody("TagId", tagId)
	if err != nil {
		return "", err
	}
	userIdField, err := ParseForJsonBody("UserId", userId)
	if err != nil {
		return "", err
	}
	stateField, err := ParseForJsonBody("State", state)
	if err != nil {
		return "", err
	}

	result := "{"
	if textField != "" {
		result += textField + ","
	}
	if topicField != "" {
		result += topicField + ","
	}
	if tagIdField != "" {
		result += tagIdField + ","
	}
	if userIdField != "" {
		result += userIdField + ","
	}
	if stateField != "" {
		result += stateField + ","
	}
	if len(result) != 1 {
		result = result[:len(result)-1]
	}
	result += "}"
	return result, nil
}

func CreateAuthenicateBody(email any, password any) (string, error) {
	emailField, err := ParseForJsonBody("Email", email)
	if err != nil {
		return "", err
	}
	passwordField, err := ParseForJsonBody("Password", password)
	if err != nil {
		return "", err
	}
	result := "{"
	if emailField != "" {
		result += emailField + ","
	}
	if passwordField != "" {
		result += passwordField + ","
	}
	if len(result) != 1 {
		result = result[:len(result)-1]
	}
	result += "}"
	return result, nil
}

func CreateRefreshTokenBody(token any) (string, error) {
	refreshTokenField, err := ParseForJsonBody("RefreshToken", token)
	if err != nil {
		return "", err
	}
	result := "{"
	if refreshTokenField != "" {
		result += refreshTokenField + ","
	}
	if len(result) != 1 {
		result = result[:len(result)-1]
	}
	result += "}"
	return result, nil
}
