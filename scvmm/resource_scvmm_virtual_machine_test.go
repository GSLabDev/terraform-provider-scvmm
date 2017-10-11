package scvmm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/masterzen/winrm"
)

func testBasicPreCheckVirtualMachine(t *testing.T) {
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

func TestAccVirtualMachineCreate_Basic(t *testing.T) {
	resourceName := "scvmm_virtual_machine.CreateVM"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testBasicPreCheckVirtualMachine(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineDestroy(resourceName),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVirtualMachineConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "vm_name", "terraformVM"),
					resource.TestCheckResourceAttr(
						resourceName, "vmm_server", os.Getenv("SCVMM_VMM_SERVER")),
				),
			},
		},
	})
}

func testAccCheckVirtualMachineDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		org := testAccProvider.Meta().(*winrm.Client)
		script := "[CmdletBinding(SupportsShouldProcess=$true)]\nparam (\n  \n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VmmServer\")]\n  [string]$vmmServer,\n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VM Name\")]\n  [string]$vmName\n)\n\nBegin\n{\n         If (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] \"Administrator\"))\n          {   \n            $arguments = \"\" + $myinvocation.mycommand.definition + \" \"\n            $myinvocation.BoundParameters.Values | foreach{$arguments += \"'$_' \" }\n            echo $arguments\n            Start-Process powershell -Verb runAs -ArgumentList $arguments\n            Break\n         }\n\u0009     try\n\u0009     {\n               if($vmName -eq $null) \n               {\n                    echo \"VM Name not entered\"\n                    exit\n               } \n               #gets virtual machine objects from the Virtual Machine Manager database\n               Set-SCVMMServer -VMMServer $vmmServer > $null\n\u0009\u0009       $VM = Get-SCVirtualMachine | Where-Object {$_.Name -Eq $vmName }   \n               #check if VM Exists\n               if($VM -eq $null)\n               {     \n                   Write-Error \"VM does not exists\"\n                   exit\n               }\n            \n         }\n\u0009     catch [Exception]\n         {\n               Write-Error $_.Exception.Message\n\u0009     }\n      \n   \n}\n"
		arguments := rs.Primary.Attributes["vmm_server"] + " " + rs.Primary.Attributes["vm_name"]
		filename := "DeleteVM_Test"
		result, err := execScript(org, script, filename, arguments)
		if err == "" {
			return fmt.Errorf("Virtual Machine still exists: %v", result)
		}
		return nil
	}
}

func testAccCheckVirtualMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vm ID is set")
		}
		org := testAccProvider.Meta().(*winrm.Client)

		script := "[CmdletBinding(SupportsShouldProcess=$true)]\nparam (\n  \n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VmmServer\")]\n  [string]$vmmServer,\n  [Parameter(Mandatory=$true,HelpMessage=\"Enter VM Name\")]\n  [string]$vmName\n)\n\nBegin\n{\n         If (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] \"Administrator\"))\n          {   \n            $arguments = \"\" + $myinvocation.mycommand.definition + \" \"\n            $myinvocation.BoundParameters.Values | foreach{$arguments += \"'$_' \" }\n            echo $arguments\n            Start-Process powershell -Verb runAs -ArgumentList $arguments\n            Break\n         }\n\u0009     try\n\u0009     {\n               if($vmName -eq $null) \n               {\n                    echo \"VM Name not entered\"\n                    exit\n               } \n               #gets virtual machine objects from the Virtual Machine Manager database\n               Set-SCVMMServer -VMMServer $vmmServer > $null\n\u0009\u0009       $VM = Get-SCVirtualMachine | Where-Object {$_.Name -Eq $vmName }   \n               #check if VM Exists\n               if($VM -eq $null)\n               {     \n                   Write-Error \"VM does not exists\"\n                   exit\n               }\n            \n         }\n\u0009     catch [Exception]\n         {\n               Write-Error $_.Exception.Message\n\u0009     }\n      \n   \n}\n"
		arguments := os.Getenv("SCVMM_VMM_SERVER") + " terraformVM"
		filename := "CreateVM_Test"
		result, err := execScript(org, script, filename, arguments)
		if err != "" {
			return fmt.Errorf("Error while getting the VM %v", result)
		}
		return nil
	}
}

func testAccCheckVirtualMachineConfigBasic() string {
	return fmt.Sprintf(`
resource "scvmm_virtual_machine" "CreateVM" {
  timeout = "1000"
  vmm_server = "%s"
  vm_name = "terraformVM"
  template_name = "%s"
  cloud_name = "%s"
}`, os.Getenv("SCVMM_VMM_SERVER"),
		os.Getenv("SCVMM_TEMPLATE_NAME"),
		os.Getenv("SCVMM_CLOUD_NAME"))
}
