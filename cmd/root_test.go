package cmd

import (
	"gha-register-build-artifact/internal/artifacts"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetDefaultValuesEmpty(t *testing.T) {
	var config = artifacts.Config{}
	setDefaultValues(&config)
	assert.Equal(t, config.ArtifactType, "")
	assert.Equal(t, config.ArtifactDigest, "")
	assert.Equal(t, config.ArtifactOperation, artifacts.PUBLISHED)
}

func Test_SetDefaultValues(t *testing.T) {
	var config = artifacts.Config{}
	os.Setenv(artifacts.ArtifactType, "123456789")
	os.Setenv(artifacts.ArtifactDigest, "123456789:digest")
	setDefaultValues(&config)
	assert.Equal(t, config.ArtifactType, "123456789")
	assert.Equal(t, config.ArtifactDigest, "123456789:digest")
	assert.Equal(t, config.ArtifactOperation, artifacts.PUBLISHED)
}

func Test_Run(t *testing.T) {
	os.Setenv(artifacts.GithubRunId, "123456789")
	os.Setenv(artifacts.GithubRunAttempt, "1")
	os.Setenv(artifacts.ArtifactName, "testartifact")
	os.Setenv(artifacts.ArtifactUrl, "https://test.com")
	os.Setenv(artifacts.ArtifactVersion, "1.0.0")
	os.Setenv(artifacts.GithubRunNumber, "123")
	os.Setenv(artifacts.GithubRepository, "HemalaDev57/TestAction")
	os.Setenv(artifacts.GithubWorkflowRef, "HemalaDev57/TestAction/.github/workflows/test_action.yml@refs/heads/main")
	os.Setenv(artifacts.GithubServerUrl, "https://github.com")
	os.Setenv(artifacts.GithubJobName, "testjob")

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Set the httpClient to use the test server
	os.Setenv(artifacts.CloudbeesApiUrl, ts.URL)
	err := run(nil, nil)
	assert.Nil(t, err)
}

func Test_Failure(t *testing.T) {
	os.Setenv(artifacts.GithubRunId, "123456789")
	os.Setenv(artifacts.GithubRunAttempt, "1")
	os.Setenv(artifacts.ArtifactName, "testartifact")
	os.Setenv(artifacts.ArtifactUrl, "https://test.com")
	os.Setenv(artifacts.ArtifactVersion, "1.0.0")
	os.Setenv(artifacts.GithubRunNumber, "123")
	os.Setenv(artifacts.GithubRepository, "HemalaDev57/TestAction")
	os.Setenv(artifacts.GithubWorkflowRef, "HemalaDev57/TestAction/.github/workflows/test_action.yml@refs/heads/main")
	os.Setenv(artifacts.GithubServerUrl, "https://github.com")
	os.Setenv(artifacts.GithubJobName, "testjob")

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Set the httpClient to use the test server
	os.Setenv(artifacts.CloudbeesApiUrl, ts.URL)
	err := run(nil, nil)
	assert.Nil(t, err)
}
