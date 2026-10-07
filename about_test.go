package dtrack

import (
	"context"
	"fmt"
	"os"

	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/mod/semver"

	"github.com/stretchr/testify/require"
)

func TestAboutService_Get(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{})

	about, err := client.About.Get(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, about)

	require.NotEmpty(t, about.Timestamp)
	require.NotEmpty(t, about.Version)
	require.NotEqual(t, uuid.Nil, about.UUID)
	if semver.Compare(about.Version, "5") < 0 {
		require.NotEqual(t, uuid.Nil, about.SystemUUID)
	}
	require.Equal(t, "Dependency-Track", about.Application)

	require.NotEmpty(t, about.Framework.Timestamp)
	require.NotEmpty(t, about.Framework.Version)
	require.NotEqual(t, uuid.Nil, about.Framework.UUID)
	require.Equal(t, "Alpine", about.Framework.Name)
}

type testContainerOptions struct {
	Version        string
	APIPermissions []string
}

func setUpContainer(t *testing.T, options testContainerOptions) *Client {
	ctx := context.Background()
	host := os.Getenv("DEPENDENCYTRACK_API_HOST")
	key := os.Getenv("DEPENDENCYTRACK_API_KEY")

	version := "latest"
	if options.Version != "" {
		version = options.Version
	}

	fmt.Printf("Using host: %v, key: %v\n.", host, key)
	if len(host) == 0 {
		fmt.Println("Starting test container")
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image: "dependencytrack/apiserver:" + version,
				Env: map[string]string{
					"JAVA_OPTIONS":                         "-Xmx1g",
					"SYSTEM_REQUIREMENT_CHECK_ENABLED":     "false",
					"TELEMETRY_SUBMISSION_ENABLED_DEFAULT": "false",
				},
				ExposedPorts: []string{"8080/tcp"},
				WaitingFor:   wait.ForLog("Dependency-Track is ready"),
			},
			Started: true,
		})
		require.NoError(t, err)

		t.Cleanup(func() {
			err = container.Terminate(ctx)
			if err != nil {
				log.Fatalf("failed to terminate container: %v", err)
			}
		})

		host, err = container.Endpoint(ctx, "http")
		require.NoError(t, err)
	}

	var tmpClient *Client
	var err error
	if len(key) > 0 {
		fmt.Println("Using existing key")
		tmpClient, err = NewClient(host, WithAPIKey(key))
		require.NoError(t, err)
	} else {
		fmt.Println("Bootstrapping authentication")
		client, err := NewClient(host)
		require.NoError(t, err)

		err = client.User.ForceChangePassword(ctx, "admin", "admin", "test")
		require.NoError(t, err)

		bearerToken, err := client.User.Login(ctx, "admin", "test")
		require.NoError(t, err)

		tmpClient, err = NewClient(host, WithBearerToken(bearerToken))
		require.NoError(t, err)
	}
	if version != "latest" {
		require.Equal(t, tmpClient.about.Version, version)
	}
	team, err := tmpClient.Team.Create(ctx, Team{Name: "test"})
	require.NoError(t, err)
	t.Cleanup(func() {
		err = tmpClient.Team.Delete(ctx, team)
		if err != nil {
			log.Fatalf("failed to delete temporary team: %v", err)
		}
	})

	for _, permissionName := range options.APIPermissions {
		_, err = tmpClient.Permission.AddPermissionToTeam(ctx, Permission{Name: permissionName}, team.UUID)
		require.NoError(t, err)
	}

	apiKey, err := tmpClient.Team.GenerateAPIKey(ctx, team.UUID)
	require.NoError(t, err)

	client, err := NewClient(host, WithAPIKey(apiKey.Key))
	require.NoError(t, err)

	return client
}
