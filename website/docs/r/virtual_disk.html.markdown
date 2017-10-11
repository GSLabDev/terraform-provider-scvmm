---
layout: "scvmm"
page_title: "Microsoft SCVMM: scvmm_virtual_disk"
sidebar_current: "docs-scvmm-resource-inventory-folder"
description: |-
  Provides a SCVMM Virtual Disk resource. This can be used to create and delete Virtual Disk in SCVMM Server.
---

# scvvm\_virtual\_disk

Provides a SCVMM Virtual Disk resource. This can be used to create and delete Virtual Disk in Virtual Machine existing in SCVMM. It adds a disk drive using a powershell script to the SCSI Adapter. 

## Example Usage

```hcl
resource "scvmm_virtual_disk" "demoDisk" {
  timeout = 10000
  vmm_server = "${scvmm_virtual_machine.CreateVM.vmm_server}"
  vm_name = "${scvmm_virtual_machine.CreateVM.vm_name}"
  virtual_disk_name = "demo_disk"
  virtual_disk_size = 10000
}
```

## Argument Reference

The following arguments are supported:

* `timeout` - (Required) The time within which the specified resource should be created
* `vmm_server` - (Required) The Virtual Machine Manager of the Virtual Machine
* `vm_name` - (Required) The Virtual Machine Name to which the Virtual Disk will be attached
* `virtual_disk_name` - (Required) The Virtual Disk name
* `virtual_disk_size` - (Required) The size of the Virtual Disk in Mega Bytes