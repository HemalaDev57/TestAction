package artifacts

type ArtifactInfo struct {
	ArtifactName      string `json:"artifact_name,omitempty"`
	ArtifactUrl       string `json:"artifact_url,omitempty"`
	ArtifactVersion   string `json:"artifact_version,omitempty"`
	ArtifactType      string `json:"artifact_type,omitempty"`
	ArtifactDigest    string `json:"artifact_digest,omitempty"`
	ArtifactOperation string `json:"artifact_operation,omitempty"`
	ArtifactLabel     string `json:"artifact_label,omitempty"`
}

type ProviderInfo struct {
	RunId      string `json:"run_id,omitempty"`
	RunAttempt string `json:"run_attempt,omitempty"`
	RunNumber  string `json:"run_number,omitempty"`
	JobName    string `json:"job_name,omitempty"`
	Provider   string `json:"provider,omitempty"`
}

type Output struct {
	ProviderInfo ProviderInfo `json:"provider_info"`
	ArtifactInfo ArtifactInfo `json:"artifact_info"`
}
