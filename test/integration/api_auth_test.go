//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/auth"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestApiAuthLogin(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var result auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &result)

		assert.Nil(t, err)

		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
		assert.NotNil(t, result.AccessTokenExpiredAt)
		assert.NotNil(t, result.RefreshTokenExpiredAt)
		assert.NotEqual(t, "", result.AccessToken)
		assert.NotEqual(t, "", result.RefreshToken)
		assert.NotEqual(t, "", result.AccessTokenExpiredAt)
		assert.NotEqual(t, "", result.RefreshTokenExpiredAt)
		assert.Equal(t, 236, len(result.AccessToken))
		assert.Equal(t, 238, len(result.RefreshToken))
		assert.NotEqual(t, result.AccessToken, result.RefreshTokenExpiredAt)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			record, err := queries.GetRefreshTokenByToken(tx, ctx, result.RefreshToken)

			assert.NotNil(t, record)
			assert.Equal(t, record.Token, result.RefreshToken)
			assert.Equal(t, record.UserId, user.Id)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, result.AccessToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()
	})))
	t.Run("RepeatAuthenication", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var authenication1 auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &authenication1)

		assert.Nil(t, err)

		time.Sleep(1 * time.Second) // tokens generated based on time.Now(), sometimes we have equal values

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var authenication2 auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &authenication2)

		assert.Nil(t, err)

		assert.NotEqual(t, authenication1.AccessToken, authenication2.AccessToken)
		assert.NotEqual(t, authenication1.RefreshToken, authenication2.RefreshToken)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, authenication1.RefreshToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			record, err := queries.GetRefreshTokenByToken(tx, ctx, authenication2.RefreshToken)

			assert.NotNil(t, record)
			assert.Equal(t, record.Token, authenication2.RefreshToken)
			assert.Equal(t, record.UserId, user.Id)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, authenication1.AccessToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, authenication2.AccessToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()
	})))
	t.Run("WrongEmail", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate("some_wrong_prefix"+user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_WRONG_PASSWORD_OR_EMAIL+"\"", body)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByUserId(tx, ctx, user.Id)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()
	})))
	t.Run("WrongPassword", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, "some_wrong_prefix"+user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_WRONG_PASSWORD_OR_EMAIL+"\"", body)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByUserId(tx, ctx, user.Id)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()
	})))
}

func TestApiAuthRefresh(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var authenication1 auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &authenication1)

		assert.Nil(t, err)

		time.Sleep(1 * time.Second) // tokens generated based on time.Now(), sometimes we have equal values

		httpStatusCode, body, err = testHttpClient.RefreshToken(authenication1.RefreshToken)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var authenication2 auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &authenication2)

		assert.Nil(t, err)

		assert.NotEqual(t, authenication1.AccessToken, authenication2.AccessToken)
		assert.NotEqual(t, authenication1.RefreshToken, authenication2.RefreshToken)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, authenication1.RefreshToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			record, err := queries.GetRefreshTokenByToken(tx, ctx, authenication2.RefreshToken)

			assert.NotNil(t, record)
			assert.Equal(t, record.Token, authenication2.RefreshToken)
			assert.Equal(t, record.UserId, user.Id)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, authenication1.AccessToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, authenication2.AccessToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()

	})))
	t.Run("ExpiredRefreshToken", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var authenication1 auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &authenication1)

		assert.Nil(t, err)

		time.Sleep(10 * time.Second) // expected that .env.test has refresh token TTL in 10 seconds

		httpStatusCode, body, err = testHttpClient.RefreshToken(authenication1.RefreshToken)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_TOKEN_IS_EXPIRED+"\"", body)
	})))
}

func TestApiAuthAccess(t *testing.T) {
	t.Run("ValidAccessToken", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var authenication1 auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &authenication1)

		assert.Nil(t, err)

		time.Sleep(1 * time.Second) // expected that .env.test has access token TTL in 10 seconds

		httpStatusCode, body, err = testHttpClient.SafePing(authenication1.AccessToken)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\"Pong!\"", body)
	})))
	t.Run("ExpiredAccessToken", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var authenication1 auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &authenication1)

		assert.Nil(t, err)

		time.Sleep(10 * time.Second) // expected that .env.test has access token TTL in 10 seconds

		httpStatusCode, body, err = testHttpClient.SafePing(authenication1.AccessToken)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, httpStatusCode)
	})))
}
