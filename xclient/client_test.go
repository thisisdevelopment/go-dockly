package xclient_test

import (
	"net/http"

	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func (s *TestSuite) Test_Client() {

	expected := map[string]string{"foo": "bar"}
	expectedPath := "/bar"
	expectedStatus := http.StatusOK

	gock.New(s.baseURL).
		Get(expectedPath).
		Reply(expectedStatus).
		JSON(expected)

	result := map[string]string{}

	actualStatus, err := s.client.Do("GET", expectedPath, nil, &result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), expected, result)
	require.Equal(s.T(), actualStatus, expectedStatus)

	// Verify that we don't have pending mocks
	require.Equal(s.T(), gock.IsDone(), true)
}
