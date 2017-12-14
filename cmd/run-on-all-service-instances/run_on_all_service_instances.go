// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/loggerfactory"
	"github.com/pivotal-cf/on-demand-service-broker/runner"
	"gopkg.in/yaml.v2"
)

func main() {
	loggerFactory := loggerfactory.New(os.Stdout, "run-on-all-service-instances", loggerfactory.Flags)
	logger := loggerFactory.New()

	configFilePath := flag.String("configFilePath", "", "path to config file")
	flag.Parse()

	rawConfig, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		logger.Fatalf("Error reading config file: %s", err)
	}

	var config runner.Config
	err = yaml.Unmarshal(rawConfig, &config)
	if err != nil {
		logger.Fatalf("Invalid config file: %s", err)
	}

	err = yaml.Unmarshal(configContents, &conf)
	if err != nil {
		logger.Fatalln(err.Error())
	}

	builder, err := runner.NewBuilder(conf, logger)
	if err != nil {
		logger.Fatalln(err.Error())
	}
	upgradeTool := upgrader.New(builder)

	err = upgradeTool.Upgrade()
	if err != nil {
		logger.Fatalln(err.Error())
	}
}
