package xclient_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"gopkg.in/h2non/gock.v1"
	"io"
	"net/http"
)

const (
	expectedPath   = "bar"
	expectedStatus = http.StatusOK
)

var (
	expected = map[string]string{"foo": "bar"}
	result   = map[string]string{}
)

func (s *TestSuite) Test_Client_Get() {

	gock.New(s.baseURL).
		Get(expectedPath).
		Reply(expectedStatus).
		JSON(expected)

	actualStatus, err := s.client.Do(context.Background(), "GET", expectedPath, nil, nil, &result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), expected, result)
	require.Equal(s.T(), actualStatus, expectedStatus)

	// Verify that we don't have pending mocks
	require.Equal(s.T(), gock.IsDone(), true)
}

func (s *TestSuite) Test_Client_Post() {
	gock.New(s.baseURL).
		Post(expectedPath).
		Reply(expectedStatus).
		JSON(expected)

	actualStatus, err := s.client.Do(context.Background(), "POST", expectedPath, expected, nil, &result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), expected, result)
	require.Equal(s.T(), actualStatus, expectedStatus)

	// Verify that we don't have pending mocks
	require.Equal(s.T(), gock.IsDone(), true)
}

func (s *TestSuite) Test_Client_Post_Reader() {

	gock.New(s.baseURL).
		Post(expectedPath).
		Reply(expectedStatus).
		JSON(expected)

	b, _ := json.Marshal(expected)

	actualStatus, err := s.client.Do(context.Background(), "POST", expectedPath, io.NopCloser(bytes.NewReader(b)), nil, &result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), expected, result)
	require.Equal(s.T(), actualStatus, expectedStatus)

	// Verify that we don't have pending mocks
	require.Equal(s.T(), gock.IsDone(), true)
}
