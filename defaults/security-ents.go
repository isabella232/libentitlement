package defaults

import (
	"github.com/docker/libentitlement/defaults/osdefs"
	"github.com/docker/libentitlement/entitlement"
	"github.com/docker/libentitlement/secprofile"
	"github.com/opencontainers/runtime-spec/specs-go"
)

const (
	securityDomain = "security"
)

const (
	// SecurityConfinedEntFullID is the ID for the security.confined entitlement
	SecurityConfinedEntFullID = securityDomain + ".confined"
	// SecurityViewEntFullID is the ID for the security.view entitlement
	SecurityViewEntFullID = securityDomain + ".view"
	// SecurityAdminEntFullID is the ID for the security.admin entitlement
	SecurityAdminEntFullID = securityDomain + ".admin"
	// SecurityMemoryLockFullID is the ID for the security.memory-lock entitlement
	SecurityMemoryLockFullID = securityDomain + ".memory-lock"
)

var (
	securityConfinedEntitlement   = entitlement.NewVoidEntitlement(SecurityConfinedEntFullID, securityConfinedEntitlementEnforce)
	securityViewEntitlement       = entitlement.NewVoidEntitlement(SecurityViewEntFullID, securityViewEntitlementEnforce)
	securityAdminEntitlement      = entitlement.NewVoidEntitlement(SecurityAdminEntFullID, securityAdminEntitlementEnforce)
	securityMemoryLockEntitlement = entitlement.NewVoidEntitlement(SecurityMemoryLockFullID, securityMemoryLockEnforce)
)

func securityConfinedEntitlementEnforce(profile secprofile.Profile) (secprofile.Profile, error) {
	ociProfile, err := ociProfileConversionCheck(profile, SecurityConfinedEntFullID)
	if err != nil {
		return nil, err
	}

	capsToRemove := []string{
		CapMacAdmin, CapMacOverride, CapDacOverride, CapDacReadSearch, CapSetfcap, CapSetfcap, CapSetuid, CapSetgid,
		CapSysPtrace, CapFsetid, CapSysModule, CapSyslog, CapSysRawio, CapSysAdmin, CapLinuxImmutable,
	}
	ociProfile.RemoveCaps(capsToRemove...)

	syscallsToBlock := []string{
		SysPtrace, SysArchPrctl, SysPersonality, SysPersonality, SysSetuid, SysSetgid, SysPrctl, SysMadvise,
	}
	ociProfile.BlockSyscalls(syscallsToBlock...)

	syscallsWithArgsToAllow := map[string][]specs.LinuxSeccompArg{
		SysPrctl: {
			{
				Index: 0,
				Value: osdefs.PrCapbsetDrop,
				Op:    specs.OpNotEqual,
			},
			{
				Index: 0,
				Value: osdefs.PrCapbsetRead,
				Op:    specs.OpNotEqual,
			},
		},
	}
	ociProfile.AllowSyscallsWithArgs(syscallsWithArgsToAllow)

	/* FIXME: Add AppArmor rules to deny RW on sensitive FS directories */

	return ociProfile, nil
}

func securityViewEntitlementEnforce(profile secprofile.Profile) (secprofile.Profile, error) {
	ociProfile, err := ociProfileConversionCheck(profile, SecurityViewEntFullID)
	if err != nil {
		return nil, err
	}

	capsToRemove := []string{
		CapSysAdmin, CapSysPtrace, CapSetuid, CapSetgid, CapSetpcap, CapSetfcap, CapMacAdmin, CapMacOverride,
		CapDacOverride, CapFsetid, CapSysModule, CapSyslog, CapSysRawio, CapLinuxImmutable,
	}
	ociProfile.RemoveCaps(capsToRemove...)

	capsToAdd := []string{CapDacReadSearch}
	ociProfile.AddCaps(capsToAdd...)

	syscallsToBlock := []string{
		SysPtrace, SysArchPrctl, SysPersonality, SysSetuid, SysSetgid, SysPrctl, SysMadvise,
	}
	ociProfile.BlockSyscalls(syscallsToBlock...)

	syscallsWithArgsToAllow := map[string][]specs.LinuxSeccompArg{
		SysPrctl: {
			{
				Index: 0,
				Value: osdefs.PrCapbsetDrop,
				Op:    specs.OpNotEqual,
			},
		},
	}
	ociProfile.AllowSyscallsWithArgs(syscallsWithArgsToAllow)

	/* FIXME: Add AppArmor rules to RO on sensitive FS directories */

	return ociProfile, nil
}

func securityAdminEntitlementEnforce(profile secprofile.Profile) (secprofile.Profile, error) {
	ociProfile, err := ociProfileConversionCheck(profile, SecurityAdminEntFullID)
	if err != nil {
		return nil, err
	}

	capsToAdd := []string{
		CapMacAdmin, CapMacOverride, CapDacOverride, CapDacReadSearch, CapSetpcap, CapSetfcap, CapSetuid, CapSetgid,
		CapSysPtrace, CapFsetid, CapSysModule, CapSyslog, CapSysRawio, CapSysAdmin, CapLinuxImmutable,
	}
	ociProfile.AddCaps(capsToAdd...)

	syscallsToAllow := []string{
		SysPtrace, SysArchPrctl, SysPersonality, SysSetuid, SysSetgid, SysPrctl, SysMadvise,
	}
	ociProfile.AllowSyscalls(syscallsToAllow...)

	return ociProfile, nil
}

func securityMemoryLockEnforce(profile secprofile.Profile) (secprofile.Profile, error) {
	ociProfile, err := ociProfileConversionCheck(profile, SecurityMemoryLockFullID)
	if err != nil {
		return nil, err
	}

	capsToAdd := []string{
		CapIpcLock,
	}
	ociProfile.AddCaps(capsToAdd...)

	syscallsToAllow := []string{
		SysMlock, SysMunlock, SysMlock2, SysMlockall, SysMunlockall,
	}
	ociProfile.AllowSyscalls(syscallsToAllow...)

	return ociProfile, nil
}
