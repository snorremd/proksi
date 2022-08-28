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

// Projects related state and types

type project struct { // Embed the rmpb struct to make it extensible for list use
	rmpb.Project
}

type projects struct {
	client   *rm.ProjectsClient
	list     list.Model
	choice   project
	quitting bool
}

type newProjectsMessage struct {
	projects []project
}

// Item delegate for organization struct model
func (p project) Title() string                             { return p.DisplayName }
func (p project) Description() string                       { return p.Description() }
func (p project) FilterValue() string                       { return p.DisplayName }
func (p project) Height() int                               { return 1 }
func (p project) Spacing() int                              { return 0 }
func (p project) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (p project) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	fmt.Fprintf(w, p.DisplayName)
}

// Organizations list View renderer

func (m *projects) Render() string {
	return m.list.View()
}

// Model initializer

func initialProjectsModel(client *rm.ProjectsClient) *projects {
	organizationsList := list.New([]list.Item{}, project{}, defaultWidth, defaultHeight)

	return &projects{
		list:     organizationsList,
		client:   client,
		quitting: false,
	}
}

// State msgs related to projects list state

func (m *projects) NewProjects(msg newProjectsMessage) tea.Msg {

	items := []list.Item{}

	for _, organization := range msg.projects {
		items = append(items, organization)
	}

	m.list.SetItems(items)

	return nil
}

// Commands related to projects list and selection

func (m *projects) GetProjects() tea.Msg {
	ctx := context.Background()
	req := rmpb.ListProjectsRequest{}
	it := m.client.ListProjects(ctx, &req)
	projects := []project{}
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
			break
		}
		projects = append(projects, project{*resp})
	}

	return newProjectsMessage{projects}
}
