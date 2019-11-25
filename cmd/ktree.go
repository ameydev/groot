/*
Copyright © 2019 Amey Deshmukh

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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"os/user"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// var cfgFile string

var indentationCount int = 1
var indentation string

// func initindentationCount() {
// 	indentationCount = 1
// }

type configs struct {
	Namespace  string
	KubeConfig string
}

func Execute() error {
	return ktree().Execute()
}

func ktree() *cobra.Command {

	c := &configs{
		Namespace:  "default",
		KubeConfig: getKubeConfig(),
	}

	cmd := &cobra.Command{
		Use:   "ktree",
		Short: "ktree is a k8s helper CLI utility tool.",
		Long: `This tool is used to find k8s resourses and their mappings with other k8s reources. For example:
		  
		ktree -n dev.`,
		PreRunE: func(cobracmd *cobra.Command, _ []string) error {
			// load current kube-config
			return initConfig(c)

		},

		RunE: func(_ *cobra.Command, _ []string) error {
			return getOverView(c)
		},
	}
	flags := cmd.Flags()

	flags.StringVar(&c.Namespace, "namespace", c.Namespace, "namespace in which we need to map k8s resources..")
	flags.StringVar(&c.KubeConfig, "kubeconfig", c.KubeConfig, "Any external kube config we want to use")

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig(c *configs) error {
	if c.KubeConfig != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c.KubeConfig)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ktree" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ktree")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	return nil
}

func getOverView(c *configs) error {
	// var contextName string
	// initIndentation()
	config, err := clientcmd.BuildConfigFromFlags("", c.KubeConfig)
	if err != nil {
		// log.Info("There was an error getting the config from c.KubeConfig.")
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// log.Info("There was an error getting clientset from config")
		return err
	}

	// contextName = clientcmd.getContextName()
	fmt.Println("Namespace : " + c.Namespace + " \n ")
	pods, err := clientset.CoreV1().Pods(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		// log.Info("There was an error getting the pod from clientset", err)
		return err
	}

	printPodDetails(pods)

	return nil
}

func getKubeConfig() string {
	var kubeconfig string

	if envVar := os.Getenv("KUBECONFIG"); len(envVar) > 0 {
		kubeconfig = envVar
	} else {
		usr, err := user.Current()
		if err != nil {
			fmt.Println(err)
		}
		kubeconfig = usr.HomeDir + "/.kube/config"
	}
	return kubeconfig
}

func getIndentation(ind int) string {

	indentation = ""
	for i := 0; i < ind; i++ {
		indentation += "\t"
	}
	indentationCount += 1
	return indentation

}
func printPodDetails(pods *v1.PodList) {
	// fmt.Printf("\nPods\n----\n")
	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	// fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "NAME", "READY", "STATUS", "RESTARTS")

	for _, pod := range pods.Items {
		fmt.Println(getIndentation(indentationCount)+"- Pods - "+pod.Name, "", pod.Status.Phase, "")
	}
	// w.Flush()
}
