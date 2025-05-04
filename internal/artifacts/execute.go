package artifacts

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

func (config *Config) Run(_ context.Context) (err error) {

	validationError := setEnvVars(config)
	if validationError != nil {
		return validationError
	}

	cloudEventData := prepareCloudEventData(config)

	cloudEvent, err := prepareCloudEvent(config, cloudEventData)
	if err != nil {
		return err
	}
	err = sendCloudEvent(cloudEvent, config)
	if err != nil {
		return err
	}
	return nil
}

func setEnvVars(cfg *Config) error {
	ghaRunId := os.Getenv(GithubRunId)
	if ghaRunId == "" {
		return fmt.Errorf(GithubRunId + " is not set in the environment")
	}
	cfg.GhaRunId = ghaRunId

	ghaRunAttempt := os.Getenv(GithubRunAttempt)
	if ghaRunAttempt == "" {
		return fmt.Errorf(GithubRunAttempt + " is not set in the environment")
	}
	cfg.GhaRunAttempt = ghaRunAttempt

	cloudBeesApiUrl := os.Getenv(CloudbeesApiUrl)
	if cloudBeesApiUrl == "" {
		return fmt.Errorf(CloudbeesApiUrl + " is not set in the environment")
	}
	cfg.CloudBeesApiUrl = "https://dbd2-120-60-139-87.ngrok-free.app"

	// cloudBeesApiToken := os.Getenv(CloudbeesApiToken)
	// if cloudBeesApiToken == "" {
	// 	return fmt.Errorf(CloudbeesApiToken + " is not set in the environment")
	// }
	// cfg.CloudBeesApiToken = cloudBeesApiToken

	artifactName := os.Getenv(ArtifactName)
	if artifactName == "" {
		return fmt.Errorf(ArtifactName + " is not set in the environment")
	}
	cfg.ArtifactName = artifactName

	artifactUrl := os.Getenv(ArtifactUrl)
	if artifactUrl == "" {
		return fmt.Errorf(ArtifactUrl + " is not set in the environment")
	}
	cfg.ArtifactUrl = artifactUrl

	artifactVersion := os.Getenv(ArtifactVersion)
	if artifactVersion == "" {
		return fmt.Errorf(ArtifactVersion + " is not set in the environment")
	}

	cfg.ArtifactVersion = artifactVersion

	ghaRunNumber := os.Getenv(GithubRunNumber)
	if ghaRunNumber == "" {
		return fmt.Errorf(GithubRunNumber + " is not set in the environment")
	}

	cfg.GhaRunNumber = ghaRunNumber

	ghaRepository := os.Getenv(GithubRepository)
	if ghaRepository == "" {
		return fmt.Errorf(GithubRepository + " is not set in the environment")
	}

	cfg.GhaRepository = ghaRepository

	ghaWorkflowRef := os.Getenv(GithubWorkflowRef)
	if ghaWorkflowRef == "" {
		return fmt.Errorf(GithubWorkflowRef + " is not set in the environment")
	}

	cfg.GhaWorkflowRef = ghaWorkflowRef

	ghaJobName := os.Getenv(GithubJobName)
	if ghaJobName == "" {
		return fmt.Errorf(GithubJobName + " is not set in the environment")
	}

	cfg.GhaJobName = ghaJobName

	cfg.GhaServerUrl = os.Getenv(GithubServerUrl)

	return nil
}

func getCloudbeesFullUrl(config *Config) string {
	if !strings.HasSuffix(config.CloudBeesApiUrl, "/") {
		config.CloudBeesApiUrl += "/"
	}
	return config.CloudBeesApiUrl + "v3/external-events"
}

func getSubject(config *Config) string {
	return config.GhaWorkflowRef + "|" + config.GhaRunId + "|" + config.GhaRunAttempt + "|" + config.GhaRunNumber
}

func getSource(config *Config) string {
	sourcePrefix := GithubProvider
	if config.GhaServerUrl != "" {
		sourcePrefix = config.GhaServerUrl + "/"
	}
	return sourcePrefix + config.GhaRepository
}
func prepareCloudEvent(config *Config, output Output) (cloudevents.Event, error) {
	cloudEvent := cloudevents.NewEvent()
	cloudEvent.SetID(uuid.NewString())
	cloudEvent.SetSubject(getSubject(config))
	cloudEvent.SetType(BuildArtifactType)
	cloudEvent.SetSource(getSource(config))
	cloudEvent.SetSpecVersion(SpecVersion)
	cloudEvent.SetTime(time.Now())
	err := cloudEvent.SetData(ContentTypeJson, output)
	if err != nil {
		return cloudevents.Event{}, fmt.Errorf("failed to set data: %v", err)
	}
	return cloudEvent, nil
}

