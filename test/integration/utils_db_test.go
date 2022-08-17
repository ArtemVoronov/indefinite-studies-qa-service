//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_TASK_NAME_1        string = "Test task 1"
	TEST_TASK_STATE_1       string = entities.TASK_STATE_NEW
	TEST_TASK_NAME_2        string = "Test task 2"
	TEST_TASK_STATE_2       string = entities.TASK_STATE_DONE
	TEST_TASK_NAME_TEMPLATE string = "Test task "

	TEST_TAG_NAME_1        string = "Test tag 1"
	TEST_TAG_STATE_1       string = entities.TAG_STATE_NEW
	TEST_TAG_NAME_2        string = "Test tag 2"
	TEST_TAG_STATE_2       string = entities.TAG_STATE_BLOCKED
	TEST_TAG_NAME_TEMPLATE string = "Test tag "

	TEST_USER_LOGIN_1    string = "Test user 1"
	TEST_USER_EMAIL_1    string = "user1@somewhere.com"
	TEST_USER_PASSWORD_1 string = "Test password1 "
	TEST_USER_ROLE_1     string = entities.USER_ROLE_OWNER
	TEST_USER_STATE_1    string = entities.USER_STATE_NEW
	TEST_USER_LOGIN_2    string = "Test user 2"
	TEST_USER_EMAIL_2    string = "user2@somewhere.com"
	TEST_USER_PASSWORD_2 string = "Tes tpassword 2"
	TEST_USER_ROLE_2     string = entities.USER_ROLE_RESIDENT
	TEST_USER_STATE_2    string = entities.USER_STATE_BLOCKED

	TEST_USER_LOGIN_TEMPLATE   string = "Test user "
	TEST_USER_EMAIL_TEMPLATE   string = "user%v@somewhere.com"
	TEST_USER_PASSORD_TEMPLATE string = "Test password "

	TEST_NOTE_TEXT_1    string = "Test text 1"
	TEST_NOTE_TOPIC_1   string = "Test topic 1"
	TEST_NOTE_TAG_ID_1  int    = 1
	TEST_NOTE_USER_ID_1 int    = 1
	TEST_NOTE_STATE_1   string = entities.NOTE_STATE_NEW
	TEST_NOTE_TEXT_2    string = "Test text 2"
	TEST_NOTE_TOPIC_2   string = "Test topic 2"
	TEST_NOTE_TAG_ID_2  int    = 2
	TEST_NOTE_USER_ID_2 int    = 2
	TEST_NOTE_STATE_2   string = entities.NOTE_STATE_BLOCKED

	TEST_NOTE_TEXT_TEMPLATE  string = "Test text "
	TEST_NOTE_TOPIC_TEMPLATE string = "Test topic "

	TEST_REFRESH_TOKEN_TEMPLATE  string = "Token "
	TEST_REFRESH_TOKEN_1                = "Token 1"
	TEST_REFRESH_TOKEN_2                = "Token 2"
	TEST_REFRESH_TOKEN_USER_ID_1        = 1
	TEST_REFRESH_TOKEN_USER_ID_2        = 2
)

var (
	TEST_REFRESH_TOKEN_EXPIRE_AT_1 = time.Now().Add(time.Hour * 2)
	TEST_REFRESH_TOKEN_EXPIRE_AT_2 = time.Now().Add(time.Hour * 1)
)

type TestAsserts struct {
}

type TestUtilsAsserts interface {
	AssertEqualTasks(t *testing.T, expected entities.Task, actual entities.Task)
	AssertEqualTaskArrays(t *testing.T, expected []entities.Task, actual []entities.Task)
	AssertEqualTags(t *testing.T, expected entities.Tag, actual entities.Tag)
	AssertEqualTagArrays(t *testing.T, expected []entities.Tag, actual []entities.Tag)
	AssertEqualUsers(t *testing.T, expected entities.User, actual entities.User)
	AssertEqualUserArrays(t *testing.T, expected []entities.User, actual []entities.User)
	AssertEqualNotes(t *testing.T, expected entities.Note, actual entities.Note)
	AssertEqualNoteArrays(t *testing.T, expected []entities.Note, actual []entities.Note)
}

func (p *TestAsserts) AssertEqualTasks(t *testing.T, expected entities.Task, actual entities.Task) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.State, actual.State)
}

func (p *TestAsserts) AssertEqualTaskArrays(t *testing.T, expected []entities.Task, actual []entities.Task) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		utils.asserts.AssertEqualTasks(t, expected[i], actual[i])
	}
}

