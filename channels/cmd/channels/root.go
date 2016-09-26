package main

import (
	goflag "flag"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

type RootCmd struct {
	configFile string

	cobraCommand *cobra.Command
}

var rootCommand = RootCmd{
	cobraCommand: &cobra.Command{
		Use:   "channels",
		Short: "channels applies software from a channel",
	},
}

func Execute() {
	goflag.Set("logtostderr", "true")
	goflag.CommandLine.Parse([]string{})
	if err := rootCommand.cobraCommand.Execute(); err != nil {
		exitWithError(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	cmd := rootCommand.cobraCommand

	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)

	cmd.PersistentFlags().StringVar(&rootCommand.configFile, "config", "", "config file (default is $HOME/.channels.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootCommand.configFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(rootCommand.configFile)
	}

	viper.SetConfigName(".channels") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func (c *RootCmd) AddCommand(cmd *cobra.Command) {
	c.cobraCommand.AddCommand(cmd)
}

func (c *RootCmd) KubernetesClient() (*release_1_3.Clientset, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot load kubecfg settings: %v", err)
	}

	k8sClient, err := release_1_3.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot build kube client: %v", err)
	}
	return k8sClient, err
}