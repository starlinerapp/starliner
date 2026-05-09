package handler

import (
	"github.com/spf13/cobra"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Long:  "Manage authentication with Starliner.",
	}

	cmd.AddCommand(
		ah.newAuthLoginCmd(),
		ah.newAuthLogoutCmd(),
	)

	return cmd
}

func (ah *AuthHandler) newAuthLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to Starliner",
		Long:    "Login to Starliner.",
		Example: "starliner auth login",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
	}
	return cmd
}

func (ah *AuthHandler) newAuthLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Logout from Starliner",
		Long:    "Logout from Starliner.",
		Example: "starliner auth logout",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
	}
	return cmd
}