func (p *TestAsserts) AssertEqualTags(t *testing.T, expected entities.Tag, actual entities.Tag) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.State, actual.State)
}

func (p *TestAsserts) AssertEqualTagArrays(t *testing.T, expected []entities.Tag, actual []entities.Tag) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		utils.asserts.AssertEqualTags(t, expected[i], actual[i])
	}
}

func (p *TestAsserts) AssertEqualUsers(t *testing.T, expected entities.User, actual entities.User) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Login, actual.Login)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Password, actual.Password)
	assert.Equal(t, expected.State, actual.State)
}

func (p *TestAsserts) AssertEqualUserArrays(t *testing.T, expected []entities.User, actual []entities.User) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		utils.asserts.AssertEqualUsers(t, expected[i], actual[i])
	}
}

func (p *TestAsserts) AssertEqualNotes(t *testing.T, expected entities.Note, actual entities.Note) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Text, actual.Text)
	assert.Equal(t, expected.Topic, actual.Topic)
	assert.Equal(t, expected.TagId, actual.TagId)
	assert.Equal(t, expected.UserId, actual.UserId)
	assert.Equal(t, expected.State, actual.State)
}

func (p *TestAsserts) AssertEqualNoteArrays(t *testing.T, expected []entities.Note, actual []entities.Note) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		utils.asserts.AssertEqualNotes(t, expected[i], actual[i])
	}
}

func (p *TestAsserts) AssertEqualRefreshTokens(t *testing.T, expected entities.RefreshToken, actual entities.RefreshToken) {
	assert.Equal(t, expected.UserId, actual.UserId)
	assert.Equal(t, expected.Token, actual.Token)
}

type TestEntityGenerators struct {
}

type TestUtilsGenerators interface {
	GenerateTask(id int) entities.Task
	GenerateTaskName(template string, id int) string
	GenerateTag(id int) entities.Tag
	GenerateTagName(template string, id int) string
	GenerateUserLogin(template string, id int) string
	GenerateUserPassword(template string, id int) string
	GenerateUserEmail(template string, id int) string
	GenerateUser(id int) entities.User
	GenerateNoteText(template string, id int) string
	GenerateNoteTopic(template string, id int) string
	GenerateNote(noteId int, userId int, tagId int) entities.Note
}

func (p *TestEntityGenerators) GenerateTask(id int) entities.Task {
	return entities.Task{
		Id:    id,
		Name:  utils.entityGenerators.GenerateTaskName(TEST_TASK_NAME_TEMPLATE, id),
		State: TEST_TASK_STATE_1,
	}
}

func (p *TestEntityGenerators) GenerateTaskName(template string, id int) string {
	return template + strconv.Itoa(id)
}

func (p *TestEntityGenerators) GenerateTag(id int) entities.Tag {
	return entities.Tag{
		Id:    id,
		Name:  utils.entityGenerators.GenerateTagName(TEST_TAG_NAME_TEMPLATE, id),
		State: TEST_TAG_STATE_1,
	}
}

func (p *TestEntityGenerators) GenerateTagName(template string, id int) string {
	return template + strconv.Itoa(id)
}

func (p *TestEntityGenerators) GenerateUserLogin(template string, id int) string {
	return template + strconv.Itoa(id)
}

func (p *TestEntityGenerators) GenerateUserPassword(template string, id int) string {
	return template + strconv.Itoa(id)
}

func (p *TestEntityGenerators) GenerateUserEmail(template string, id int) string {
	return fmt.Sprintf(template, id)
}

func (p *TestEntityGenerators) GenerateUser(id int) entities.User {
	return entities.User{
		Id:       id,
		Login:    utils.entityGenerators.GenerateUserLogin(TEST_USER_LOGIN_TEMPLATE, id),
		Email:    utils.entityGenerators.GenerateUserEmail(TEST_USER_EMAIL_TEMPLATE, id),
		Password: utils.entityGenerators.GenerateUserPassword(TEST_USER_PASSORD_TEMPLATE, id),
		Role:     TEST_USER_ROLE_1,
		State:    TEST_USER_STATE_1,
	}
}

func (p *TestEntityGenerators) GenerateRefreshToken(id int) entities.RefreshToken {
	return entities.RefreshToken{
		UserId:   id,
		Token:    TEST_REFRESH_TOKEN_TEMPLATE + strconv.Itoa(id),
		ExpireAt: time.Now().Add(time.Minute * 30),
	}
}

