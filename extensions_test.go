package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtensionsList(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		Version: "v5-latest",
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})
	points, err := client.Extension.ListExtensionPoints(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, points.Items)
	require.Nil(t, points.NextPageToken)

	for _, point := range points.Items {
		extensions, err := client.Extension.ListExtensions(ctx, point.Name)
		require.NoError(t, err)
		require.NotEmpty(t, extensions.Items)
		require.Nil(t, extensions.NextPageToken)
		for _, extension := range extensions.Items {
			require.NotZero(t, extension.Name)
		}
	}
}

func TestExtensionsConfig(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		Version: "v5-latest",
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})
	points, err := client.Extension.ListExtensionPoints(ctx)
	require.NoError(t, err)

	for _, point := range points.Items {
		extensions, err := client.Extension.ListExtensions(ctx, point.Name)
		require.NoError(t, err)
		for _, extension := range extensions.Items {
			if !extension.Configurable {
				continue
			}
			identity := ExtensionIdentity{
				ExtensionPointName: point.Name,
				ExtensionName:      extension.Name,
			}
			{
				config, err := client.Extension.GetExtensionConfig(ctx, identity)
				require.NoError(t, err)
				require.NotZero(t, config)
			}
			{
				schema, err := client.Extension.GetExtensionConfigSchema(ctx, identity)
				require.NoError(t, err)
				require.NotZero(t, schema.SchemaURL)
				require.NotZero(t, schema.Schema)
			}
		}
	}
}

func TestExtensionsConfigUpdate(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		Version: "v5-latest",
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})
	identity := ExtensionIdentity{
		ExtensionPointName: "notification-publisher",
		ExtensionName:      "email",
	}
	config, err := client.Extension.GetExtensionConfig(ctx, identity)
	require.NoError(t, err)

	sslEnabled, ok := config.Config["sslEnabled"].(bool)
	require.True(t, ok)
	require.False(t, sslEnabled)

	// Change
	{
		config.Config["sslEnabled"] = !sslEnabled
		err = client.Extension.UpdateExtensionConfig(ctx, identity, config)
		require.NoError(t, err)

		updated, err := client.Extension.GetExtensionConfig(ctx, identity)
		require.NoError(t, err)

		sslEnabledUpdated, ok := updated.Config["sslEnabled"].(bool)
		require.True(t, ok)
		require.NotEqual(t, sslEnabled, sslEnabledUpdated)
	}
	// Restore
	{
		config.Config["sslEnabled"] = sslEnabled
		err = client.Extension.UpdateExtensionConfig(ctx, identity, config)
		require.NoError(t, err)

		updated, err := client.Extension.GetExtensionConfig(ctx, identity)
		require.NoError(t, err)

		sslEnabledUpdated, ok := updated.Config["sslEnabled"].(bool)
		require.True(t, ok)
		require.Equal(t, sslEnabled, sslEnabledUpdated)
	}
	// No Change
	{
		err = client.Extension.UpdateExtensionConfig(ctx, identity, config)
		require.EqualError(t, err, "api error (status: 304)")

		unchanged, err := client.Extension.GetExtensionConfig(ctx, identity)
		require.NoError(t, err)

		sslEnabledUnchanged, ok := unchanged.Config["sslEnabled"].(bool)
		require.True(t, ok)
		require.Equal(t, sslEnabled, sslEnabledUnchanged)
	}
}
