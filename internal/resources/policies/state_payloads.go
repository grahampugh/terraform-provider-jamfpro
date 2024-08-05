package policies

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"log"
	"reflect"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Parent func for stating payloads. Constructs var with prep funcs and states as one here.
func statePayloads(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	out := make([]map[string]interface{}, 0)
	out = append(out, make(map[string]interface{}, 1))

	// DiskEncryption
	prepStatePayloadDiskEncryption(&out, resp)

	// Packages
	prepStatePayloadPackages(&out, resp)

	// Scripts
	prepStatePayloadScripts(&out, resp)

	// Printers
	prepStatePayloadPrinters(&out, resp)

	// Dock Items
	prepStatePayloadDockItems(&out, resp)

	// Account Maintenance
	prepStatePayloadAccountMaintenance(&out, resp)

	// Files Processes
	prepStatePayloadFilesProcesses(&out, resp)

	// User Interaction
	prepStatePayloadUserInteraction(&out, resp)

	// Reboot
	prepStatePayloadReboot(&out, resp)

	// Maintenance
	prepStatePayloadMaintenance(&out, resp)

	// State
	err := d.Set("payloads", out)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}

// prepStatePayloadDiskEncryption reads response and preps disk encryption payload items for stating
func prepStatePayloadDiskEncryption(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if &resp.DiskEncryption == nil {
		log.Println("No disk encryption configuration found")
		return
	}

	// Define default values
	defaults := map[string]interface{}{
		"action":                           "none",
		"disk_encryption_configuration_id": 0,
		"auth_restart":                     false,
		"remediate_key_type":               "",
		"remediate_disk_encryption_configuration_id": 0,
	}

	diskEncryptionBlock := map[string]interface{}{
		"action":                           resp.DiskEncryption.Action,
		"disk_encryption_configuration_id": resp.DiskEncryption.DiskEncryptionConfigurationID,
		"auth_restart":                     resp.DiskEncryption.AuthRestart,
		"remediate_key_type":               resp.DiskEncryption.RemediateKeyType,
		"remediate_disk_encryption_configuration_id": resp.DiskEncryption.RemediateDiskEncryptionConfigurationID,
	}

	// Check if all values are default
	allDefault := true
	for key, value := range diskEncryptionBlock {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		log.Println("All disk encryption values are default, skipping state")
		return
	}

	log.Println("Initializing disk encryption in state")
	(*out)[0]["disk_encryption"] = []map[string]interface{}{diskEncryptionBlock}
	log.Printf("Final state disk encryption: %+v\n", diskEncryptionBlock)
}

// Reads response and preps package payload items
func prepStatePayloadPackages(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if &resp.PackageConfiguration == nil || resp.PackageConfiguration.Packages == nil {
		log.Println("No package configuration found")
		return
	}

	log.Println("Initializing packages in state")

	packagesMap := make(map[string]interface{})
	packagesMap["distribution_point"] = resp.PackageConfiguration.DistributionPoint

	packagesMap["package"] = make([]map[string]interface{}, 0)
	for _, v := range resp.PackageConfiguration.Packages {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["action"] = v.Action
		outMap["fill_user_template"] = v.FillUserTemplate
		outMap["fill_existing_user_template"] = v.FillExistingUsers
		packagesMap["package"] = append(packagesMap["package"].([]map[string]interface{}), outMap)
	}

	(*out)[0]["packages"] = []map[string]interface{}{packagesMap}
	log.Printf("Final state packages: %+v\n", (*out)[0]["packages"])
}

