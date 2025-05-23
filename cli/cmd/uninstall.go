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
	"log"
	"strings"

	"drasi.io/cli/installers"
	"drasi.io/cli/sdk/registry"

	"github.com/spf13/cobra"
)

func NewUninstallCommand() *cobra.Command {
	var uninstallCommand = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Drasi",
		Long: `Uninstall the Drasi environment from the the default or a specific namespace on the current Kubernetes cluster.
		
Usage examples:
  drasi uninstall
  drasi uninstall -n my-namespace
`,
		Args: cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var currentNamespace string
			if currentNamespace, err = cmd.Flags().GetString("namespace"); err != nil {
				return err
			}

			fmt.Println("Uninstalling Drasi")
			fmt.Println("Deleting namespace: ", currentNamespace)

			// Ask for confirmation if the user didn't pass the -y flag
			if !cmd.Flags().Changed("yes") {
				fmt.Printf("Are you sure you want to uninstall Drasi from the namespace %s? (yes/no): ", currentNamespace)
				if !askForConfirmation(currentNamespace) {
					fmt.Println("Uninstall cancelled")
					return nil
				}
			}

			uninstallDapr, _ := cmd.Flags().GetBool("uninstall-dapr")

			reg, err := registry.LoadCurrentRegistrationWithNamespace(currentNamespace)
			if err != nil {
				return err
			}

			uninstaller, err := installers.MakeUninstaller(reg)
			if err != nil {
				return err
			}

			err = uninstaller.Uninstall(uninstallDapr)
			if err != nil {
				return err
			}

			fmt.Println("Drasi uninstalled successfully")

			return nil
		},
	}

	uninstallCommand.Flags().BoolP("yes", "y", false, "Automatic yes to prompts")
	uninstallCommand.Flags().BoolP("uninstall-dapr", "d", false, "Uninstall Dapr by deleting the Dapr system namespace")

	return uninstallCommand
}

func askForConfirmation(namespace string) bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return askForConfirmation(namespace)
	}
}
