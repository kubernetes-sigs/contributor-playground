/*
Copyright 2018 The Kubernetes Authors.

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
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	//_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	_ "k8s.io/kubernetes/pkg/features" // add the kubernetes feature gates
	utilflag "k8s.io/kubernetes/pkg/util/flag"
	//_ "k8s.io/kubernetes/pkg/version/prometheus" // for version metric registration

	"icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/cloud-controller-manager/app"
	"icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/cloud-controller-manager/app/options"
	cloud_provider "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/cloud-provider"
	_ "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/cloud-provider"
)

var version string

func init() {
	cloud_provider.CCMVersion = version
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	logs.InitLogs()
	defer logs.FlushLogs()

	s, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	command := &cobra.Command{
		Use:  "cce-cloud-controller-manager",
		Long: `cce-cloud-controller-manager manages BAIDU cloud resources for a Kubernetes cluster.`,
		Run: func(cmd *cobra.Command, args []string) {

			// Use our version instead of the Kubernetes formatted version
			versionFlag := cmd.Flags().Lookup("version")
			if versionFlag.Value.String() == "true" {
				fmt.Printf("%s version: %s\n", cmd.Name(), version)
				os.Exit(0)
			}

			// Hard code aws cloud provider
			cloudProviderFlag := cmd.Flags().Lookup("cloud-provider")
			cloudProviderFlag.Value.Set("cce")

			utilflag.PrintFlags(cmd.Flags())

			c, err := s.Config(app.KnownControllers(), app.ControllersDisabledByDefault.List())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := app.Run(c.Complete(), wait.NeverStop); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	fs := command.Flags()
	namedFlagSets := s.Flags(app.KnownControllers(), app.ControllersDisabledByDefault.List())
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), command.Name())

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func GetVersion() string {
	return version
}
