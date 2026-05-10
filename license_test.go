package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLicenseGetAll(t *testing.T) {
	po := PageOptions{
		PageSize:   10,
		PageNumber: 1,
	}
	client := setUpContainer(t, testContainerOptions{})

	licenses, err := client.License.GetAll(context.Background(), po)
	require.NoError(t, err)
	require.Equal(t, len(licenses.Items), po.PageSize)
	require.NotZero(t, licenses.TotalCount)
	for _, license := range licenses.Items {
		require.NotZero(t, license.UUID)
		require.NotEmpty(t, license.Name)
		require.NotEmpty(t, license.LicenseID)
	}
}

func TestLicenseLifecycle(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})
	licenseSPDX := "DependencyTrack-ClientGo"
	// Check absence
	{
		l, err := client.License.Get(context.Background(), licenseSPDX)
		require.Error(t, err)
		require.Zero(t, l)
	}
	// Create
	license, err := client.License.Create(context.Background(), License{
		Name:      "Client Go Test License",
		LicenseID: licenseSPDX,
	})
	{
		require.NoError(t, err)
		require.NotZero(t, license.UUID)
		require.Equal(t, license.LicenseID, licenseSPDX)
		require.Equal(t, license.Name, "Client Go Test License")
	}
	// Check presence
	{
		l, err := client.License.Get(context.Background(), licenseSPDX)
		require.NoError(t, err)
		require.Equal(t, l, license)
	}
	// Delete
	{
		err := client.License.Delete(context.Background(), licenseSPDX)
		require.NoError(t, err)
	}
	// Check absence
	{
		l, err := client.License.Get(context.Background(), licenseSPDX)
		require.Error(t, err)
		require.Zero(t, l)
	}
}

func TestLicenseGetConcise(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{})
	licenses, err := client.License.GetConcise(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, licenses)
	for _, license := range licenses {
		require.NotZero(t, license.UUID)
		require.NotEmpty(t, license.Name)
		require.NotEmpty(t, license.LicenseID)
	}
}
