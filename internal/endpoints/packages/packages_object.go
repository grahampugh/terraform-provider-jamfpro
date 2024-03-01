// packages_data_object.go
package packages

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func constructJamfProJCDS2Package(d *schema.ResourceData) (string, error) {
	// Extract the 'package_file_path' attribute from the Terraform resource data
	filePath, ok := d.Get("package_file_path").(string)
	if !ok || filePath == "" {
		// Return an error if the file path is not provided or is empty
		return "", fmt.Errorf("file path for the Jamf Pro package is required and cannot be empty")
	}

	return filePath, nil
}

func constructJamfProPackage(d *schema.ResourceData) (*jamfpro.ResourcePackage, error) {

	packageResource := &jamfpro.ResourcePackage{
		Name:                       d.Get("name").(string),
		Filename:                   d.Get("filename").(string),
		Category:                   d.Get("category").(string),
		Info:                       d.Get("info").(string),
		Notes:                      d.Get("notes").(string),
		Priority:                   d.Get("priority").(int),
		RebootRequired:             d.Get("reboot_required").(bool),
		FillUserTemplate:           d.Get("fill_user_template").(bool),
		FillExistingUsers:          d.Get("fill_existing_users").(bool),
		BootVolumeRequired:         d.Get("boot_volume_required").(bool),
		AllowUninstalled:           d.Get("allow_uninstalled").(bool),
		OSRequirements:             d.Get("os_requirements").(string),
		RequiredProcessor:          d.Get("required_processor").(string),
		SwitchWithPackage:          d.Get("switch_with_package").(string),
		InstallIfReportedAvailable: d.Get("install_if_reported_available").(bool),
		ReinstallOption:            d.Get("reinstall_option").(string),
		TriggeringFiles:            d.Get("triggering_files").(string),
		SendNotification:           d.Get("send_notification").(bool),
	}

	return packageResource, nil
}
