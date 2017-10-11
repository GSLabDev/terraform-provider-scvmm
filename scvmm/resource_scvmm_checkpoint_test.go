package scvmm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/masterzen/winrm"
)

func testBasicPreCheckSP(t *testing.T) {
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

func TestAccsp_Basic(t *testing.T) {

	resourceName := "scvmm_checkpoint.CreateCheckpoint"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testBasicPreCheckSP(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCheckpointDestroy(resourceName),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckCheckpointConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCheckpointExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "vm_name", "terraformVM"),
					resource.TestCheckResourceAttr(
						resourceName, "checkpoint_name", "terraformCheckpoint"),
				),
			},
		},
	})
}

func testAccCheckCheckpointDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		org := testAccProvider.Meta().(*winrm.Client)

		script := "[CmdletBinding(SupportsShouldProcess=$true)]\nparam(\n    [parameter(Mandatory=$true,HelpMessage=\"Enter VMMServer\")]\n    [string]$vmmServer,\n\n    [parameter(Mandatory=$true,HelpMessage=\"Enter Virtual Machine Name\")]\n    [string]$vmName,\n\n    [parameter(Mandatory=$true,HelpMessage=\"Enter Checkpoint Name\")]\n    [string]$checkpointName\n)\nBegin\n{  \n       If (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] \"Administrator\"))\n    {   $arguments = \"\" + $myinvocation.mycommand.definition + \" \"\n        $myinvocation.BoundParameters.Values | foreach{\n            $arguments += \"'$_' \"\n        }\n        echo $arguments\n        Start-Process powershell -Verb runAs -ArgumentList $arguments\n        Break\n    }                 \n    \n        try\n\u0009     {\n\u0009\u0009 Set-SCVMMServer -VMMServer $vmmServer > $null\n                 \n                 $checkpoint = Get-SCVMCheckpoint -VM $vmName | Where-Object {$_.Name -eq $checkpointName}\n                 if($checkpoint-eq $null)\n                  {\n                    Write-Error \"No Checkpoint found\"\n                  }             \n             }catch [Exception]\n\u0009        {\n\u0009\u0009        echo $_.Exception.Message\n                }    \n}"
		arguments := rs.Primary.Attributes["vmm_server"] + " " + rs.Primary.Attributes["vm_name"] + " " + rs.Primary.Attributes["checkpoint_name"]
		filename := "deletecp"
		result, err := execScript(org, script, filename, arguments)

		if err == "" {
			return fmt.Errorf("Checkpoint  still exists: %v", result)
		}
		return nil
	}

}

func testAccCheckCheckpointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vm ID is set")
		}

		org := testAccProvider.Meta().(*winrm.Client)

		script := "[CmdletBinding(SupportsShouldProcess=$true)]\nparam(\n    [parameter(Mandatory=$true,HelpMessage=\"Enter VMMServer\")]\n    [string]$vmmServer,\n\n    [parameter(Mandatory=$true,HelpMessage=\"Enter Virtual Machine Name\")]\n    [string]$vmName,\n\n    [parameter(Mandatory=$true,HelpMessage=\"Enter Checkpoint Name\")]\n    [string]$checkpointName\n)\nBegin\n{  \n       If (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] \"Administrator\"))\n    {   $arguments = \"\" + $myinvocation.mycommand.definition + \" \"\n        $myinvocation.BoundParameters.Values | foreach{\n            $arguments += \"'$_' \"\n        }\n        echo $arguments\n        Start-Process powershell -Verb runAs -ArgumentList $arguments\n        Break\n    }                 \n    \n        try\n\u0009     {\n\u0009\u0009 Set-SCVMMServer -VMMServer $vmmServer > $null\n                 \n                 $checkpoint = Get-SCVMCheckpoint -VM $vmName | Where-Object {$_.Name -eq $checkpointName}\n                 if($checkpoint-eq $null)\n                  {\n                    Write-Error \"No Checkpoint found\"\n                  }             \n             }catch [Exception]\n\u0009        {\n\u0009\u0009        echo $_.Exception.Message\n                }    \n}"
		arguments := rs.Primary.Attributes["vmm_server"] + " " + rs.Primary.Attributes["vm_name"] + " " + rs.Primary.Attributes["checkpoint_name"]
		filename := "createcp"
		result, err := execScript(org, script, filename, arguments)

		if err != "" {
			return fmt.Errorf("Error while getting the checkpoint %v", result)
		}

		return nil
	}
}

func testAccCheckCheckpointConfigBasic() string {
	return fmt.Sprintf(`
resource "scvmm_virtual_machine" "CreateVM" {
  timeout = "1000"
  vmm_server = "%s"
  vm_name = "terraformVM"
  template_name = "%s"
  cloud_name = "%s"
}
resource "scvmm_checkpoint" "CreateCheckpoint" {
  timeout = "1000"
  vmm_server = "${scvmm_virtual_machine.CreateVM.vmm_server}"
  vm_name = "${scvmm_virtual_machine.CreateVM.vm_name}"
  checkpoint_name= "terraformCheckpoint"
}`, os.Getenv("SCVMM_VMM_SERVER"),
		os.Getenv("SCVMM_TEMPLATE_NAME"),
		os.Getenv("SCVMM_CLOUD_NAME"))

}
