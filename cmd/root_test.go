package cmd

import (
	"fmt"
	"gha-register-build-artifact/internal/artifacts"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetDefaultValuesEmpty(t *testing.T) {
	var config = artifacts.Config{}
	setDefaultValues(&config)
	assert.Equal(t, config.ArtifactType, "")
	assert.Equal(t, config.ArtifactDigest, "")
}

func Test_SetDefaultValues(t *testing.T) {
	var config = artifacts.Config{}
	os.Setenv(artifacts.ArtifactType, "123456789")
	os.Setenv(artifacts.ArtifactDigest, "123456789:digest")
	setDefaultValues(&config)
	assert.Equal(t, config.ArtifactType, "123456789")
	assert.Equal(t, config.ArtifactDigest, "123456789:digest")
}

func Test_Run(t *testing.T) {
	os.Setenv(artifacts.GithubRunId, "123456789")
	os.Setenv(artifacts.GithubRunAttempt, "1")
	os.Setenv(artifacts.ArtifactName, "testartifact")
	os.Setenv(artifacts.ArtifactUrl, "https://test.com")
	os.Setenv(artifacts.ArtifactVersion, "1.0.0")
	os.Setenv(artifacts.GithubRunNumber, "123")
	os.Setenv(artifacts.GithubRepository, "SrimanPadmanabanCB/gha-action")
	os.Setenv(artifacts.GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
	os.Setenv(artifacts.GithubServerUrl, "https://github.com")
	os.Setenv(artifacts.GithubJobName, "testjob")
	os.Setenv(artifacts.ArtifactLabel, "labelA,labelB,labelC")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Mock Request URL:", r.URL.String())
		switch {
		case r.Method == "GET" && strings.HasPrefix(r.URL.String(), "/?audience="):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"value": "mock-oidc-token"}`))
		case r.Method == "POST" && r.URL.Path == "/token-exchange/external-oidc-id-token":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"accessToken": "mock-cbp-token"}`))
		case r.Method == "POST" && r.URL.Path == "/v3/external-events":
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "unexpected request: "+r.URL.Path, http.StatusNotFound)
		}
	}))
	defer ts.Close()

	os.Setenv(artifacts.CloudbeesApiUrl, ts.URL)
	os.Setenv(artifacts.ActionIdTokenRequestUrl, ts.URL)
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
	os.Setenv(artifacts.GithubRepository, "SrimanPadmanabanCB/gha-action")
	os.Setenv(artifacts.GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
	os.Setenv(artifacts.GithubServerUrl, "https://github.com")
	os.Setenv(artifacts.GithubJobName, "testjob")
	os.Setenv(artifacts.ArtifactLabel, "labelA,labelB,labelC")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Mock Request URL:", r.URL.String())
		switch {
		case r.Method == "GET" && strings.HasPrefix(r.URL.String(), "/?audience="):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"value": "mock-oidc-token"}`))
		case r.Method == "POST" && r.URL.Path == "/token-exchange/external-oidc-id-token":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"accessToken": "mock-cbp-token"}`))
		case r.Method == "POST" && r.URL.Path == "/v3/external-events":
			w.WriteHeader(http.StatusBadGateway)
		default:
			http.Error(w, "unexpected request: "+r.URL.Path, http.StatusNotFound)
		}
	}))
	defer ts.Close()

	os.Setenv(artifacts.CloudbeesApiUrl, ts.URL)
	os.Setenv(artifacts.ActionIdTokenRequestUrl, ts.URL)
	err := run(nil, nil)
	assert.Contains(t, err.Error(), "error sending CloudEvent to platform")
}

func Test_Failure_1(t *testing.T) {
	err := run(nil, []string{"test, command"})
	assert.Contains(t, err.Error(), "unknown arguments:")
}