func (p *TestEntityGenerators) GenerateNoteText(template string, id int) string {
	return template + strconv.Itoa(id)
}

func (p *TestEntityGenerators) GenerateNoteTopic(template string, id int) string {
	return template + strconv.Itoa(id)
}

func (p *TestEntityGenerators) GenerateNote(noteId int, userId int, tagId int) entities.Note {
	return entities.Note{
		Id:     noteId,
		Text:   utils.entityGenerators.GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, noteId),
		Topic:  utils.entityGenerators.GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, noteId),
		TagId:  tagId,
		UserId: userId,
		State:  TEST_USER_STATE_1,
	}
}

type TestUtilsQueries interface {
	CreateTaskInDB(t *testing.T, tx *sql.Tx, ctx context.Context, name string, state string) (int, error)
	CreateTasksInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, nameTemplate string, state string) error
	CreateTagInDB(t *testing.T, tx *sql.Tx, ctx context.Context, name string, state string) (int, error)
	CreateTagsInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, nameTemplate string, state string) error
	CreateUserInDB(t *testing.T, tx *sql.Tx, ctx context.Context, login string, email string, password string, role string, state string) (int, error)
	CreateUsersInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, loginTemplate string, emailTemplate string, passwordTemplate string, role string, state string) error
	CreateNoteInDB(t *testing.T, tx *sql.Tx, ctx context.Context, text string, topic string, tagId int, userId int, state string) (int, error)
	CreateNotesInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, textTemplate string, topicTemplate string, tagId int, userId int, state string) error
}

func CreateTaskInDB(t *testing.T, tx *sql.Tx, ctx context.Context, name string, state string) (int, error) {
	taskId, err := queries.CreateTask(tx, ctx, name, state)
	assert.Nil(t, err)
	assert.NotEqual(t, taskId, -1)
	return taskId, err
}

func CreateTasksInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, nameTemplate string, state string) error {
	var lastErr error
	for i := 1; i <= count; i++ {
		_, err := CreateTaskInDB(t, tx, ctx, utils.entityGenerators.GenerateTaskName(TEST_TASK_NAME_TEMPLATE, i), state)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func CreateTagInDB(t *testing.T, tx *sql.Tx, ctx context.Context, name string, state string) (int, error) {
	tagId, err := queries.CreateTag(tx, ctx, name, state)
	assert.Nil(t, err)
	assert.NotEqual(t, tagId, -1)
	return tagId, err
}

func CreateTagsInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, nameTemplate string, state string) error {
	var lastErr error
	for i := 1; i <= count; i++ {
		_, err := CreateTagInDB(t, tx, ctx, utils.entityGenerators.GenerateTagName(TEST_TAG_NAME_TEMPLATE, i), state)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func CreateUserInDB(t *testing.T, tx *sql.Tx, ctx context.Context, login string, email string, password string, role string, state string) (int, error) {
	userId, err := queries.CreateUser(tx, ctx, login, email, password, role, state)
	assert.Nil(t, err)
	assert.NotEqual(t, userId, -1)
	return userId, err
}

func CreateUsersInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, loginTemplate string, emailTemplate string, passwordTemplate string, role string, state string) error {
	var lastErr error
	for i := 1; i <= count; i++ {
		_, err := CreateUserInDB(t, tx, ctx,
			utils.entityGenerators.GenerateUserLogin(loginTemplate, i),
			utils.entityGenerators.GenerateUserEmail(emailTemplate, i),
			utils.entityGenerators.GenerateUserPassword(passwordTemplate, i),
			role,
			state,
		)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func CreateNoteInDB(t *testing.T, tx *sql.Tx, ctx context.Context, text string, topic string, tagId int, userId int, state string) (int, error) {
	noteId, err := queries.CreateNote(tx, ctx, text, topic, tagId, userId, state)
	assert.Nil(t, err)
	assert.NotEqual(t, noteId, -1)
	return noteId, err
}

func CreateNotesInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, textTemplate string, topicTemplate string, tagId int, userId int, state string) error {
	var lastErr error
	for i := 1; i <= count; i++ {
		_, err := CreateNoteInDB(t, tx, ctx,
			utils.entityGenerators.GenerateNoteText(textTemplate, i),
			utils.entityGenerators.GenerateNoteTopic(topicTemplate, i),
			tagId,
			userId,
			state,
		)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}
