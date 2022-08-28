package ui

import (
	"context"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	run "cloud.google.com/go/run/apiv2"
	tea "github.com/charmbracelet/bubbletea"
	runpb "google.golang.org/genproto/googleapis/cloud/run/v2"
)

func NewProgram(config ProgramConfig) *tea.Program {
	return tea.NewProgram(initialModel(config))
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.organizations.GetOrganizations
}

func initialModel(config ProgramConfig) tea.Model {
	return &model{
		organizations: initialOrganizationsModel(config.OrganizationsClient),
		projects:      initialProjectsModel(config.ProjectsClient),
		services:      []runpb.Service{},
		serverLog:     []string{},
	}
}

type ProgramConfig struct {
	OrganizationsClient *resourcemanager.OrganizationsClient
	ProjectsClient      *resourcemanager.ProjectsClient
	ServicesClient      *run.ServicesClient
	Ctx                 *context.Context
}

// The main UI state
type model struct {
	organizations *organizations
	projects      *projects

	// Service state for handling services list
	services []runpb.Service

	serverLog []string
}

func (m model) View() string {
	return "\n" +
		m.organizations.View()

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case newOrganizationsMsg:
		return onNewOrganizations(msg, m)

	default:
		return m, nil
	}

	return m, nil

}
