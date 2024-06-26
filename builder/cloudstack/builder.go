// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudstack

import (
	"context"
	"fmt"
	"github.com/apache/cloudstack-go/v2/cloudstack"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

const BuilderId = "packer.cloudstack"

// Builder represents the CloudStack builder.
type Builder struct {
	config Config
	runner multistep.Runner
	ui     packersdk.Ui
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	errs := b.config.Prepare(raws...)
	if errs != nil {
		return nil, nil, errs
	}

	return nil, nil, nil
}

// Run implements the packersdk.Builder interface.
func (b *Builder) Run(ctx context.Context, ui packersdk.Ui, hook packersdk.Hook) (packersdk.Artifact, error) {
	b.ui = ui

	// Create a CloudStack API client.
	client := cloudstack.NewAsyncClient(
		b.config.APIURL,
		b.config.APIKey,
		b.config.SecretKey,
		!b.config.SSLNoVerify,
	)

	// Set the time to wait before timing out
	client.AsyncTimeout(int64(b.config.AsyncTimeout.Seconds()))

	// Some CloudStack service providers only allow HTTP GET calls.
	client.HTTPGETOnly = b.config.HTTPGetOnly

	// Set up the state.
	state := new(multistep.BasicStateBag)
	state.Put("client", client)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Build the steps.
	steps := []multistep.Step{
		&stepPrepareConfig{},
		commonsteps.HTTPServerFromHTTPConfig(&b.config.HTTPConfig),
		&stepKeypair{
			Debug:        b.config.PackerDebug,
			Comm:         &b.config.Comm,
			DebugKeyPath: fmt.Sprintf("cs_%s.pem", b.config.PackerBuildName),
		},
		&stepCreateSecurityGroup{},
		&stepCreateInstance{
			Ctx:   b.config.ctx,
			Debug: b.config.PackerDebug,
		},
		&stepSetupNetworking{},
		&stepDetachIso{},
		&stepSetUpVNC{
			VNCEnabled:         !b.config.DisableVNC,
			WebsocketURL:       b.config.WebsocketURL,
			InsecureConnection: b.config.InsecureConnection,
		},
		&stepBootCommandVNC{
			VNCEnabled: !b.config.DisableVNC,
			Config:     b.config.VNCConfig,
			BootWait:   b.config.BootWait,
			Ctx:        b.config.ctx,
		},
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      communicator.CommHost(b.config.Comm.Host(), "ipaddress"),
			SSHConfig: b.config.Comm.SSHConfigFunc(),
			SSHPort:   commPort,
			WinRMPort: commPort,
		},
		&commonsteps.StepProvision{},
		&commonsteps.StepCleanupTempKeys{
			Comm: &b.config.Comm,
		},
		&stepShutdownInstance{},
		&stepCreateTemplate{},
	}

	// Configure the runner and run the steps.
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		ui.Error(rawErr.(error).Error())
		return nil, rawErr.(error)
	}

	// If there was no template created, just return
	if _, ok := state.GetOk("template"); !ok {
		return nil, nil
	}

	// Build the artifact and return it
	artifact := &Artifact{
		client:    client,
		config:    &b.config,
		template:  state.Get("template").(*cloudstack.CreateTemplateResponse),
		StateData: map[string]interface{}{"generated_data": state.Get("generated_data")},
	}

	return artifact, nil
}
