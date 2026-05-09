package cli

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"starliner.app/cli/internal/presentation/cli/handler"
)

func Register(
	lc fx.Lifecycle,
	sd fx.Shutdowner,
	demo *handler.DemoHandler,
	auth *handler.AuthHandler,
) {
	rootCmd := &cobra.Command{
		Version:       "v0.0.1",
		Use:           "starliner",
		Example:       "starliner",
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(
		auth.NewAuthCmd(),
		demo.NewDemoCmd(),
	)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				err := rootCmd.ExecuteContext(context.Background())
				_ = sd.Shutdown(fx.ExitCode(exitCodeFor(err)))
			}()
			return nil
		},
	})
}

func exitCodeFor(err error) int {
	if err != nil {
		return 1
	}
	return 0
}
