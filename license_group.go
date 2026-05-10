package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type LicenseGroup struct {
	UUID       uuid.UUID `json:"uuid"`
	Name       string    `json:"name"`
	Licenses   []License `json:"licenses,omitempty"`
	RiskWeight int32     `json:"riskWeight,omitempty"`
}

type LicenseGroupService struct {
	client *Client
}

func (ls LicenseGroupService) AddLicense(ctx context.Context, groupUUID, licenseUUID uuid.UUID) (lg LicenseGroup, err error) {
	err = ls.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("api/v1/licenseGroup/%s/license/%s", groupUUID, licenseUUID))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &lg)
	return
}

func (ls LicenseGroupService) RemoveLicense(ctx context.Context, groupUUID, licenseUUID uuid.UUID) (lg LicenseGroup, err error) {
	err = ls.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/licenseGroup/%s/license/%s", groupUUID, licenseUUID))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &lg)
	return
}

func (ls LicenseGroupService) GetAll(ctx context.Context, po PageOptions, so SortOptions) (p Page[LicenseGroup], err error) {
	err = ls.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodGet, "api/v1/licenseGroup")
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

func (ls LicenseGroupService) Create(ctx context.Context, group LicenseGroup) (lg LicenseGroup, err error) {
	err = ls.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodPut, "api/v1/licenseGroup", withBody(group))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &lg)
	return
}

func (ls LicenseGroupService) Update(ctx context.Context, group LicenseGroup) (lg LicenseGroup, err error) {
	err = ls.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodPost, "api/v1/licenseGroup", withBody(group))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &lg)
	return
}

func (ls LicenseGroupService) Get(ctx context.Context, groupUUID uuid.UUID) (lg LicenseGroup, err error) {
	err = ls.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v1/licenseGroup/%s", groupUUID))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, &lg)
	return
}

func (ls LicenseGroupService) Delete(ctx context.Context, groupUUID uuid.UUID) (err error) {
	err = ls.client.assertServerVersionAtLeast("4.0.0")
	if err != nil {
		return
	}

	req, err := ls.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/licenseGroup/%s", groupUUID))
	if err != nil {
		return
	}

	_, err = ls.client.doRequest(req, nil)
	return
}
