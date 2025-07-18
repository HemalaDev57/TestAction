package artifacts

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}

type TokenRequest struct {
	Provider string `json:"provider"`
	Audience string `json:"audience"`
}

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
	cfg.CloudBeesApiUrl = cloudBeesApiUrl

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

	cfg.ArtifactLabel = os.Getenv(ArtifactLabel)

	return nil
}

func getExternalEventlUrl(config *Config) string {
	if !strings.HasSuffix(config.CloudBeesApiUrl, "/") {
		config.CloudBeesApiUrl += "/"
	}
	return config.CloudBeesApiUrl + "v3/external-events"
}

func getExternalTokenExchangeUrl(config *Config) string {
	if !strings.HasSuffix(config.CloudBeesApiUrl, "/") {
		config.CloudBeesApiUrl += "/"
	}
	return config.CloudBeesApiUrl + "token-exchange/external-oidc-id-token"
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
		ArtifactName:    config.ArtifactName,
		ArtifactUrl:     config.ArtifactUrl,
		ArtifactVersion: config.ArtifactVersion,
		ArtifactType:    config.ArtifactType,
		ArtifactDigest:  config.ArtifactDigest,
		ArtifactLabel:   config.ArtifactLabel,
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
	// Fetch the OIDC token
	// This token is used to authenticate the request to the CloudBees API
	fmt.Println("Started fetching OIDC Token...")
	oidcToken, err := getOIDCToken(config.CloudBeesApiUrl)
	if err != nil {
		return fmt.Errorf("failed to create oidc token - %s", err.Error())
	}
	fmt.Println("OIDC Token fetched successfully!")

	fmt.Println("Initiated exchanging the OIDC Token with CBP token...")
	tokenRequestObj := TokenRequest{
		Provider: GithubProvider,
		Audience: strings.TrimSuffix(config.CloudBeesApiUrl, "/"), // Optional: omit or override
	}
	tokenReqJSON, err := json.Marshal(tokenRequestObj)
	if err != nil {
		return fmt.Errorf("error encoding CloudEvent JSON %s", err)
	}

	tokenReq, _ := http.NewRequest(PostMethod, getExternalTokenExchangeUrl(config), bytes.NewBuffer(tokenReqJSON))
	tokenReq.Header.Set(ContentTypeHeaderKey, ContentTypeCloudEventsJson)
	tokenReq.Header.Set(AuthorizationHeaderKey, Bearer+oidcToken)

	client := &http.Client{}
	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return fmt.Errorf("error sending CloudEvent to platform - %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(tokenResp.Body)

	bodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	if tokenResp.StatusCode != http.StatusOK {
		bodyObj := ErrorResponse{}
		msg := string(bodyBytes)
		if err := json.Unmarshal(bodyBytes, &bodyObj); err == nil && bodyObj.Message != "" {
			msg = bodyObj.Message
		}
		return fmt.Errorf("error during token exchange - %s : %s", tokenResp.Status, msg)
	}

	var respMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &respMap); err != nil {
		return fmt.Errorf("failed to parse token exchange response: %w", err)
	}

	accessToken, ok := respMap[AccessToken].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("accessToken missing or invalid in response")
	}
	log.Println(base64.StdEncoding.EncodeToString([]byte(accessToken)))
	// Write the token to a file
	err = os.WriteFile("access_token.txt", []byte(accessToken), 0644)
	if err != nil {
		return fmt.Errorf("failed to write token to file: %w", err)
	}
	fmt.Println("Token exchange successful!")
	data, err := os.ReadFile("access_token.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	fmt.Println("Access token from file:", string(data))

	fmt.Println("Initiated sending the CloudEvent to platform...")
	eventJSON, err := json.Marshal(cloudEvent)
	fmt.Println(PrettyPrint(cloudEvent))

	eventReq, err := http.NewRequest(PostMethod, getExternalEventlUrl(config), bytes.NewBuffer(eventJSON))
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
	if eventResp.StatusCode != http.StatusOK {
		bodyObj := ErrorResponse{}
		msg := string(eventBodyBytes)
		if err := json.Unmarshal(eventBodyBytes, &bodyObj); err == nil && bodyObj.Message != "" {
			msg = bodyObj.Message
		}
		return fmt.Errorf("error sending CloudEvent to platform - %s : %s", eventResp.Status, msg)
	}
	fmt.Println("CloudEvent sent successfully!")
	return nil
}

func getOIDCToken(cloudbeesUrl string) (string, error) {
	log.Println(ActionIdTokenRequestToken)
	encoded := base64.StdEncoding.EncodeToString([]byte(os.Getenv(ActionIdTokenRequestToken)))
	log.Println(encoded)
	log.Println(ActionIdTokenRequestUrl)
	log.Println(os.Getenv(ActionIdTokenRequestUrl))
	oidcToken := os.Getenv(ActionIdTokenRequestToken)
	oidcBaseURL := os.Getenv(ActionIdTokenRequestUrl)
	oidcAudience := url.QueryEscape(strings.TrimSuffix(cloudbeesUrl, "/"))
	oidcURL := fmt.Sprintf("%s?audience=%s", oidcBaseURL, oidcAudience)

	oidcTokenReq, err := http.NewRequest("GET", oidcURL, nil)
	if err != nil {
		log.Printf("Failed to create OIDC request: %v", err)
		return "", err
	}
	oidcTokenReq.Header.Add(AuthorizationHeaderKey, Bearer+oidcToken)
	client := &http.Client{}
	oidcTokenResp, err := client.Do(oidcTokenReq)
	if err != nil {
		log.Printf("Failed to execute OIDC request: %v", err)
		return "", err
	}
	defer oidcTokenResp.Body.Close()

	if oidcTokenResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(oidcTokenResp.Body)
		log.Printf("OIDC token request failed. Status: %d, Body: %s", oidcTokenResp.StatusCode, string(body))
		return "", errors.New("OIDC token request failed")
	}

	var oidcResp struct{ Value string }
	if err := json.NewDecoder(oidcTokenResp.Body).Decode(&oidcResp); err != nil {
		log.Printf("Failed to decode OIDC response: %v", err)
		return "", err
	}
	if oidcResp.Value == "" {
		log.Printf("OIDC token value is empty")
		return "", errors.New("OIDC token value is empty")
	}
	log.Print("Response Value : ")
	log.Println(base64.StdEncoding.EncodeToString([]byte(oidcResp.Value)))
	return oidcResp.Value, nil
}

// PrettyPrint converts the input to JSON string
func PrettyPrint(in any) string {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		fmt.Println("error marshalling response", err)
	}
	return string(data)
}
