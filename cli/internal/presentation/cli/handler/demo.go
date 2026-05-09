package handler

import (
	"github.com/spf13/cobra"
	"starliner.app/cli/internal/application"
)

type DemoHandler struct {
	demoApplication *application.DemoApplication
}

func NewDemoHandler(demoApplication *application.DemoApplication) *DemoHandler {
	return &DemoHandler{
		demoApplication: demoApplication,
	}
}

func (dh *DemoHandler) NewDemoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Manage the local Starliner demo",
		Long:  "Manage the Starliner demo environment on your local machine.",
	}

	cmd.AddCommand(
		dh.newDemoUpCmd(),
		dh.newDemoDownCmd(),
	)

	return cmd
}

func (dh *DemoHandler) newDemoUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up",
		Short:   "Start the local Starliner demo",
		Long:    "Start the Starliner demo environment on your local machine.",
		Example: "starliner demo up",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := dh.demoApplication.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

func (dh *DemoHandler) newDemoDownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Short:   "Stop the local Starliner demo",
		Long:    "Stop the Starliner demo environment on your local machine.",
		Example: "starliner demo down",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dh.demoApplication.Stop(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