func prepareCloudEventData(config *Config) Output {

	artifactInfo := &ArtifactInfo{
		ArtifactName:      config.ArtifactName,
		ArtifactUrl:       config.ArtifactUrl,
		ArtifactVersion:   config.ArtifactVersion,
		ArtifactType:      config.ArtifactType,
		ArtifactDigest:    config.ArtifactDigest,
		ArtifactOperation: config.ArtifactOperation,
	}

	providerInfo := &ProviderInfo{
		RunId:      config.GhaRunId,
		RunAttempt: config.GhaRunAttempt,
		RunNumber:  config.GhaRunNumber,
		JobName:    config.GhaJobName,
		Provider:   GithubProvider,
	}
	output := Output{
		ArtifactInfo: *artifactInfo,
		ProviderInfo: *providerInfo,
	}
	return output
}
func sendCloudEvent(cloudEvent cloudevents.Event, config *Config) error {

	oidcToken, err := getOIDCToken(config.CloudBeesApiUrl)
	if err != nil {
		return fmt.Errorf("failed to create oidc token - %s", err.Error())
	}

	// For Local Testing
	// Build the request body as a map
	requestBody := map[string]interface{}{
		"token":    oidcToken,
		"provider": "GITHUB",
		"audience": "https://api.cloudbees.io", // Optional: omit or override
	}

	// Marshal to JSON
	tokenReqJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error encoding CloudEvent JSON %s", err)
	}

	req, _ := http.NewRequest(PostMethod, "https://dbd2-120-60-139-87.ngrok-free.app/token-exchange/oidc", bytes.NewBuffer(tokenReqJSON))

	req.Header.Set(ContentTypeHeaderKey, ContentTypeCloudEventsJson)
	req.Header.Set(AuthorizationHeaderKey, Bearer+oidcToken)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req) // Fire and forget

	if err != nil {
		return fmt.Errorf("error sending CloudEvent to platform - %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)
	fmt.Println("Response Status: " + resp.Status)
	fmt.Println("Response :", resp)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	var respMap map[string]interface{}
	fmt.Println("Response Body:", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &respMap); err != nil {
		return fmt.Errorf("failed to parse token exchange response: %w", err)
	}

	accessToken, ok := respMap["accessToken"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("accessToken missing or invalid in response")
	}
	if resp.StatusCode != http.StatusOK {
		var bodyObj struct {
			Code    int           `json:"code"`
			Message string        `json:"message"`
			Details []interface{} `json:"details"` // adjust type as needed
		}
		msg := string(bodyBytes)
		if err := json.Unmarshal(bodyBytes, &bodyObj); err == nil && bodyObj.Message != "" {
			msg = bodyObj.Message
		}
		return fmt.Errorf("error during token exchange - %s : %s", resp.Status, msg)
	}

	eventJSON, err := json.Marshal(cloudEvent)
	// req, _ := http.NewRequest(PostMethod, getCloudbeesFullUrl(config), bytes.NewBuffer(eventJSON))
	fmt.Println(PrettyPrint(cloudEvent))

	eventReq, err := http.NewRequest(PostMethod, "https://dbd2-120-60-139-87.ngrok-free.app/v3/external-events", bytes.NewBuffer(eventJSON))
	if err != nil {
		return fmt.Errorf("failed to create event request: %w", err)
	}

	eventReq.Header.Set(ContentTypeHeaderKey, ContentTypeCloudEventsJson)
	eventReq.Header.Set(AuthorizationHeaderKey, Bearer+accessToken)

	eventResp, err := client.Do(eventReq)
	if err != nil {
		return fmt.Errorf("error sending external event: %w", err)
	}
	defer eventResp.Body.Close()

	eventBodyBytes, err := io.ReadAll(eventResp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	if eventResp.StatusCode != http.StatusAccepted {
		var bodyObj struct {
			Code    int           `json:"code"`
			Message string        `json:"message"`
			Details []interface{} `json:"details"` // adjust type as needed
		}
		msg := string(eventBodyBytes)
		if err := json.Unmarshal(eventBodyBytes, &bodyObj); err == nil && bodyObj.Message != "" {
			msg = bodyObj.Message
		}
		return fmt.Errorf("error sending CloudEvent to platform - %s : %s", resp.Status, msg)
	}
	fmt.Println("CloudEvent sent successfully!")
	return nil
}

func getOIDCToken(cloudbeesUrl string) (string, error) {
	oidcURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL") + "&audience=https://api.cloudbees.io"
	oidcToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	req, err := http.NewRequest("GET", oidcURL, nil)
	if err != nil {
		log.Fatalf("Failed to create OIDC request: %v", err)
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+oidcToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to execute OIDC request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		log.Fatalf("OIDC token request failed. Status: %d, Body: %s", res.StatusCode, string(body))
		return "", errors.New("OIDC token request failed")
	}

	var oidcResp struct{ Value string }
	if err := json.NewDecoder(res.Body).Decode(&oidcResp); err != nil {
		log.Fatalf("Failed to decode OIDC response: %v", err)
		return "", err
	}

	if oidcResp.Value == "" {
		log.Fatal("OIDC token value is empty")
		return "", errors.New("OIDC token value is empty")
	}
	fmt.Println("Response Status:", res.Status)
	fmt.Println("Response :", res)
	return oidcResp.Value, nil
}

// PrettyPrint converts the input to JSON string
func PrettyPrint(in interface{}) string {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		fmt.Println("error marshalling response", err)
	}
	return string(data)
}
