package dtrack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Since v5.0

type ExtensionService struct {
	client *Client
}

type ListExtensionPointsResponseItem struct {
	Name string `json:"name"`
}

type ListExtensionPointsResponse struct {
	Items         []ListExtensionPointsResponseItem `json:"items"`
	Total         TotalCount                        `json:"total"`
	NextPageToken *string                           `json:"next_page_token,omitempty"`
}

type ListExtensionsResponseItem struct {
	Configurable bool    `json:"configurable"`
	DisplayName  *string `json:"display_name,omitempty"` // Since v5.1
	Name         string  `json:"name"`
	Testable     bool    `json:"testable"`
}

type ListExtensionsResponse struct {
	Items         []ListExtensionsResponseItem `json:"items"`
	Total         TotalCount                   `json:"total"`
	NextPageToken *string                      `json:"next_page_token,omitempty"`
}

type ExtensionIdentity struct {
	ExtensionPointName string `json:"extension_point_name"`
	ExtensionName      string `json:"extension_name"`
}

type ExtensionConfig struct {
	Config map[string]any `json:"config"`
}

// Custom UnmarshalJSON
type GetExtensionConfigSchemaResponse struct {
	SchemaURL string
	Schema    map[string]any
}

type TestExtensionResponse struct {
	Checks []ExtensionTestCheck `json:"checks"`
}

type ExtensionTestCheck struct {
	Name    string                   `json:"name"`
	Status  ExtensionTestCheckStatus `json:"status"`
	Message *string                  `json:"message,omitempty"`
}

type ExtensionTestCheckStatus string

const (
	ExtensionTestCheckStatusPassed  ExtensionTestCheckStatus = "PASSED"
	ExtensionTestCheckStatusFailed  ExtensionTestCheckStatus = "FAILED"
	ExtensionTestCheckStatusSkipped ExtensionTestCheckStatus = "SKIPPED"
)

func (es ExtensionService) ListExtensionPoints(ctx context.Context) (r ListExtensionPointsResponse, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodGet, "api/v2/extension-points")
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &r)
	return
}

func (es ExtensionService) ListExtensions(ctx context.Context, extensionPointName string) (r ListExtensionsResponse, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v2/extension-points/%s/extensions", extensionPointName))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &r)
	return
}

func (es ExtensionService) GetExtensionConfig(ctx context.Context, extension ExtensionIdentity) (cfg ExtensionConfig, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v2/extension-points/%s/extensions/%s/config", extension.ExtensionPointName, extension.ExtensionName))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &cfg)
	return
}

func (es ExtensionService) GetExtensionConfigSchema(ctx context.Context, extension ExtensionIdentity) (schema GetExtensionConfigSchemaResponse, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v2/extension-points/%s/extensions/%s/config-schema", extension.ExtensionPointName, extension.ExtensionName))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &schema)
	if err != nil {
		return
	}

	return
}

func (es ExtensionService) UpdateExtensionConfig(ctx context.Context, extension ExtensionIdentity, config ExtensionConfig) (err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("api/v2/extension-points/%s/extensions/%s/config", extension.ExtensionPointName, extension.ExtensionName), withBody(config), withAcceptContentType(""))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, nil)
	return
}

func (es ExtensionService) TestExtension(ctx context.Context, extension ExtensionIdentity, config ExtensionConfig) (r TestExtensionResponse, err error) {
	err = es.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := es.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("api/v2/extension-points/%s/extensions/%s/test", extension.ExtensionPointName, extension.ExtensionName), withBody(config))
	if err != nil {
		return
	}

	_, err = es.client.doRequest(req, &r)
	return
}

func (r *GetExtensionConfigSchemaResponse) UnmarshalJSON(data []byte) error {
	var response map[string]any
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	r.SchemaURL, _ = response["$schema"].(string)
	delete(response, "$schema")
	r.Schema = response
	return nil
}
