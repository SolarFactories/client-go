package dtrack

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"
)

func TestBOMService_Upload(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionBOMUpload,
			PermissionProjectCreationUpload,
			PermissionViewPortfolio,
		},
	})

	_, err := client.BOM.Upload(context.Background(), BOMUploadRequest{
		ProjectName:    "acme-app",
		ProjectVersion: "1.2.3",
		ProjectTags: []Tag{
			{Name: "foo"},
			{Name: "bar"},
		},
		IsLatest:   OptionalBoolOf(true),
		AutoCreate: true,
		BOM: base64.StdEncoding.EncodeToString([]byte(`
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "version": 1,
  "components": []
}`)),
	})
	require.NoError(t, err)

	project, err := client.Project.Lookup(context.Background(), "acme-app", "1.2.3")
	require.NoError(t, err)
	require.Contains(t, project.Tags, Tag{Name: "foo"})
	require.Contains(t, project.Tags, Tag{Name: "bar"})
	require.NotNil(t, project.IsLatest)
	require.True(t, *project.IsLatest)
}

func TestBOMService_PostBom(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionBOMUpload,
			PermissionProjectCreationUpload,
			PermissionViewPortfolio,
		},
	})

	_, err := client.BOM.PostBom(context.Background(), BOMUploadRequest{
		ProjectName:    "acme-app",
		ProjectVersion: "1.2.3",
		ProjectTags: []Tag{
			{Name: "foo"},
			{Name: "bar"},
		},
		IsLatest:   OptionalBoolOf(true),
		AutoCreate: true,
		BOM: `
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "version": 1,
  "components": []
}`,
	})
	require.NoError(t, err)

	project, err := client.Project.Lookup(context.Background(), "acme-app", "1.2.3")
	require.NoError(t, err)
	if semver.Compare("v"+client.about.Version, "v4.12") >= 0 {
		require.Contains(t, project.Tags, Tag{Name: "foo"})
		require.Contains(t, project.Tags, Tag{Name: "bar"})
		require.NotNil(t, project.IsLatest)
		require.True(t, *project.IsLatest)
	}
}
