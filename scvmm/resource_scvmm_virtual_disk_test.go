package scvmm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/masterzen/winrm"
)

func testBasicPreCheckVolume(t *testing.T) {
	if v := os.Getenv("SCVMM_VMM_SERVER"); v == "" {
		t.Fatal("SCVMM_VMM_SERVER must be set for acceptance tests")
	}
	if v := os.Getenv("SCVMM_TEMPLATE_NAME"); v == "" {
		t.Fatal("SCVMM_TEMPLATE_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("SCVMM_CLOUD_NAME"); v == "" {
		t.Fatal("SCVMM_CLOUD_NAME must be set for acceptance tests")
	}
}

func TestAccVolume_Basic(t *testing.T) {
	resourceName := "scvmm_virtual_disk.CreateDiskDrive"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testBasicPreCheckVolume(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy(resourceName),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVDConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "vm_name", "terraformVM"),
					resource.TestCheckResourceAttr(
						resourceName, "vmm_server", os.Getenv("SCVMM_VMM_SERVER")),
					resource.TestCheckResourceAttr(
						resourceName, "virtual_disk_name", "terraformDisk"),
				),
			},
		},
	})
}

func testAccCheckDiskDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		org := testAccProvider.Meta().(*winrm.Client)

		script := "\n[CmdletBinding(SupportsShouldProcess=$true)]\nparam (\n \n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VM Name\")]\n  [string]$vmName,\n\n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VmmServer\")]\n  [string]$vmmServer,\n\n   [Parameter(Mandatory=$true,HelpMessage=\"Enter Volume Name\")]\n  [string]$diskdriveName\n)\n\nBegin\n{  \n           If (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] \"Administrator\"))\n          {   \n                $arguments = \"\" + $myinvocation.mycommand.definition + \" \"\n                $myinvocation.BoundParameters.Values | foreach{$arguments += \"'$_' \" }\n            echo $arguments\n            Start-Process powershell -Verb runAs -ArgumentList $arguments\n            Break\n         }\n\u0009    try\n\u0009     {     if($vmName -eq $null) \n               {\n                    echo \"VM Name not entered\"\n                    exit\n                } \n                #gets virtual machine objects from the Virtual Machine Manager database\n                Set-SCVMMServer -VMMServer $vmmServer\n\u0009\u0009$VM = Get-SCVirtualMachine | Where-Object {$_.Name -Eq $vmName }   \n                #check if VM Exists\n               \n                    \n                          if($diskdriveName -ne $null)\n                          {\n                             #gets the specified volume object\n                             $diskdrive=Get-SCVirtualDiskDrive -VM $VM | Where-Object { $_.VirtualHardDisk.Name -Eq $diskdriveName} \n                             if($diskdrive -eq $null)\n                             {\n                               #This will delete specified volume attached which is to that vm \n                               Write-Error \"Virtual Disk Not Found\"\n                             }\n\n                          }  \n                          else\n                           {\n                                echo \"Name of the disk drive  is not entered\"\n                           }                             \n         }\n\u0009     catch [Exception]\n\u0009       {\n\u0009\u0009        echo $_.Exception.Message\n\u0009        }\n}\n"

		arguments := rs.Primary.Attributes["vm_name"] + " " + rs.Primary.Attributes["vmm_server"] + " " + rs.Primary.Attributes["virtual_disk_name"]
		filename := "deletediskdrive"
		result, err := execScript(org, script, filename, arguments)

		if err == "" {
			return fmt.Errorf("Error, Volume is still exist  %v", result)
		}
		return nil
	}
}

func testAccCheckDiskExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vm ID is set")
		}
		org := testAccProvider.Meta().(*winrm.Client)

		script := "\n[CmdletBinding(SupportsShouldProcess=$true)]\nparam (\n \n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VM Name\")]\n  [string]$vmName,\n\n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VmmServer\")]\n  [string]$vmmServer,\n\n   [Parameter(Mandatory=$true,HelpMessage=\"Enter Volume Name\")]\n  [string]$diskdriveName\n)\n\nBegin\n{  \n           If (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] \"Administrator\"))\n          {   \n                $arguments = \"\" + $myinvocation.mycommand.definition + \" \"\n                $myinvocation.BoundParameters.Values | foreach{$arguments += \"'$_' \" }\n            echo $arguments\n            Start-Process powershell -Verb runAs -ArgumentList $arguments\n            Break\n         }\n\u0009    try\n\u0009     {     if($vmName -eq $null) \n               {\n                    echo \"VM Name not entered\"\n                    exit\n                } \n                #gets virtual machine objects from the Virtual Machine Manager database\n                Set-SCVMMServer -VMMServer $vmmServer\n\u0009\u0009$VM = Get-SCVirtualMachine | Where-Object {$_.Name -Eq $vmName }   \n                #check if VM Exists\n               \n                    \n                          if($diskdriveName -ne $null)\n                          {\n                             #gets the specified volume object\n                             $diskdrive=Get-SCVirtualDiskDrive -VM $VM | Where-Object { $_.VirtualHardDisk.Name -Eq $diskdriveName} \n                             if($diskdrive -eq $null)\n                             {\n                               #This will delete specified volume attached which is to that vm \n                               Write-Error \"Virtual Disk Not Found\"\n                             }\n\n                          }  \n                          else\n                           {\n                                echo \"Name of the disk drive  is not entered\"\n                           }                             \n         }\n\u0009     catch [Exception]\n\u0009       {\n\u0009\u0009        echo $_.Exception.Message\n\u0009        }\n}\n"
		arguments := "terraformVM " + os.Getenv("SCVMM_VMM_SERVER") + " terraformDisk"
		filename := "creatediskdrive"
		result1, err := execScript(org, script, filename, arguments)
		if err != "" {
			return fmt.Errorf("Error, Volume is not created  %v", result1)
		}
		return nil
	}
}

func testAccCheckVDConfigBasic() string {
	return fmt.Sprintf(`
resource "scvmm_virtual_machine" "CreateVM" {
  timeout = "1000"
  vmm_server = "%s"
  vm_name = "terraformVM"
  template_name = "%s"
  cloud_name = "%s"
}

resource "scvmm_virtual_disk" "CreateDiskDrive" {
  timeout = 1000
  vmm_server = "${scvmm_virtual_machine.CreateVM.vmm_server}"
  vm_name = "${scvmm_virtual_machine.CreateVM.vm_name}"
  virtual_disk_name = "terraformDisk"
  virtual_disk_size = "1000"
}
`, os.Getenv("SCVMM_VMM_SERVER"),
		os.Getenv("SCVMM_TEMPLATE_NAME"),
		os.Getenv("SCVMM_CLOUD_NAME"))
}
