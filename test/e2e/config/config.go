/*
Copyright 2022 The cert-manager Authors.

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

package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	kubeConfig string

	RestConfig *rest.Config
}

func New(fs *flag.FlagSet) *Config {
	return new(Config).addFlags(fs)
}

func (c *Config) Complete() error {
	if c.kubeConfig == "" {
		return fmt.Errorf("--kubeconfig-path must be specified")
	}

	var err error
	c.RestConfig, err = clientcmd.BuildConfigFromFlags("", c.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to build kubernetes rest config: %s", err)
	}

	return nil
}

func (c *Config) addFlags(fs *flag.FlagSet) *Config {
	kubeConfigFile := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if kubeConfigFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Failed to get user home directory: " + err.Error())
		}
		kubeConfigFile = filepath.Join(homeDir, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName)
	}

	fs.StringVar(&c.kubeConfig, "kubeconfig-path", kubeConfigFile, "Path to config containing embedded authinfo for kubernetes. Default value is from environment variable "+clientcmd.RecommendedConfigPathEnvVar)
	return c
}
