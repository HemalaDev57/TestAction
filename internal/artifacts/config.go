package artifacts

import "context"

type Config struct {
	context.Context
	ArtifactName    string `json:"artifact-name,omitempty"`
	ArtifactUrl     string `json:"artifact-url,omitempty"`
	ArtifactVersion string `json:"artifact-version,omitempty"`
	ArtifactType    string `json:"artifact-type,omitempty"`
	ArtifactDigest  string `json:"artifact-digest,omitempty"`
	ArtifactLabel   string `json:"artifact-label,omitempty"`
	GhaRunId        string `json:"gha-run-id,omitempty"`
	GhaRunAttempt   string `json:"gha-run-attempt,omitempty"`
	GhaRunNumber    string `json:"gha-run-number,omitempty"`
	CloudBeesApiUrl string `json:"cloudbees-api-url,omitempty"`
	GhaRepository   string `json:"gha-repository,omitempty"`
	GhaWorkflowRef  string `json:"gha-workflow-ref,omitempty"`
	GhaServerUrl    string `json:"gha-server-url,omitempty"`
	GhaJobName      string `json:"gha-job-name,omitempty"`
}
