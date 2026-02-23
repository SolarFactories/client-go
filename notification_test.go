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

func TestRules(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})
	publisher, err := client.Notification.CreatePublisher(ctx, NotificationPublisher{
		Name:             "Test_Rule_Publisher",
		Description:      "Test_Rule_Description",
		PublisherClass:   "org.dependencytrack.notification.publisher.ConsolePublisher",
		TemplateMIMEType: "text/plain",
		Template:         "Test_Rule_Template",
	})
	require.NoError(t, err)
	// Create
	rule, err := client.Notification.CreateRule(ctx, NotificationRule{
		Name:        "Test_Rule_Name",
		Scope:       NotificationRuleScopePortfolio,
		TriggerType: NotificationRuleTriggerTypeEvent,
		Publisher:   publisher,
	})
	{
		require.NoError(t, err)
		require.NotZero(t, rule.UUID)
		require.Equal(t, rule.Name, "Test_Rule_Name")
		require.Equal(t, rule.Scope, NotificationRuleScopePortfolio)
		require.Equal(t, rule.TriggerType, NotificationRuleTriggerTypeEvent)
		require.Equal(t, rule.Enabled, true)
		require.Equal(t, rule.NotifyChildren, true)
		require.Equal(t, rule.LogSuccessfulPublish, false)
		require.Empty(t, rule.NotificationLevel)
		require.Empty(t, rule.NotifyOn)
		require.Empty(t, rule.Message)
		require.Empty(t, rule.PublisherConfig)
		require.Empty(t, rule.ScheduleLastTriggeredAt)
		require.Empty(t, rule.ScheduleNextTriggerAt)
		require.Empty(t, rule.ScheduleCron)
		require.Empty(t, rule.ScheduleSkipUnchanged)
		require.Equal(t, rule.Publisher.UUID, publisher.UUID)
		require.Empty(t, rule.Projects)
		require.Empty(t, rule.Tags)
		require.Empty(t, rule.Teams)
	}
	// Update
	{
		updatedReq := rule
		updatedReq.PublisherConfig = "{\"Key\": \"Publisher Config\"}"
		updatedRes, err := client.Notification.UpdateRule(ctx, updatedReq)
		require.NoError(t, err)
		require.Equal(t, updatedRes, updatedReq)
	}
	// Fetch
	{
		allRules, err := FetchAll(func(po PageOptions) (Page[NotificationRule], error) {
			return client.Notification.GetAllRules(ctx, po, SortOptions{}, GetAllRulesFilterOptions{})
		})
		require.NoError(t, err)
		found := NotificationRule{}
		for _, rule_ := range allRules {
			if rule_.UUID == rule.UUID {
				found = rule_
				break
			}
		}
		require.Equal(t, found.UUID, rule.UUID)
		require.Equal(t, found.Name, rule.Name)
		require.Equal(t, found.Enabled, rule.Enabled)
		require.Equal(t, found.NotifyChildren, rule.NotifyChildren)
		require.Equal(t, found.LogSuccessfulPublish, rule.LogSuccessfulPublish)
		require.Equal(t, found.Scope, rule.Scope)
		require.Equal(t, found.NotificationLevel, rule.NotificationLevel)
		require.Equal(t, found.NotifyOn, rule.NotifyOn)
		require.Equal(t, found.TriggerType, rule.TriggerType)
		require.Equal(t, found.Message, rule.Message)
		require.Equal(t, found.PublisherConfig, "{\"Key\": \"Publisher Config\"}")
		require.Equal(t, found.ScheduleLastTriggeredAt, rule.ScheduleLastTriggeredAt)
		require.Equal(t, found.ScheduleNextTriggerAt, rule.ScheduleNextTriggerAt)
		require.Equal(t, found.ScheduleCron, rule.ScheduleCron)
		require.Equal(t, found.ScheduleSkipUnchanged, rule.ScheduleSkipUnchanged)
		require.Equal(t, found.Publisher, rule.Publisher)
		require.Equal(t, found.Projects, rule.Projects)
		require.Equal(t, found.Tags, rule.Tags)
		require.Equal(t, found.Teams, rule.Teams)
	}
	// Delete
	{
		err := client.Notification.DeleteRule(ctx, rule)
		require.NoError(t, err)
	}
	// Check Absence
	{
		allRules, err := FetchAll(func(po PageOptions) (Page[NotificationRule], error) {
			return client.Notification.GetAllRules(ctx, po, SortOptions{}, GetAllRulesFilterOptions{})
		})
		require.NoError(t, err)
		found := NotificationRule{}
		for _, rule_ := range allRules {
			if rule_.UUID == rule.UUID {
				found = rule_
				break
			}
		}
		require.Zero(t, found.UUID)
	}
}

