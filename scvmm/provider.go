package scvmm

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider ... provides scvmm capability to terraform
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SCVMM Server IP",
				DefaultFunc: schema.EnvDefaultFunc("SCVMM_SERVER_IP", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server Port",
				DefaultFunc: schema.EnvDefaultFunc("SCVMM_SERVER_PORT", nil),
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name",
				DefaultFunc: schema.EnvDefaultFunc("SCVMM_SERVER_USER", nil),
			},
			"user_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password for provided user_name",
				DefaultFunc: schema.EnvDefaultFunc("SCVMM_SERVER_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"scvmm_virtual_machine": resourceSCVMMVirtualMachine(),
			"scvmm_virtual_disk":    resourceSCVMMVirtualDisk(),
			"scvmm_checkpoint":      resourceSCVMMCheckpoint(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ServerIP: d.Get("server_ip").(string),
		Port:     d.Get("port").(int),
		Username: d.Get("user_name").(string),
		Password: d.Get("user_password").(string),
	}

	log.Println("[DEBUG] Initializing Winrm Connection")
	return config.Connection()
}
