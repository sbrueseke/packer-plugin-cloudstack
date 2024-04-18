package cloudstack

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/sbrueseke/go-vnc"
	"golang.org/x/net/websocket"
	"log"
	"net/url"
)

type stepSetUpVNC struct {
	VNCEnabled         bool
	WebsocketURL       string
	InsecureConnection bool
}

func (s stepSetUpVNC) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if !s.VNCEnabled {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packersdk.Ui)
	ui.Say("Setting up VNC...")

	var wsURL string

	if s.WebsocketURL != "" {
		wsURL = s.WebsocketURL
	} else {
		consoleEndpointURL, err := setUpWithCreateConsoleEndpoint(state)
		if err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		wsURL = consoleEndpointURL
	}

	//connect to websocket

	log.Printf("[DEBUG] websocket url: %s", wsURL)
	console, err := url.Parse(wsURL)
	if err != nil {
		state.Put("error", fmt.Errorf("Error parsing websocket url: %s\n", err))
		return multistep.ActionHalt
	}
	ui.Say(wsURL)
	origin, err := url.Parse("http://localhost")
	if err != nil {
		state.Put("error", fmt.Errorf("Error parsing websocket origin url: %s\n", err))
		return multistep.ActionHalt
	}

	// Create the websocket connection and set it to a BinaryFrame
	websocketConfig := &websocket.Config{
		Location:  console,
		Origin:    origin,
		TlsConfig: &tls.Config{InsecureSkipVerify: s.InsecureConnection},
		Version:   websocket.ProtocolVersionHybi13,
		Protocol:  []string{"binary"},
	}
	nc, err := websocket.DialConfig(websocketConfig)
	if err != nil {
		state.Put("error", fmt.Errorf("Error Dialing: %s\n", err))
		return multistep.ActionHalt
	}
	nc.PayloadType = websocket.BinaryFrame

	// Setup the VNC connection over the websocket
	ccconfig := &vnc.ClientConfig{
		Auth:      []vnc.ClientAuth{new(vnc.VencryptAuth)},
		Exclusive: false,
	}
	connection, err := vnc.Client(nc, ccconfig)
	if err != nil {
		state.Put("error", fmt.Errorf("Error setting the VNC over websocket client: %s\n", err))
		return multistep.ActionHalt
	}

	state.Put("vnc_conn", connection)
	return multistep.ActionContinue
}

func setUpWithCreateConsoleEndpoint(state multistep.StateBag) (string, error) {
	client := state.Get("client").(*cloudstack.CloudStackClient)

	virtualMachineId := state.Get("instance_id").(string)
	p := client.ConsoleEndpoint.NewCreateConsoleEndpointParams(virtualMachineId)

	endpoint, err := client.ConsoleEndpoint.CreateConsoleEndpoint(p)
	if err != nil {
		return "", fmt.Errorf("failed to create console endpoint: %s", err)
	}

	host := endpoint.Websocket["host"].(string)
	path := endpoint.Websocket["path"].(string)
	token := endpoint.Websocket["token"].(string)
	port := endpoint.Websocket["port"].(string)

	websocketUrl := fmt.Sprintf("wss://%s:%s/%s?token=%s", host, port, path, token)

	return websocketUrl, nil
}

func (s stepSetUpVNC) Cleanup(bag multistep.StateBag) {}
