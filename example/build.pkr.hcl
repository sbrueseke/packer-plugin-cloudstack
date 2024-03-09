# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

packer {
#  required_plugins {
#    cloudstack = {
#      version = ">= 1.0.0"
#      source  = "github.com/hashicorp/cloudstack"
#    }
#  }
}

source "cloudstack" "ubuntu" {
  # cloud api login
  api_url    = "https://portal.proio.cloud/client/api"
  api_key    = "fPuqVQHLdJke6Sl8XttiaqZnIcXkSCPOp-2NN1F_aXgZYQaevMpFekv38Qew2akiXRnmWTAUXm9xZ6QMy5Hh8A"
  secret_key = "auQozjMnHarn7X_RazxHP94FOjpffmmqVXc-kRyoRMkyfUl4ZY_WiPWLHmVgJLZbUcrrFCy9JxJRwQ7A7HISPg"
  # vm settings
  zone             = "fra1"
  hypervisor       = "KVM"
  network          = "999-go-testing-network"
  source_iso       = "Ubuntu 22.04.3 LTS live server"
  service_offering = "packer"
  disk_offering    = "Custom"
  disk_size        = "9"
  # template settings
  ssh_username              = "ubuntu"
  ssh_password              = "dVvK8R4!nbf(DK2T"
  ssh_timeout               = "30m"
  template_name             = "packer Ubuntu 22.04 LTS {{isotime}}"
  template_display_text     = "user: ubuntu created: {{isotime}}"
  template_public           = false
  template_featured         = false
  template_os               = "Ubuntu 22.04 LTS"
  template_password_enabled = true
  template_scalable         = true
  # other
  expunge                   = true
  use_local_ip_address      = true
  # http_directory            = "httpdir/"
  # http_port_min             = 80
  # http_port_max             = 80
}

build {
  sources = ["source.cloudstack.ubuntu"]
}
