---
layout: "scvmm"
page_title: "Provider: SCVMM"
sidebar_current: "docs-scvmm-index"
description: |-
  The SCVMM provider is used to interact with the resources supported by
  Microsoft SCVMM. The provider needs to be configured with the proper credentials
  before it can be used.
---

# SCVMM Provider

The SCVMM provider is used to interact with the resources supported by
SCVMM.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

~> **NOTE:** The SCVMM Provider currently represents _initial support_
and therefore may undergo significant changes as the community improves it. This
provider at this time only supports adding Virtual Machine, Virtual Disk and Checkpoint Resource

## Example Usage

```hcl
# Configure the Microsoft SCVMM Provider
provider "scvmm" {
  server_ip = "${var.scvmm_server_ip}"
  port = ${var.scvmm_server_port}
  user_name = "${var.scvmm_server_user}"
  user_password = "${var.scvmm_server_password}"
}

# Add a Virtual Machine
resource "scvmm_virtual_machine" "CreateVM" {
  timeout = "10000"
  vmm_server = "WIN-2F929KU8HIU"
  vm_name = "Test_VM_demo01"
  template_name = "TestVMTemplate"
  cloud_name = "GSL Cloud"
}

#Add a Virtual Disk
resource "scvmm_virtual_disk" "demoDisk" {
  timeout = 10000
  vmm_server = "${scvmm_virtual_machine.CreateVM.vmm_server}"
  vm_name = "${scvmm_virtual_machine.CreateVM.vm_name}"
  virtual_disk_name = "demo_disk"
  virtual_disk_size = 10000
}

resource "scvmm_checkpoint" "demoCheck" {
        timeout=1000
        vmm_server="${scvmm_virtual_disk.demoDisk.vmm_server}"
        vm_name="${scvmm_virtual_machine.CreateVM.vm_name}"
        checkpoint_name="demo_checkpoint"
}
```

## Argument Reference

The following arguments are used to configure the SCVMM Provider:

* `user_name` - (Required) This is the username for SCVMM server login. Can also
  be specified with the `SCVMM_SERVER_USER` environment variable.
* `user_password` - (Required) This is the password for Open Daylight API operations. Can
  also be specified with the `SCVMM_SERVER_PASSWORD` environment variable.
* `server_ip` - (Required) This is the SCVMM server ip for SCVMM
  operations. Can also be specified with the `SCVMM_SERVER_IP` environment
  variable.
* `port` - (Required) This is the port for winrm operations of the SCVMM Server.

## Acceptance Tests

The Active Directory provider's acceptance tests require the above provider
configuration fields to be set using the documented environment variables.

Once all these variables are in place, the tests can be run like this:

```
make testacc TEST=./builtin/providers/scvmm
```
