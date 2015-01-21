package main

import (
	"testing"
)

func TestMergeBuildParams(t *testing.T) {
	bp1 := BuildParams{
		ProjectFileDir:      "pfd1",
		BuildTool:           "bt1",
		ProjectFile:         "pf1",
		SchemeName:          "sn1",
		DeviceDestination:   "dd1",
		BuildOutputFilePath: "bofp1",
	}

	bp2 := BuildParams{
		ProjectFileDir:      "pfd2",
		BuildTool:           "bt2",
		ProjectFile:         "pf2",
		SchemeName:          "sn2",
		DeviceDestination:   "dd2",
		BuildOutputFilePath: "bofp2",
	}

	mergedBuildParams := MergeBuildParams(bp1, bp2)

	if mergedBuildParams.ProjectFileDir != bp1.ProjectFileDir {
		t.Error("Merge failed: ProjectFileDir")
	}
}

func TestValidate(t *testing.T) {
	// valid BuildParams
	bParams := BuildParams{
		ProjectFileDir:      "pfd1",
		BuildTool:           "bt1",
		ProjectFile:         "pf1",
		SchemeName:          "sn1",
		DeviceDestination:   "dd1",
		BuildOutputFilePath: "bofp1",
		CodeSignIdentity:    "sign-identity",
		KeychainName:        "keychain-name",
		KeychainPassword:    "keychain-psw",
		ProvisioningProfile: "prov-profile",
	}

	if err := bParams.Validate(); err != nil {
		t.Error("BuildParams should be valid but got: ", err)
	}

	// remove ProvProfile - should be still valid (it's optional)
	bParams.ProvisioningProfile = ""
	if err := bParams.Validate(); err != nil {
		t.Error("BuildParams should be valid but got: ", err)
	}

	// remove SchemeName - it should NOT be valid anymore
	bParams.SchemeName = ""
	if err := bParams.Validate(); err == nil {
		t.Error("BuildParams should NOT be valid but got no error!")
	}
}
