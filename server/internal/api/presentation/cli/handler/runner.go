package handler

import (
	"errors"
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
)

type RunnerHandler struct {
	runnerApplication *application.RunnerApplication
}

func NewRunnerHandler(runnerApplication *application.RunnerApplication) *RunnerHandler {
	return &RunnerHandler{
		runnerApplication: runnerApplication,
	}
}

func (h *RunnerHandler) NewRunnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Manage global runners",
	}

	cmd.AddCommand(
		h.newCreateGlobalCmd(),
		h.newListGlobalCmd(),
		h.newDeleteGlobalCmd(),
	)

	return cmd
}

func (h *RunnerHandler) newCreateGlobalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-global",
		Short: "Create a global runner and print its registration token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := h.runnerApplication.CreateGlobalRunner(cmd.Context())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if _, err := fmt.Fprintf(out, "Created global runner id=%d\n", result.Id); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(out, "Registration token (shown once, expires at %s):\n", result.ExpiresAt.Format("2006-01-02 15:04:05 MST")); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(out, result.Token); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(out, "Register with: sudo runner register --token <token>"); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (h *RunnerHandler) newListGlobalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-global",
		Short: "List all global runners",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			runners, err := h.runnerApplication.ListGlobalRunners(cmd.Context())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if len(runners) == 0 {
				if _, err := fmt.Fprintln(out, "No global runners found."); err != nil {
					return err
				}
				return nil
			}

			w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
			if _, err := fmt.Fprintln(w, "ID\tNAME\tSTATUS\tLABELS\tMAX JOBS"); err != nil {
				return err
			}
			for _, runner := range runners {
				name := ""
				if runner.Name != nil {
					name = *runner.Name
				}
				if _, err := fmt.Fprintf(w, "%d\t%s\t%s\t%v\t%d\n", runner.Id, name, runner.Status, runner.Labels, runner.MaxConcurrentJobs); err != nil {
					return err
				}
			}

			return w.Flush()
		},
	}

	return cmd
}

func (h *RunnerHandler) newDeleteGlobalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-global [runner-id]",
		Short: "Delete a global runner by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runnerId, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid runner id: %w", err)
			}

			if err := h.runnerApplication.DeleteGlobalRunner(cmd.Context(), runnerId); err != nil {
				if errors.Is(err, value.ErrRunnerNotFound) {
					return fmt.Errorf("global runner %d not found", runnerId)
				}
				return err
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Deleted global runner id=%d\n", runnerId); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