// Reads response and preps script payload items
func prepStatePayloadScripts(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Scripts == nil {
		log.Println("No scripts found")
		return
	}

	log.Println("Initializing scripts in state")
	(*out)[0]["scripts"] = make([]map[string]interface{}, 0)

	for _, v := range resp.Scripts {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["priority"] = v.Priority

		if v.Parameter4 != "" {
			outMap["parameter4"] = v.Parameter4
		}
		if v.Parameter5 != "" {
			outMap["parameter5"] = v.Parameter5
		}
		if v.Parameter6 != "" {
			outMap["parameter6"] = v.Parameter6
		}
		if v.Parameter7 != "" {
			outMap["parameter7"] = v.Parameter7
		}
		if v.Parameter8 != "" {
			outMap["parameter8"] = v.Parameter8
		}
		if v.Parameter9 != "" {
			outMap["parameter9"] = v.Parameter9
		}
		if v.Parameter10 != "" {
			outMap["parameter10"] = v.Parameter10
		}
		if v.Parameter11 != "" {
			outMap["parameter11"] = v.Parameter11
		}
		log.Printf("Adding script to state: %+v\n", outMap)
		(*out)[0]["scripts"] = append((*out)[0]["scripts"].([]map[string]interface{}), outMap)
	}

	log.Printf("Final state scripts: %+v\n", (*out)[0]["scripts"])
}

// prepStatePayloadPrinters reads response and preps printer payload items for stating
func prepStatePayloadPrinters(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.Printers.Printer == nil {
		log.Println("No printers found")
		return
	}

	log.Println("Initializing printers in state")
	(*out)[0]["printers"] = make([]map[string]interface{}, 0)

	for _, v := range *resp.Printers.Printer {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["name"] = v.Name
		outMap["action"] = v.Action
		outMap["make_default"] = v.MakeDefault

		log.Printf("Adding printer to state: %+v\n", outMap)
		(*out)[0]["printers"] = append((*out)[0]["printers"].([]map[string]interface{}), outMap)
	}

	log.Printf("Final state printers: %+v\n", (*out)[0]["printers"])
}

// Reads response and preps dock items payload items
func prepStatePayloadDockItems(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if resp.DockItems == nil {
		log.Println("No dock items found")
		return
	}

	log.Println("Initializing dock items in state")
	(*out)[0]["dock_items"] = make([]map[string]interface{}, 0)

	for _, v := range resp.DockItems {
		outMap := make(map[string]interface{})
		outMap["id"] = v.ID
		outMap["name"] = v.Name
		outMap["action"] = v.Action

		log.Printf("Adding dock item to state: %+v\n", outMap)
		(*out)[0]["dock_items"] = append((*out)[0]["dock_items"].([]map[string]interface{}), outMap)
	}

	log.Printf("Final state dock items: %+v\n", (*out)[0]["dock_items"])
}

