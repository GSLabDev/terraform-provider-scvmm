# Terraform Microsoft SCVMM Provider

This is the repository for the Terraform [Microsoft SCVMM][1] Provider, which one can use
with Terraform to work with Microsoft SCVMM.

[1]: https://docs.microsoft.com/en-us/system-center/vmm/

Coverage is currently only limited to Virtual Machine, Virtual Disk and Checkpoint but in the coming months we are planning release coverage for most essential Microsoft SCVMM workflows.
Watch this space!

For general information about Terraform, visit the [official website][3] and the
[GitHub project page][4].

[3]: https://terraform.io/
[4]: https://github.com/hashicorp/terraform

# Using the Provider

The current version of this provider requires Terraform v0.10.2 or higher to
run.

Note that you need to run `terraform init` to fetch the provider before
deploying. Read about the provider split and other changes to TF v0.10.0 in the
official release announcement found [here][4].

[4]: https://www.hashicorp.com/blog/hashicorp-terraform-0-10/

## Full Provider Documentation

The provider is usefull in adding Virtual Machine, Virtual Disk and Checkpoint using Microsoft SCVMM.

### Example
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

# Building The Provider

**NOTE:** Unless you are [developing][7] or require a pre-release bugfix or feature,
you will want to use the officially released version of the provider (see [the
section above][8]).

[7]: #developing-the-provider
[8]: #using-the-provider


## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/GSLabDev/terraform-provider-scvmm`:

```sh
mkdir -p $GOPATH/src/github.com/GSLabDev/
cd $GOPATH/src/github.com/GSLabDev/
git clone git@github.com:terraform-providers/terraform-provider-scvmm
```

## Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/src/github.com/GSLabDev/terraform-provider-scvmm
make
```

## Installing the Local Plugin

After the build is complete, copy the `terraform-provider-scvmm` binary into
the same path as your `terraform` binary, and re-run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Developing the Provider

If you wish to work on the provider, you'll first need [Go][9] installed on your
machine (version 1.9+ is **required**). You'll also need to correctly setup a
[GOPATH][10], as well as adding `$GOPATH/bin` to your `$PATH`.

[9]: https://golang.org/
[10]: http://golang.org/doc/code.html#GOPATH

See [Building the Provider][11] for details on building the provider.

[11]: #building-the-provider

# Testing the Provider

**NOTE:** Testing the Microsoft SCVMM provider is currently a complex operation as it
requires having a Microsoft SCVMM Server to test against.

## Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment
variables to run. See the individual `*_test.go` files in the
[`scvmm/`](scvmm/) directory for more details. The next section also
describes how you can manage a configuration file of the test environment
variables.

### Using the `.tf-scvmm-devrc.mk` file

The [`tf-scvmm-devrc.mk.example`](tf-scvmm-devrc.mk.example) file contains
an up-to-date list of environment variables required to run the acceptance
tests. Copy this to `$HOME/.tf-scvmm-devrc.mk` and change the permissions to
something more secure (ie: `chmod 600 $HOME/.tf-scvmm-devrc.mk`), and
configure the variables accordingly.

## Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccSCVMM"
```

This following example would run all of the acceptance tests matching
`TestAccSCVMM`. Change this for the specific tests you want to
run.
