// Copyright 2021, The Fuchsia Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	subcommands "github.com/google/subcommands"
	fg "go.fuchsia.dev/fuchsia/tools/emu-control/emu-grpc"
)

var (
	server_addr = flag.String("server", "localhost", "server address")
	server_port = flag.Int("port", 5556, "gRPC port")
	timeout     = flag.Duration("timeout", 1*time.Second, "gRPC connection timeout")
)

func newClient() (fg.FemuGrpcClientInterface, error) {
	config := fg.FemuGrpcClientConfig{
		ServerAddr: fmt.Sprintf("%s:%d", *server_addr, *server_port),
		Timeout:    *timeout,
	}

	client, err := fg.NewFemuGrpcClient(config)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&recordAudioCmd{}, "")
	subcommands.Register(&recordScreenCmd{}, "")
	subcommands.Register(&keyboardCmd{}, "")
	flag.Parse()

	client, err := newClient()
	if err != nil {
		panic(fmt.Sprintf("error while creating gRPC client: %v", err))
	}
	ctx := context.WithValue(context.Background(), "client", client)
	os.Exit(int(subcommands.Execute(ctx)))
}
