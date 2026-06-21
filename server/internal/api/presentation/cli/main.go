package cli

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"starliner.app/internal/api/presentation/cli/handler"
)

func Register(
	lc fx.Lifecycle,
	sd fx.Shutdowner,
	runner *handler.RunnerHandler,
) {
	rootCmd := &cobra.Command{
		Use:          "starliner",
		Short:        "Starliner administration CLI",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(
		runner.NewRunnerCmd(),
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
	if err == nil {
		return 0
	}
	return 1
}
