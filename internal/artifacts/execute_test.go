package artifacts

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	t.Run("Missing Env:"+GithubRunId, func(t *testing.T) {
		var config = Config{}

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), GithubRunId+" is not set in the environment")
	})

	t.Run("Missing Env:"+GithubRunAttempt, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), GithubRunAttempt+" is not set in the environment")
	})

	t.Run("Missing Env:"+CloudbeesApiUrl, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), CloudbeesApiUrl+" is not set in the environment")
	})

	t.Run("Missing Env:"+ArtifactName, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(CloudbeesApiUrl, "https://api-test.cloudbees.com")
		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), ArtifactName+" is not set in the environment")
	})

	t.Run("Missing Env:"+ArtifactUrl, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(CloudbeesApiUrl, "https://api-test.cloudbees.com")
		os.Setenv(ArtifactName, "testartifact")
		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), ArtifactUrl+" is not set in the environment")
	})

	t.Run("Missing Env:"+ArtifactVersion, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(CloudbeesApiUrl, "https://api-test.cloudbees.com")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), ArtifactVersion+" is not set in the environment")
	})

	t.Run("Missing Env:"+GithubRunNumber, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(CloudbeesApiUrl, "https://api-test.cloudbees.com")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), GithubRunNumber+" is not set in the environment")
	})

	t.Run("Missing Env:"+GithubRepository, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(CloudbeesApiUrl, "https://api-test.cloudbees.com")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), GithubRepository+" is not set in the environment")
	})

	t.Run("Missing Env:"+GithubWorkflowRef, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(CloudbeesApiUrl, "https://api-test.cloudbees.com")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "test/test")
		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), GithubWorkflowRef+" is not set in the environment")
	})

	t.Run("Missing Env:"+GithubJobName, func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(CloudbeesApiUrl, "https://api-test.cloudbees.com")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "test/test")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), GithubJobName+" is not set in the environment")
	})
	t.Run("Success", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "HemalaDev57/TestAction")
		os.Setenv(GithubWorkflowRef, "HemalaDev57/TestAction/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		// Set the httpClient to use the test server
		os.Setenv(CloudbeesApiUrl, ts.URL)

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.Nil(t, err)
		//assert.Equal(t, err.Error(), GithubWorkflowRef+" is not set in the environment")
	})

	t.Run("Success All Fields", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "HemalaDev57/TestAction")
		os.Setenv(GithubWorkflowRef, "HemalaDev57/TestAction/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubServerUrl, "https://github.com")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		// Set the httpClient to use the test server
		os.Setenv(CloudbeesApiUrl, ts.URL)

		err := config.Run(context.Background())
		assert.Nil(t, err)
		//assert.Equal(t, err.Error(), GithubWorkflowRef+" is not set in the environment")
	})

	t.Run("Failed Sending cloudevent", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "HemalaDev57/TestAction")
		os.Setenv(GithubWorkflowRef, "HemalaDev57/TestAction/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer ts.Close()

		// Set the httpClient to use the test server
		os.Setenv(CloudbeesApiUrl, ts.URL)

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "error sending CloudEvent to platform - 502 Bad Gateway : ")
	})

	t.Run("Error sending CloudEvent", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "HemalaDev57/TestAction")
		os.Setenv(GithubWorkflowRef, "HemalaDev57/TestAction/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Set the httpClient to use the test server
		os.Setenv(CloudbeesApiUrl, "testurl")

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
	})
}
