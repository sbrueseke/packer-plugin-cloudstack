package cloudstack

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type stepBootCommandVNC struct {
}

func (s stepBootCommandVNC) Run(ctx context.Context, bag multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s stepBootCommandVNC) Cleanup(bag multistep.StateBag) {
	//TODO implement me
}
