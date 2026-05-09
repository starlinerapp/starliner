package handler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"starliner.app/cli/internal/application"
)

type AuthHandler struct {
	authApplication *application.AuthApplication
}

func NewAuthHandler(authApplication *application.AuthApplication) *AuthHandler {
	return &AuthHandler{
		authApplication: authApplication,
	}
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
			in := cmd.InOrStdin()
			out := cmd.OutOrStdout()

			reader := bufio.NewReader(in)

			if _, err := fmt.Fprint(out, "Email: "); err != nil {
				return err
			}
			email, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			if _, err := fmt.Fprint(out, "Password: "); err != nil {
				return err
			}

			stdinFile, ok := in.(*os.File)
			if !ok {
				return fmt.Errorf("password input requires a terminal")
			}

			passwordBytes, err := term.ReadPassword(int(stdinFile.Fd()))
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintln(out); err != nil {
				return err
			}

			if err := ah.authApplication.Login(
				cmd.Context(),
				strings.TrimSpace(email),
				string(passwordBytes),
			); err != nil {
				return err
			}

			if _, err := fmt.Fprintln(out, "Successfully logged in."); err != nil {
				return err
			}
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
			if err := ah.authApplication.Logout(cmd.Context()); err != nil {
				return err
			}
			fmt.Println("Successfully logged out.")
			return nil
		},
	}
	return cmd
}