// prepStatePayloadAccountMaintenance reads response and preps account maintenance payload items.
// If all values are default, do not set the account_maintenance block
func prepStatePayloadAccountMaintenance(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if &resp.AccountMaintenance == nil {
		log.Println("No account maintenance configuration found")
		return
	}

	log.Println("Initializing account maintenance in state")
	accountMaintenanceMap := make(map[string]interface{})

	// Handle accounts
	if resp.AccountMaintenance.Accounts != nil {
		localAccounts := make([]map[string]interface{}, 0)
		for _, v := range *resp.AccountMaintenance.Accounts {
			accountMap := make(map[string]interface{})
			accountMap["action"] = v.Action
			accountMap["username"] = v.Username
			accountMap["realname"] = v.Realname
			accountMap["password"] = v.Password
			accountMap["archive_home_directory"] = v.ArchiveHomeDirectory
			accountMap["archive_home_directory_to"] = v.ArchiveHomeDirectoryTo
			accountMap["home"] = v.Home
			accountMap["hint"] = v.Hint
			accountMap["picture"] = v.Picture
			accountMap["admin"] = v.Admin
			accountMap["filevault_enabled"] = v.FilevaultEnabled

			log.Printf("Adding account to state: %+v\n", accountMap)
			localAccounts = append(localAccounts, accountMap)
		}

		if len(localAccounts) > 0 {
			accountMaintenanceMap["local_accounts"] = []map[string]interface{}{
				{"account": localAccounts},
			}
		}
	}

	// Handle directory bindings
	if resp.AccountMaintenance.DirectoryBindings != nil {
		directoryBindings := make([]map[string]interface{}, 0)
		for _, v := range *resp.AccountMaintenance.DirectoryBindings {
			bindingMap := make(map[string]interface{})
			bindingMap["id"] = v.ID
			bindingMap["name"] = v.Name

			log.Printf("Adding directory binding to state: %+v\n", bindingMap)
			directoryBindings = append(directoryBindings, bindingMap)
		}

		if len(directoryBindings) > 0 {
			accountMaintenanceMap["directory_bindings"] = []map[string]interface{}{
				{"binding": directoryBindings},
			}
		}
	}

	// Handle management account
	if resp.AccountMaintenance.ManagementAccount != nil {
		managementAccountMap := make(map[string]interface{})
		if resp.AccountMaintenance.ManagementAccount.Action != "doNotChange" || resp.AccountMaintenance.ManagementAccount.ManagedPassword != "" || resp.AccountMaintenance.ManagementAccount.ManagedPasswordLength != 0 {
			managementAccountMap["action"] = resp.AccountMaintenance.ManagementAccount.Action
			managementAccountMap["managed_password"] = resp.AccountMaintenance.ManagementAccount.ManagedPassword
			managementAccountMap["managed_password_length"] = resp.AccountMaintenance.ManagementAccount.ManagedPasswordLength

			log.Printf("Adding management account to state: %+v\n", managementAccountMap)
			accountMaintenanceMap["management_account"] = []map[string]interface{}{managementAccountMap}
		}
	}

	// Handle open firmware/EFI password
	if resp.AccountMaintenance.OpenFirmwareEfiPassword != nil {
		openFirmwareMap := make(map[string]interface{})
		if resp.AccountMaintenance.OpenFirmwareEfiPassword.OfMode != "none" || resp.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword != "" {
			openFirmwareMap["of_mode"] = resp.AccountMaintenance.OpenFirmwareEfiPassword.OfMode
			openFirmwareMap["of_password"] = resp.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword

			log.Printf("Adding open firmware/EFI password to state: %+v\n", openFirmwareMap)
			accountMaintenanceMap["open_firmware_efi_password"] = []map[string]interface{}{openFirmwareMap}
		}
	}

	if len(accountMaintenanceMap) > 0 {
		(*out)[0]["account_maintenance"] = []map[string]interface{}{accountMaintenanceMap}
		log.Printf("Final state account maintenance: %+v\n", (*out)[0]["account_maintenance"])
	}
}

// prepStatePayloadFilesProcesses reads response and preps files and processes payload items.
func prepStatePayloadFilesProcesses(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if &resp.FilesProcesses == nil {
		log.Println("No files and processes configuration found")
		return
	}

	// Define default values
	defaults := map[string]interface{}{
		"search_by_path":         "",
		"delete_file":            false,
		"locate_file":            "",
		"update_locate_database": false,
		"spotlight_search":       "",
		"search_for_process":     "",
		"kill_process":           false,
		"run_command":            "",
	}

	filesProcessesBlock := map[string]interface{}{
		"search_by_path":         resp.FilesProcesses.SearchByPath,
		"delete_file":            resp.FilesProcesses.DeleteFile,
		"locate_file":            resp.FilesProcesses.LocateFile,
		"update_locate_database": resp.FilesProcesses.UpdateLocateDatabase,
		"spotlight_search":       resp.FilesProcesses.SpotlightSearch,
		"search_for_process":     resp.FilesProcesses.SearchForProcess,
		"kill_process":           resp.FilesProcesses.KillProcess,
		"run_command":            resp.FilesProcesses.RunCommand,
	}

	// Check if all values are default
	allDefault := true
	for key, value := range filesProcessesBlock {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		log.Println("All files and processes values are default, skipping state")
		return
	}

	log.Println("Initializing files and processes in state")
	(*out)[0]["files_processes"] = []map[string]interface{}{filesProcessesBlock}
	log.Printf("Final state files and processes: %+v\n", filesProcessesBlock)
}

