package cmd

import (
	"context"
	"os"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	run "cloud.google.com/go/run/apiv2"
	"github.com/snorremd/proksi/ui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "proksi",
	Short: "Creates a http proxy injecting your identity token into the Authorization header",
	Long: `Proksi is a simple command line and TUI application that
creates a http proxy injecting your identity token into the Authorization header.
It filters out requests for URLs not matching a user specified whitelist.
The TUI provides a simple  interface to configure the proxy.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		organizationsClient, err := resourcemanager.NewOrganizationsClient(ctx)
		if err != nil {
			return err
		}

		projectsClient, err := resourcemanager.NewProjectsClient(ctx)
		if err != nil {
			return err
		}

		servicesClient, err := run.NewServicesClient(ctx)
		if err != nil {
			return err
		}

		defer projectsClient.Close()
		p := ui.NewProgram(ui.ProgramConfig{
			OrganizationsClient: organizationsClient,
			ProjectsClient:      projectsClient,
			ServicesClient:      servicesClient,
			Ctx:                 &ctx,
		})
		return p.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.proksi.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
