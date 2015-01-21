package main

import (
	"testing"
)

func MergeBuildParamsTest(t *testing.T) {
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
