// Copyright 2024 The Drasi Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"drasi.io/cli/sdk"
	"drasi.io/cli/sdk/registry"
	"github.com/spf13/cobra"
)

func NewNamespaceCommand() *cobra.Command {
	var namespaceCommand = &cobra.Command{
		Use:   "namespace",
		Short: "Manage CLI namespace settings",
		Long:  `Manage the default Kubernetes namespace used by the Drasi CLI.`,
	}
	namespaceCommand.AddCommand(setNamespaceCommand())
	namespaceCommand.AddCommand(getNamespaceCommand())
	namespaceCommand.AddCommand(listNamespaceCommand())
	return namespaceCommand
}

func setNamespaceCommand() *cobra.Command {
	var setNamespaceCommand = &cobra.Command{
		Use:   "set [namespace]",
		Short: "Set the default Drasi environment",
		Long: `Set the default namespace used as the target for all Drasi CLI commands. 
This namespace is used as the target for all future Drasi CLI commands unless overridden using the '-n' flag.
If both a default namespace is configured and the '-n' flag is used, the '-n' flag takes precedence.

If a default namespace is never set, the Drasi CLI will use the default namespace 'drasi-system'.

Arguments:
  namespace  The name of the Kubernetes namespace to configure as the default Drasi environment.
             This commands assumes that Drasi has been installed to the namespace specified and does not verify it is there.

Usage examples:
  drasi namespace get
  drasi namespace set my-namespace
  drasi namespace list
`,
		Args: cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("too many arguments")
			}
			var err error
			var namespace string

			if len(args) == 1 {
				namespace = args[0]
			}

			if cmd.Flags().Changed("namespace") {
				if namespace, err = cmd.Flags().GetString("namespace"); err != nil {
					return err
				}
			}

			if namespace != "" {
				current, err := registry.GetCurrentRegistration()
				if err != nil {
					return err
				}
				if current == "" {
					return fmt.Errorf("no current registration found")
				}
				reg, err := registry.LoadRegistration(current)
				if err != nil {
					return err
				}

				if k8sConfig, ok := reg.(*registry.KubernetesConfig); ok {
					k8sConfig.Namespace = namespace
					if err := registry.SaveRegistration(current, k8sConfig); err != nil {
						return err
					}
					fmt.Println("Namespace set to " + namespace)
				} else {
					return fmt.Errorf("current registration is not for Kubernetes")
				}

			} else {
				return fmt.Errorf("namespace cannot be empty")
			}

			return nil
		},
	}
	return setNamespaceCommand
}

func getNamespaceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show the current default Drasi environment",
		Long:  `Get the current default namespace used for all Drasi CLI commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := registry.LoadCurrentRegistration()
			if err != nil {
				return err
			}

			if k8sConfig, ok := reg.(*registry.KubernetesConfig); ok {
				fmt.Println("Current namespace: " + k8sConfig.Namespace)
			} else {
				return fmt.Errorf("current registration is not for Kubernetes")
			}

			return nil
		},
	}
}

func listNamespaceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all Drasi environments",
		Long:  `List all namespaces on the default Kubernetes cluster that have Drasi installed in them.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			reg, err := registry.LoadCurrentRegistration()
			if err != nil {
				return err
			}

			platformClient, err := sdk.NewPlatformClient(reg)
			if err != nil {
				return err
			}

			if k8sClient, ok := platformClient.(*sdk.KubernetesPlatformClient); ok {
				namespaces, err := k8sClient.ListNamespaces()
				if err != nil {
					return err
				}

				fmt.Println("Namespaces:")
				for _, ns := range namespaces {
					fmt.Println(ns)
				}
			}

			return nil
		},
	}
}
