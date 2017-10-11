---
layout: "scvmm"
page_title: "Microsoft SCVMM: scvmm_virtual_machine"
sidebar_current: "docs-scvmm-resource-inventory-folder"
description: |-
  Provides a SCVMM Virtual Machine resource. This can be used to create and delete Virtual Machine in SCVMM Server.
---

# scvvm\_virtual\_machine

Provides a SCVMM Virtual Machine resource. This can be used to create and delete Virtual Machine from SCVMM. It does this by running script on remote machine which makes use of vmm server, vm name, cloud name and template name as parameters.

## Example Usage

```hcl
resource "scvmm_virtual_machine" "CreateVM" {
  timeout = "10000"
  vmm_server = "WIN-2F929KU8HIU"
  vm_name = "Test_VM_demo01"
  template_name = "TestVMTemplate"
  cloud_name = "GSL Cloud"
}
```

## Argument Reference

The following arguments are supported:

* `timeout` - (Required) The time within which the specified resource should be created
* `vmm_server` - (Required) The Virtual Machine Manager of the Virtual Machine
* `vm_name` - (Required) The Virtual Machine Name
* `template_name` - (Required) The template name containg all the configuration for Virtual Machine
* `cloud_name` - (Required) The cloud in which the Virtual Machine should be created