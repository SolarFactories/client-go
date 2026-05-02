package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type License struct {
	UUID                uuid.UUID `json:"uuid"`
	Name                string    `json:"name"`
	Text                string    `json:"text"`
	Template            string    `json:"template"`
	Header              string    `json:"header"`
	Comment             string    `json:"comment"`
	LicenseID           string    `json:"licenseId"`
	OSIApproved         bool      `json:"isOsiApproved"`
	FSFLibre            bool      `json:"isFsfLibre"`
	DeprecatedLicenseID bool      `json:"isDeprecatedLicenseId"`
	SeeAlso             []string  `json:"seeAlso"`
}

type LicenseService struct {
	client *Client
}

func (ls LicenseService) GetAll(ctx context.Context, po PageOptions) (p Page[License], err error) {
	err = ls.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodGet, "api/v1/license", withPageOptions(po))
	if err != nil {
		return
	}

	res, err := ls.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func (ls LicenseService) Get(ctx context.Context, licenseSPDX string) (l License, err error) {
	err = ls.client.assertServerVersionAtLeast("3.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v1/license/%s", licenseSPDX))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &l)
	return
}

func (ls LicenseService) Create(ctx context.Context, license License) (l License, err error) {
	err = ls.client.assertServerVersionAtLeast("4.7.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodPut, "api/v1/license", withBody(license))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &l)
	return
}

func (ls LicenseService) Delete(ctx context.Context, licenseSPDX string) (err error) {
	err = ls.client.assertServerVersionAtLeast("4.7.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/license/%s", licenseSPDX))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, nil)
	return
}

func (ls LicenseService) GetConcise(ctx context.Context) (licenses []License, err error) {
	err = ls.client.assertServerVersionAtLeast("3.4.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodGet, "api/v1/license/concise")
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &licenses)
	return
}