// prepStatePayloadUserInteraction Reads response and preps user interaction payload items. If all values are default, do not set the user_interaction block
func prepStatePayloadUserInteraction(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if &resp.UserInteraction == nil {
		log.Println("No user interaction configuration found")
		return
	}

	// Define default values
	defaults := map[string]interface{}{
		"message_start":            "",
		"allow_user_to_defer":      false,
		"allow_deferral_until_utc": "",
		"allow_deferral_minutes":   0,
		"message_finish":           "",
	}

	userInteractionBlock := map[string]interface{}{
		"message_start":            resp.UserInteraction.MessageStart,
		"allow_user_to_defer":      resp.UserInteraction.AllowUserToDefer,
		"allow_deferral_until_utc": resp.UserInteraction.AllowDeferralUntilUtc,
		"allow_deferral_minutes":   resp.UserInteraction.AllowDeferralMinutes,
		"message_finish":           resp.UserInteraction.MessageFinish,
	}

	// Check if all values are default
	allDefault := true
	for key, value := range userInteractionBlock {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		log.Println("All user interaction values are default, skipping state")
		return
	}

	log.Println("Initializing user interaction in state")
	(*out)[0]["user_interaction"] = []map[string]interface{}{userInteractionBlock}
	log.Printf("Final state user interaction: %+v\n", userInteractionBlock)
}

// Reads response and preps reboot payload items
func prepStatePayloadReboot(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if &resp.Reboot == nil {
		log.Println("No reboot configuration found")
		return
	}

	defaults := map[string]interface{}{
		"message":                        "This computer will restart in 5 minutes. Please save anything you are working on and log out by choosing Log Out from the bottom of the Apple menu.",
		"specify_startup":                "",
		"startup_disk":                   "Current Startup Disk",
		"no_user_logged_in":              "Do not restart",
		"user_logged_in":                 "Do not restart",
		"minutes_until_reboot":           5,
		"start_reboot_timer_immediately": false,
		"file_vault_2_reboot":            false,
	}

	rebootBlock := map[string]interface{}{
		"message":                        resp.Reboot.Message,
		"specify_startup":                resp.Reboot.SpecifyStartup,
		"startup_disk":                   resp.Reboot.StartupDisk,
		"no_user_logged_in":              resp.Reboot.NoUserLoggedIn,
		"user_logged_in":                 resp.Reboot.UserLoggedIn,
		"minutes_until_reboot":           resp.Reboot.MinutesUntilReboot,
		"start_reboot_timer_immediately": resp.Reboot.StartRebootTimerImmediately,
		"file_vault_2_reboot":            resp.Reboot.FileVault2Reboot,
	}

	allDefault := true
	for key, value := range rebootBlock {
		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault {
		log.Println("All user interaction values are default, skipping state")
		return
	}

	log.Println("Initializing reboot in state")
	(*out)[0]["reboot"] = []map[string]interface{}{rebootBlock}
	log.Printf("Final state reboot: %+v\n", rebootBlock)
}

// prepStatePayloadMaintenance Reads response and preps maintenance payload items. If all values are default, do not set the maintenance block
func prepStatePayloadMaintenance(out *[]map[string]interface{}, resp *jamfpro.ResourcePolicy) {
	if &resp.Maintenance == nil {
		return
	}

	// Do not set the maintenance block if all values are default (false)
	v := reflect.ValueOf(resp.Maintenance)
	allDefault := true

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Bool() {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}
	// Else, set the maintenance block
	(*out)[0]["maintenance"] = make([]map[string]interface{}, 0)
	outMap := make(map[string]interface{})
	outMap["recon"] = resp.Maintenance.Recon
	outMap["reset_name"] = resp.Maintenance.ResetName
	outMap["install_all_cached_packages"] = resp.Maintenance.InstallAllCachedPackages
	outMap["heal"] = resp.Maintenance.Heal
	outMap["prebindings"] = resp.Maintenance.Prebindings
	outMap["permissions"] = resp.Maintenance.Permissions
	outMap["byhost"] = resp.Maintenance.Byhost
	outMap["system_cache"] = resp.Maintenance.SystemCache
	outMap["user_cache"] = resp.Maintenance.UserCache
	outMap["verify"] = resp.Maintenance.Verify
	(*out)[0]["maintenance"] = append((*out)[0]["maintenance"].([]map[string]interface{}), outMap)
}
