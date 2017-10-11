---
layout: "scvmm"
page_title: "Microsoft SCVMM: scvmm_checkpoint"
sidebar_current: "docs-scvmm-resource-inventory-folder"
description: |-
  Provides a SCVMM Checkpoint resource. This can be used to create and delete Checkpoint in Virtual Machine existing on SCVMM Server.
---

# scvvm\_checkpoint

Provides a SCVMM Checkpoint resource. This can be used to create and delete Checkpoint in Virtual Machine existing on SCVMM Server. It adds a checkpoint to virtual machine using a powershell script. 

## Example Usage

```hcl
resource "scvmm_checkpoint" "demoCheck" {
        timeout=1000
        vmm_server="${scvmm_virtual_disk.demoDisk.vmm_server}"
        vm_name="${scvmm_virtual_machine.CreateVM.vm_name}"
        checkpoint_name="demo_checkpoint"
}
```

## Argument Reference

The following arguments are supported:

* `timeout` - (Required) The time within which the specified resource should be created
* `vmm_server` - (Required) The Virtual Machine Manager of the Virtual Machine
* `vm_name` - (Required) The Virtual Machine Name to which the Virtual Disk will be attached
* `checkpoint_name` - (Required) The Checkpoint name