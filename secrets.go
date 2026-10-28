package dtrack

import (
	"context"
	"net/http"
	"strconv"
)

type SecretsService struct {
	client *Client
}

type ListSecretMetadataQuery struct {
	Filter    string
	PageToken string
	Limit     uint32
}

type SecretMetadata struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	CreatedAt   *uint64 `json:"created_at,omitempty"`
	UpdatedAt   *uint64 `json:"updated_at,omitempty"`
}

type ListSecretsMetadataResponse struct {
	Items         []SecretMetadata `json:"items"`
	Total         TotalCount       `json:"total"`
	NextPageToken string           `json:"next_page_token,omitempty"`
}

type CreateSecretRequest struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

type UpdateSecretRequest struct {
	Description *string `json:"description,omitempty"`
	Value       *string `json:"value,omitempty"`
}

func (ss SecretsService) ListMetadata(ctx context.Context, query ListSecretMetadataQuery) (r ListSecretsMetadataResponse, err error) {
	err = ss.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	params := map[string]string{}
	if len(query.Filter) > 0 {
		params["q"] = query.Filter
	}
	if len(query.PageToken) > 0 {
		params["page_token"] = query.PageToken
	}
	if query.Limit > 0 {
		params["limit"] = strconv.Itoa(int(query.Limit))
	}
	req, err := ss.client.newRequest(ctx, http.MethodGet, "api/v2/secrets", withParams(params))
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, &r)
	return
}

func (ss SecretsService) Create(ctx context.Context, secret CreateSecretRequest) (url string, err error) {
	err = ss.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := ss.client.newRequest(ctx, http.MethodPost, "api/v2/secrets", withBody(secret))
	if err != nil {
		return
	}

	res, err := ss.client.doRequest(req, nil)
	if err != nil {
		return
	}

	url = res.Header.Get("Location")
	return
}

func (ss SecretsService) GetMetadata(ctx context.Context, name string) (md SecretMetadata, err error) {
	err = ss.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := ss.client.newRequest(ctx, http.MethodGet, "api/v2/secrets/"+name)
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, &md)
	return
}

func (ss SecretsService) Delete(ctx context.Context, name string) (err error) {
	err = ss.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := ss.client.newRequest(ctx, http.MethodDelete, "api/v2/secrets/"+name)
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, nil)
	return
}

func (ss SecretsService) Update(ctx context.Context, name string, secret UpdateSecretRequest) (err error) {
	err = ss.client.assertServerVersionAtLeast("5.0.0")
	if err != nil {
		return
	}

	req, err := ss.client.newRequest(ctx, http.MethodPatch, "api/v2/secrets/"+name, withBody(secret))
	if err != nil {
		return
	}

	_, err = ss.client.doRequest(req, nil)
	return
}
