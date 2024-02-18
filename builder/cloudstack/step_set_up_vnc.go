package cloudstack

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type stepSetUpVNC struct {
}

func (s stepSetUpVNC) Run(ctx context.Context, bag multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s stepSetUpVNC) Cleanup(bag multistep.StateBag) {
	//TODO implement me
}
