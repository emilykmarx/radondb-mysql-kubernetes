/*
Copyright 2021 RadonDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/radondb/radondb-mysql-kubernetes/sidecar"
	"github.com/radondb/radondb-mysql-kubernetes/utils"
)

const (
	// The name of the sidecar.
	sidecarName = "sidecar"
	// The short description of the sidecar.
	sidecarShort = "A simple helper for mysql operator."
)

var (
	log = logf.Log.WithName("sidecar")
	// A command for sidecar.
	cmd = &cobra.Command{
		Use:   sidecarName,
		Short: sidecarShort,
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("run the sidecar, see help section")
			os.Exit(1)
		},
	}
)

func main() {
	// setup logging
	logf.SetLogger(zap.New(zap.UseDevMode(true)))
	stop := make(chan struct{})

	containerName := sidecar.GetContainerType()

	if containerName == utils.ContainerBackupName {
		backupCfg := sidecar.NewBackupConfig()
		httpCmd := &cobra.Command{
			Use:   "http",
			Short: "start http server",
			Run: func(cmd *cobra.Command, args []string) {
				if err := sidecar.RunHttpServer(backupCfg, stop); err != nil {
					log.Error(err, "run command failed")
					os.Exit(1)
				}
			},
		}
		cmd.AddCommand(httpCmd)
	} else if containerName == utils.ContainerBackupJobName {
		reqBackupCfg := sidecar.NewReqBackupConfig()
		reqBackupCmd := &cobra.Command{
			Use:   "request_a_backup",
			Short: "start request a backup",
			Args: func(cmd *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("require one arguments. ")
				}
				return nil
			},
			Run: func(cmd *cobra.Command, args []string) {
				if err := sidecar.RunRequestBackup(reqBackupCfg, args[0]); err != nil {
					log.Error(err, "run command failed")
					os.Exit(1)
				}
			},
		}
		cmd.AddCommand(reqBackupCmd)
	} else {
		initCfg := sidecar.NewInitConfig()
		initCmd := sidecar.NewInitCommand(initCfg)
		cmd.AddCommand(initCmd)
	}

	if err := cmd.Execute(); err != nil {
		log.Error(err, "failed to execute command", "cmd", cmd)
		os.Exit(1)
	}
}
