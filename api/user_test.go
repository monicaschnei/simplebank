package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_sqlc "github.com/monicaschnei/simplebank/db/mock"
	db "github.com/monicaschnei/simplebank/db/sqlc"
	util "github.com/monicaschnei/simplebank/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_sqlc.NewMockStore(ctrl)
	arg := db.CreateUserParams{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}

	store.EXPECT().
		CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
		Times(1).
		Return(user, nil)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	body := gin.H{
		"username":  user.Username,
		"password":  password,
		"full_name": user.FullName,
		"email":     user.Email,
	}
	//Marshal body data to JSON
	data, err := json.Marshal(body)
	require.NoError(t, err)

	url := "/users"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchUser(t, recorder.Body, user)

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandonString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
