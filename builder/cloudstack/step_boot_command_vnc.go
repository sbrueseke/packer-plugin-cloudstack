package cloudstack

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/bootcommand"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/sbrueseke/go-vnc"
	"time"
)

type stepBootCommandVNC struct {
	VNCEnabled bool
	Config     bootcommand.VNCConfig
	BootWait   time.Duration
	Ctx        interpolate.Context
}

func (s stepBootCommandVNC) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if !s.VNCEnabled {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packersdk.Ui)
	conn := state.Get("vnc_conn").(*vnc.ClientConn)
	defer conn.Close()

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

	d := bootcommand.NewVNCDriver(conn, s.Config.BootKeyInterval)

	ui.Say("Typing the boot command over VNC...")
	flatBootCommand := s.Config.FlatBootCommand()
	command, err := interpolate.Render(flatBootCommand, &s.Ctx)
	if err != nil {
		err := fmt.Errorf("Error preparing boot command: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	seq, err := bootcommand.GenerateExpressionSequence(command)
	if err != nil {
		err := fmt.Errorf("Error generating boot command: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if err := seq.Do(ctx, d); err != nil {
		err := fmt.Errorf("Error running boot command: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s stepBootCommandVNC) Cleanup(bag multistep.StateBag) {
}
