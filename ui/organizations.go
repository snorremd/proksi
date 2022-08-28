package ui

import (
	"context"
	"fmt"
	"io"

	rm "cloud.google.com/go/resourcemanager/apiv3"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/iterator"
	rmpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"
)

// Organizations related state and types

type organization struct { // Embed the rmpb struct to make it extensible for list use
	rmpb.Organization
}

type organizations struct {
	client   *rm.OrganizationsClient
	list     list.Model
	choice   organization
	quitting bool
}

type newOrganizationsMsg struct {
	organizations []organization
}

// Item delegate for organization struct model
func (o organization) Title() string                             { return o.DisplayName }
func (o organization) Description() string                       { return o.Organization.Name }
func (o organization) FilterValue() string                       { return o.DisplayName }
func (o organization) Height() int                               { return 1 }
func (o organization) Spacing() int                              { return 0 }
func (o organization) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (o organization) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, _ := listItem.(organization)
	fmt.Fprintln(w, i.Title())
}

// Organizations list View renderer

func (m *organizations) View() string {
	return m.list.View()
}

// Model initializer

func initialOrganizationsModel(client *rm.OrganizationsClient) *organizations {
	organizationsList := list.New([]list.Item{}, organization{}, defaultWidth, defaultHeight)
	organizationsList.Title = "Organizations"

	return &organizations{
		list:     organizationsList,
		client:   client,
		quitting: false,
	}
}

// State msgs related to organizations list state
func onNewOrganizations(msg newOrganizationsMsg, model model) (tea.Model, tea.Cmd) {

	items := []list.Item{}

	for _, organization := range msg.organizations {
		item := list.Item(organization)
		items = append(items, item)
	}

	cmd := model.organizations.list.SetItems(items)

	return model, cmd
}

// Commands related to organizations list and selection

func (m *organizations) GetOrganizations() tea.Msg {
	ctx := context.Background()
	req := rmpb.SearchOrganizationsRequest{}
	it := m.client.SearchOrganizations(ctx, &req)
	organizations := []organization{}
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
			break
		}
		organizations = append(organizations, organization{*resp})
	}

	return newOrganizationsMsg{organizations}
}
