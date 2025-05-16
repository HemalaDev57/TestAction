package artifacts

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

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

		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL
		err := config.Run(context.Background())
		assert.Nil(t, err)
	})

	t.Run("Success All Fields", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubServerUrl, "https://github.com")
		os.Setenv(GithubJobName, "testjob")

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

		// Set the httpClient to use the test server
		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL

		err := config.Run(context.Background())
		assert.Nil(t, err)
		//assert.Equal(t, err.Error(), GithubWorkflowRef+" is not set in the environment")
	})

	t.Run("Failed OIDC token request", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer ts.Close()

		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL

		err := config.Run(context.Background())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "OIDC token request failed")
	})

	t.Run("Failed OIDC token request - Empty token", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"value": ""}`))
		}))
		defer ts.Close()

		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "OIDC token value is empty")
	})

	t.Run("Failed OIDC token request - Invalid token", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"value": 1}`))
		}))
		defer ts.Close()

		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL

		err := config.Run(context.Background())
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to create oidc token")
	})

	t.Run("Error OIDC token request", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "error", http.StatusBadGateway)
		}))
		defer ts.Close()

		os.Setenv(CloudbeesApiUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL
		err := config.Run(context.Background())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to create oidc token")
	})

	t.Run("Failed token exchange - Access Token missing", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"value": "mock-oidc-token"}`))
		}))
		defer ts.Close()

		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL

		err := config.Run(context.Background())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "accessToken missing or invalid in response")
	})

	t.Run("Error token exchange", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"value": "mock-oidc-token"}`))
		}))
		defer ts.Close()

		os.Setenv(CloudbeesApiUrl, "test")
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		err := config.Run(context.Background())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error sending CloudEvent to platform")
	})

	t.Run("Failed token exchange - Non 200", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Mock Request URL:", r.URL.String())
			switch {
			case r.Method == "GET" && strings.HasPrefix(r.URL.String(), "/?audience="):
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"value": "mock-oidc-token"}`))
			case r.Method == "POST" && r.URL.Path == "/token-exchange/external-oidc-id-token":
				w.WriteHeader(http.StatusBadGateway)
			default:
				http.Error(w, "unexpected request: "+r.URL.Path, http.StatusNotFound)
			}
		}))
		defer ts.Close()

		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL

		err := config.Run(context.Background())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error during token exchange - 502 Bad Gateway :")
	})
	t.Run("Error sending CloudEvent", func(t *testing.T) {
		var config = Config{}
		os.Setenv(GithubRunId, "123456789")
		os.Setenv(GithubRunAttempt, "1")
		os.Setenv(ArtifactName, "testartifact")
		os.Setenv(ArtifactUrl, "https://test.com")
		os.Setenv(ArtifactVersion, "1.0.0")
		os.Setenv(GithubRunNumber, "123")
		os.Setenv(GithubRepository, "SrimanPadmanabanCB/gha-action")
		os.Setenv(GithubWorkflowRef, "SrimanPadmanabanCB/gha-action/.github/workflows/test_action.yml@refs/heads/main")
		os.Setenv(GithubJobName, "testjob")

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
				http.Error(w, "test error", http.StatusBadGateway)
			default:
				http.Error(w, "unexpected request: "+r.URL.Path, http.StatusNotFound)
			}
		}))
		defer ts.Close()

		// Set the httpClient to use the test server
		os.Setenv(CloudbeesApiUrl, ts.URL)
		os.Setenv(ActionIdTokenRequestUrl, ts.URL)
		config.CloudBeesApiUrl = ts.URL

		err := config.Run(context.Background())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error sending CloudEvent to platform - 502 Bad Gateway :")
	})
}
