package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLicenseGroupGetAll(t *testing.T) {
	po := PageOptions{
		// DepedencyTrack is preloaded with 4 license groups. This need only be <= the number preloaded in DT.
		PageSize:   4,
		PageNumber: 1,
	}
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPolicyManagement,
		},
	})

	groups, err := client.LicenseGroup.GetAll(context.Background(), po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, len(groups.Items), po.PageSize)
	require.NotZero(t, groups.TotalCount)
	for _, group := range groups.Items {
		require.NotZero(t, group.UUID)
		require.NotEmpty(t, group.Name)
		require.NotEmpty(t, group.Licenses)
	}
}

func TestLicenseGroupLifecycle(t *testing.T) {
	po := PageOptions{
		PageNumber: 1,
		PageSize:   10,
	}
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPolicyManagement,
		},
	})
	// Check absence
	{
		groups, err := client.LicenseGroup.GetAll(context.Background(), po, SortOptions{})
		require.NoError(t, err)
		for _, group := range groups.Items {
			require.NotEqual(t, group.Name, "TestLicenseGroupLifecycle")
		}
	}
	// Create
	group, err := client.LicenseGroup.Create(context.Background(), LicenseGroup{
		Name: "TestLicenseGroupLifecycle",
	})
	{
		require.NoError(t, err)
		require.NotZero(t, group.UUID)
		require.Equal(t, group.Name, "TestLicenseGroupLifecycle")
	}
	// Check presence
	{
		newGroup, err := client.LicenseGroup.Get(context.Background(), group.UUID)
		require.NoError(t, err)
		require.Equal(t, newGroup, group)
	}
	// Update
	{
		copyGroup := group
		copyGroup.Name = "UpdatedName"
		require.NotEqual(t, copyGroup.Name, group.Name)

		newGroup, err := client.LicenseGroup.Update(context.Background(), copyGroup)
		require.NoError(t, err)
		require.Equal(t, newGroup, copyGroup)
	}
	// Check update
	{
		newGroup, err := client.LicenseGroup.Get(context.Background(), group.UUID)
		require.NoError(t, err)
		require.Equal(t, newGroup.Name, "UpdatedName")
		newGroup.Name = group.Name
		require.Equal(t, newGroup, group)
	}
	// Delete
	{
		err := client.LicenseGroup.Delete(context.Background(), group.UUID)
		require.NoError(t, err)
	}
	// Check absence
	{
		groups, err := client.LicenseGroup.GetAll(context.Background(), po, SortOptions{})
		require.NoError(t, err)
		for _, check := range groups.Items {
			require.NotEqual(t, check.UUID, group.UUID)
		}
	}
}

func TestLicenseGroupLicense(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
			PermissionPolicyManagement,
		},
	})
	// Create License, License Group
	group, err := client.LicenseGroup.Create(context.Background(), LicenseGroup{
		Name: "TestLicenseGroup",
	})
	require.NoError(t, err)
	license, err := client.License.Create(context.Background(), License{
		Name:      "TestLicense",
		LicenseID: "TestLicenseID",
	})
	require.NoError(t, err)
	// Add license
	{
		newGroup, err := client.LicenseGroup.AddLicense(context.Background(), group.UUID, license.UUID)
		require.NoError(t, err)
		require.Equal(t, newGroup.Licenses, []License{license})
		newGroup.Licenses = []License{}
		require.Equal(t, newGroup, group)
	}
	// Check presence in group
	{
		newGroup, err := client.LicenseGroup.Get(context.Background(), group.UUID)
		require.NoError(t, err)
		require.Equal(t, newGroup.Licenses, []License{license})
		newGroup.Licenses = []License{}
		require.Equal(t, newGroup, group)
	}
	// Remove license
	{
		newGroup, err := client.LicenseGroup.RemoveLicense(context.Background(), group.UUID, license.UUID)
		require.NoError(t, err)
		require.Empty(t, newGroup.Licenses)
		require.Equal(t, newGroup, group)
	}
	// Check empty group
	{
		newGroup, err := client.LicenseGroup.Get(context.Background(), group.UUID)
		require.NoError(t, err)
		require.Empty(t, newGroup.Licenses)
		require.Equal(t, newGroup, group)
	}
	// Cleanup
	{
		err := client.License.Delete(context.Background(), "TestLicenseID")
		require.NoError(t, err)

		err = client.LicenseGroup.Delete(context.Background(), group.UUID)
		require.NoError(t, err)
	}
}
