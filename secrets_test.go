package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecretsLifecycle(t *testing.T) {
	ctx := context.Background()
	name := "test-secret-name"
	client := setUpContainer(t, testContainerOptions{
		Version: "v5-latest",
		APIPermissions: []string{
			PermissionSystemConfiguration,
			PermissionSecretManagement,
		},
	})
	{
		url, err := client.Secrets.Create(ctx, CreateSecretRequest{
			Name:  name,
			Value: "test-secret-value",
		})
		require.NoError(t, err)
		require.NotZero(t, url)
	}
	{
		md, err := client.Secrets.GetMetadata(ctx, name)
		require.NoError(t, err)
		require.Equal(t, md.Name, name)
		require.Zero(t, md.Description)
		require.NotNil(t, md.CreatedAt)
		require.Nil(t, md.UpdatedAt)
	}
	{
		err := client.Secrets.Update(ctx, name, UpdateSecretRequest{})
		require.NoError(t, err)
	}
	{
		description := "test-secret-description"
		err := client.Secrets.Update(ctx, name, UpdateSecretRequest{
			Description: &description,
		})
		require.NoError(t, err)
	}
	{
		mds, err := client.Secrets.ListMetadata(ctx, ListSecretMetadataQuery{Filter: name})
		require.NoError(t, err)
		require.Equal(t, len(mds.Items), 1)
		require.Equal(t, mds.Total.Count, uint64(len(mds.Items)))
		md := mds.Items[0]
		require.Equal(t, md.Description, "test-secret-description")
		require.NotNil(t, md.UpdatedAt)
	}
	{
		err := client.Secrets.Delete(ctx, name)
		require.NoError(t, err)
	}
}