func TestRuleProjects(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
			PermissionPortfolioManagement,
		},
	})
	publisher, err := client.Notification.CreatePublisher(ctx, NotificationPublisher{
		Name:             "Test_Rule_Projects_Publisher",
		Description:      "Test_Rule_Description",
		PublisherClass:   "org.dependencytrack.notification.publisher.ConsolePublisher",
		TemplateMIMEType: "text/plain",
		Template:         "Test_Rule_Template",
	})
	require.NoError(t, err)
	rule, err := client.Notification.CreateRule(ctx, NotificationRule{
		Name:        "Test_Rule_Projects_Name",
		Scope:       NotificationRuleScopePortfolio,
		TriggerType: NotificationRuleTriggerTypeEvent,
		Publisher:   publisher,
	})
	require.NoError(t, err)
	project, err := client.Project.Create(ctx, Project{
		Name: "Test_Rule_Projects_Project",
	})
	require.NoError(t, err)
	// Add Project
	{
		updated, err := client.Notification.AddProjectToRule(ctx, rule.UUID, project.UUID)
		require.NoError(t, err)
		require.Equal(t, updated.UUID, rule.UUID)
		require.Equal(t, updated.Projects, []Project{project})
	}
	// Fetch
	{
		allRules, err := FetchAll(func(po PageOptions) (Page[NotificationRule], error) {
			return client.Notification.GetAllRules(ctx, po, SortOptions{}, GetAllRulesFilterOptions{})
		})
		require.NoError(t, err)
		found := NotificationRule{}
		for _, rule_ := range allRules {
			if rule_.UUID == rule.UUID {
				found = rule_
				break
			}
		}
		require.NotZero(t, found.UUID)
		require.Empty(t, project.Tags)
		require.Empty(t, project.Properties)

		project.Tags = nil
		project.Properties = nil
		require.Equal(t, found.Projects, []Project{project})
	}
	// Remove Project
	{
		updated, err := client.Notification.RemoveProjectFromRule(ctx, rule.UUID, project.UUID)
		require.NoError(t, err)
		require.Equal(t, updated, rule)
	}
	// Check Absence
	{
		allRules, err := FetchAll(func(po PageOptions) (Page[NotificationRule], error) {
			return client.Notification.GetAllRules(ctx, po, SortOptions{}, GetAllRulesFilterOptions{})
		})
		require.NoError(t, err)
		found := NotificationRule{}
		for _, rule_ := range allRules {
			if rule_.UUID == rule.UUID {
				found = rule_
				break
			}
		}
		require.Empty(t, found.Projects)
	}
}

func TestRuleTeams(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
			PermissionAccessManagement,
		},
	})
	publisher, err := client.Notification.CreatePublisher(ctx, NotificationPublisher{
		Name:             "Test_Rule_Tags_Publisher",
		Description:      "Test_Rule_Description",
		PublisherClass:   "org.dependencytrack.notification.publisher.SendMailPublisher",
		TemplateMIMEType: "text/plain",
		Template:         "Test_Rule_Template",
	})
	require.NoError(t, err)
	rule, err := client.Notification.CreateRule(ctx, NotificationRule{
		Name:        "Test_Rule_Tags_Name",
		Scope:       NotificationRuleScopePortfolio,
		TriggerType: NotificationRuleTriggerTypeEvent,
		Publisher:   publisher,
	})
	require.NoError(t, err)
	team, err := client.Team.Create(ctx, Team{
		Name: "Test_Rule_Teams_Team",
	})
	require.NoError(t, err)
	// Add Team
	{
		updated, err := client.Notification.AddTeamToRule(ctx, rule.UUID, team.UUID)
		require.NoError(t, err)

		require.Empty(t, team.APIKeys)
		require.Empty(t, team.MappedOIDCGroups)
		team.APIKeys = nil
		team.MappedOIDCGroups = nil

		require.Equal(t, updated.Teams, []Team{team})
		updated.Teams = []Team{}
		require.Equal(t, updated, rule)
	}
	// Fetch
	{
		allRules, err := FetchAll(func(po PageOptions) (Page[NotificationRule], error) {
			return client.Notification.GetAllRules(ctx, po, SortOptions{}, GetAllRulesFilterOptions{})
		})
		require.NoError(t, err)
		found := NotificationRule{}
		for _, rule_ := range allRules {
			if rule_.UUID == rule.UUID {
				found = rule_
				break
			}
		}
		require.Empty(t, team.Permissions)
		team.Permissions = nil
		require.Equal(t, found.Teams, []Team{team})
	}
	// Remove Team
	{
		updated, err := client.Notification.RemoveTeamFromRule(ctx, rule.UUID, team.UUID)
		require.NoError(t, err)
		require.Equal(t, updated, rule)
	}
	// Check Absence
	{
		allRules, err := FetchAll(func(po PageOptions) (Page[NotificationRule], error) {
			return client.Notification.GetAllRules(ctx, po, SortOptions{}, GetAllRulesFilterOptions{})
		})
		require.NoError(t, err)
		found := NotificationRule{}
		for _, rule_ := range allRules {
			if rule_.UUID == rule.UUID {
				found = rule_
				break
			}
		}
		require.Empty(t, found.Teams)
	}
}
