The goal is to add the following parameters incl. there respectiv functions to the packer cloudstack plugin

boot_wait (duration string | ex: "1h5m2s")

boot_command ([]string)

websocket_url (string)

websocket_port (int)


With a newly integrated cloudstack feature (https://www.shapeblue.com/api-driven-console-access/) it should be possible for third party software to get access to the VNC console via the data of response of api call createConsoleEndpoint.
Here is a response:
{"createconsoleendpointresponse":{"consoleendpoint":{"success":true,"url":"http://console.proio.cloud/resource/noVNC/vnc.html?autoconnect=true&port=8080&token=Eyp00oc+FQmQECo/QvnDP0fO5cXzRD668vIdVcuP3GFT0iqHk87y/egYf8PYYjXDHGdjS/79IQvO+FJpqlSWpqvXhqgmKYSf3yamioIf9XCJdN9NjFPcElHPhptPE+YV2ifhCmbRFO0LrVXFgC0IN/ryxR6131fJXPX4G9YRCy+Ms6v/bM+RZRWveaw4wIp1swOwhuerSGAMy65vV1dttn4KkJWprOMvBAynoNB3eQDrJfYuR8osM2J+fdaceGCPqV2Q9j62oVGuJnOBQl5jJUSL/oYseMuB18Gk+wLe6Liun9AUNnVJASftnmncSZohCo59/+grsU12l/4ikrN7cqlQ7r2VGy/tf3o895sldapIN5O1XKiJSoMJDQbjUC7v/cW3T3k8thBdAF9dywd6FZyO6930+LzLbipwzBP00fXJu/dRruC1MKVmLt7pnD+oQXzO7qg3kaFW3c2SU2wT9LpW+pC2sQIjPSVUmwN+Nc/Tv+vUOh4JCfa8gfH55crUf0krNisRxn/J5DXpHOuyRYMXKda6OiA=","websocket":{"token":"Eyp00oc+FQmQECo/QvnDP0fO5cXzRD668vIdVcuP3GFT0iqHk87y/egYf8PYYjXDHGdjS/79IQvO+FJpqlSWpqvXhqgmKYSf3yamioIf9XCJdN9NjFPcElHPhptPE+YV2ifhCmbRFO0LrVXFgC0IN/ryxR6131fJXPX4G9YRCy+Ms6v/bM+RZRWveaw4wIp1swOwhuerSGAMy65vV1dttn4KkJWprOMvBAynoNB3eQDrJfYuR8osM2J+fdaceGCPqV2Q9j62oVGuJnOBQl5jJUSL/oYseMuB18Gk+wLe6Liun9AUNnVJASftnmncSZohCo59/+grsU12l/4ikrN7cqlQ7r2VGy/tf3o895sldapIN5O1XKiJSoMJDQbjUC7v/cW3T3k8thBdAF9dywd6FZyO6930+LzLbipwzBP00fXJu/dRruC1MKVmLt7pnD+oQXzO7qg3kaFW3c2SU2wT9LpW+pC2sQIjPSVUmwN+Nc/Tv+vUOh4JCfa8gfH55crUf0krNisRxn/J5DXpHOuyRYMXKda6OiA=","host":"console.proio.cloud","port":"8080","path":"websockify"}}}}

There are 2 different ways to access the console. First way is via https link provided by the response. Second way is to open a websocket connection via data provided by the response. Prefered (and maybe only way) is to use the websocket connection to send keyboard commands to the console.

The packer-plugin-vmware (https://github.com/hashicorp/packer-plugin-vmware/blob/main/builder/vmware/common/step_vnc_connect.go#L54) uses a websocket connection to send keyboard commands (connectOverWebsocket). Mybe we use some code of this.

The packer qemu plugin also uses some kind of way to send keyboard commands to the VNC console. As far as I understand you need a direct connection to noVNC console to use this. This direct way is not available in Cloudstack. But maybe we can use the code from the packer qemu plugin as a template for the code for the packer cloudstack plugin and extend or rewrite some code to get it work using websocket connection to send keyboard commands. You will find a list of boot_commands here: https://developer.hashicorp.com/packer/integrations/hashicorp/qemu/latest/components/builder/qemu#boot-configuration


Parameter explanation

boot_wait (duration string | ex: "1h5m2s")

The time to wait after booting the initial virtual machine before typing the boot_command. The value of this should be a duration. Examples are 5s and 1m30s which will cause Packer to wait five seconds and one minute 30 seconds, respectively. If this isn't specified, the default is 10s or 10 seconds. To set boot_wait to 0s, use a negative number, such as "-1s"

boot_command ([]string)

This is an array of commands to send via websocket connection when the virtual machine is first booted. The goal of these commands should be to type just enough to initialize the operating system installer. Special keys can be typed as well, and are covered in the section below on the boot command. If this is not specified, it is assumed the installer will start itself.

websocket_url (string)

With this parameter you can override the default value of the url the packer cloudstack plugin should use to establish the websocket connection. Default is the host value from api response of createConsoleEndpoint.

websocket_port (int)

With this parameter you can override the default value of the port the packer cloudstack plugin should use to establish the websocket connection. Default is the port value from api response of createConsoleEndpoint.



More info
https://github.com/sbrueseke/packer-plugin-cloudstack
https://developer.hashicorp.com/packer/integrations/hashicorp/qemu/latest/components/builder/qemu
https://developer.hashicorp.com/packer/integrations/hashicorp/qemu/latest/components/builder/qemu#boot-configuration
https://developer.hashicorp.com/packer/integrations/hashicorp/cloudstack/latest/components/builder/cloudstack
https://www.shapeblue.com/api-driven-console-access/
https://cloudstack.apache.org/api/apidocs-4.18/apis/createConsoleEndpoint.html



Cloudstack API call and response example
https://portal.proio.cloud/client/api/?command=createConsoleEndpoint&virtualmachineid=2762aa77-ebcd-448d-8236-ec9a9031bd9a&response=json

{"createconsoleendpointresponse":{"consoleendpoint":{"success":true,"url":"http://console.proio.cloud/resource/noVNC/vnc.html?autoconnect=true&port=8080&token=Eyp00oc+FQmQECo/QvnDP0fO5cXzRD668vIdVcuP3GFT0iqHk87y/egYf8PYYjXDHGdjS/79IQvO+FJpqlSWpqvXhqgmKYSf3yamioIf9XCJdN9NjFPcElHPhptPE+YV2ifhCmbRFO0LrVXFgC0IN/ryxR6131fJXPX4G9YRCy+Ms6v/bM+RZRWveaw4wIp1swOwhuerSGAMy65vV1dttn4KkJWprOMvBAynoNB3eQDrJfYuR8osM2J+fdaceGCPqV2Q9j62oVGuJnOBQl5jJUSL/oYseMuB18Gk+wLe6Liun9AUNnVJASftnmncSZohCo59/+grsU12l/4ikrN7cqlQ7r2VGy/tf3o895sldapIN5O1XKiJSoMJDQbjUC7v/cW3T3k8thBdAF9dywd6FZyO6930+LzLbipwzBP00fXJu/dRruC1MKVmLt7pnD+oQXzO7qg3kaFW3c2SU2wT9LpW+pC2sQIjPSVUmwN+Nc/Tv+vUOh4JCfa8gfH55crUf0krNisRxn/J5DXpHOuyRYMXKda6OiA=","websocket":{"token":"Eyp00oc+FQmQECo/QvnDP0fO5cXzRD668vIdVcuP3GFT0iqHk87y/egYf8PYYjXDHGdjS/79IQvO+FJpqlSWpqvXhqgmKYSf3yamioIf9XCJdN9NjFPcElHPhptPE+YV2ifhCmbRFO0LrVXFgC0IN/ryxR6131fJXPX4G9YRCy+Ms6v/bM+RZRWveaw4wIp1swOwhuerSGAMy65vV1dttn4KkJWprOMvBAynoNB3eQDrJfYuR8osM2J+fdaceGCPqV2Q9j62oVGuJnOBQl5jJUSL/oYseMuB18Gk+wLe6Liun9AUNnVJASftnmncSZohCo59/+grsU12l/4ikrN7cqlQ7r2VGy/tf3o895sldapIN5O1XKiJSoMJDQbjUC7v/cW3T3k8thBdAF9dywd6FZyO6930+LzLbipwzBP00fXJu/dRruC1MKVmLt7pnD+oQXzO7qg3kaFW3c2SU2wT9LpW+pC2sQIjPSVUmwN+Nc/Tv+vUOh4JCfa8gfH55crUf0krNisRxn/J5DXpHOuyRYMXKda6OiA=","host":"console.proio.cloud","port":"8080","path":"websockify"}}}}
