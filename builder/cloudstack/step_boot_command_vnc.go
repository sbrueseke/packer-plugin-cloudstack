package cloudstack

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"time"
)

type stepBootCommandVNC struct {
	VNCEnabled bool
	BootWait   time.Duration
}

func (s stepBootCommandVNC) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if !s.VNCEnabled {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packersdk.Ui)

	// Wait the for the vm to boot.
	if int64(s.BootWait) > 0 {
		ui.Say(fmt.Sprintf("Waiting %s for boot...", s.BootWait.String()))
		select {
		case <-time.After(s.BootWait):
			break
		case <-ctx.Done():
			return multistep.ActionHalt
		}
	}
	return multistep.ActionContinue
}

func (s stepBootCommandVNC) Cleanup(bag multistep.StateBag) {
	//TODO implement me
}
