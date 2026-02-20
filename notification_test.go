package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublishers(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})
	// Create
	publisher, err := client.Notification.CreatePublisher(ctx, NotificationPublisher{
		Name:             "Test_Publisher",
		Description:      "Test_Description",
		PublisherClass:   "org.dependencytrack.notification.publisher.ConsolePublisher",
		TemplateMIMEType: "text/plain",
		Template:         "Test_Template",
	})
	{
		require.NoError(t, err)
		require.NotZero(t, publisher.UUID)
		require.Equal(t, publisher.Name, "Test_Publisher")
		require.Equal(t, publisher.Description, "Test_Description")
		require.Equal(t, publisher.PublisherClass, "org.dependencytrack.notification.publisher.ConsolePublisher")
		require.Equal(t, publisher.TemplateMIMEType, "text/plain")
		require.Equal(t, publisher.Template, "Test_Template")
		require.Equal(t, publisher.DefaultPublisher, false)
	}
	// Update
	{
		updatedReq := publisher
		updatedReq.Description = "Test_Updated_Description"
		updated, err := client.Notification.UpdatePublisher(ctx, updatedReq)
		require.NoError(t, err)
		require.Equal(t, updated.UUID, publisher.UUID)
		require.Equal(t, updated.Name, publisher.Name)
		require.Equal(t, updated.Description, "Test_Updated_Description")
		require.Equal(t, updated.PublisherClass, publisher.PublisherClass)
		require.Equal(t, updated.Template, publisher.Template)
		require.Equal(t, updated.TemplateMIMEType, publisher.TemplateMIMEType)
		require.Equal(t, updated.DefaultPublisher, publisher.DefaultPublisher)
	}
	// Fetch
	{
		allPublishers, err := client.Notification.GetAllPublishers(ctx)
		require.NoError(t, err)
		found := NotificationPublisher{}
		for _, pub := range allPublishers {
			if pub.UUID == publisher.UUID {
				found = pub
				break
			}
		}
		require.Equal(t, found.UUID, publisher.UUID)
		require.Equal(t, found.Name, publisher.Name)
		require.Equal(t, found.Description, "Test_Updated_Description")
		require.Equal(t, found.PublisherClass, publisher.PublisherClass)
		require.Equal(t, found.Template, publisher.Template)
		require.Equal(t, found.TemplateMIMEType, publisher.TemplateMIMEType)
		require.Equal(t, found.DefaultPublisher, publisher.DefaultPublisher)
	}
	// Delete
	{
		err := client.Notification.DeletePublisher(ctx, publisher.UUID)
		require.NoError(t, err)
	}
	// Check absence
	{
		allPublishers, err := client.Notification.GetAllPublishers(ctx)
		require.NoError(t, err)
		found := NotificationPublisher{}
		for _, pub := range allPublishers {
			if pub.UUID == publisher.UUID {
				found = pub
				break
			}
		}
		require.Zero(t, found.UUID)
	}
}
