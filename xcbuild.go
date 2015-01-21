package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunCommandInDirWithArgsAndWriters(dirPath string, command string, cmdArgs []string, stdOutWriter, stdErrWriter io.Writer) error {
	c := exec.Command(command, cmdArgs...)
	c.Stdout = stdOutWriter
	c.Stderr = stdErrWriter
	if dirPath != "" {
		c.Dir = dirPath
	}

	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func RunCommandInDirWithArgs(dirPath string, command string, cmdArgs []string) error {
	return RunCommandInDirWithArgsAndWriters(dirPath, command, cmdArgs, os.Stdout, os.Stderr)
}

func ExecuteUnlockKeychain(keychainName, keychainPsw string) error {
	cargs := []string{
		"-v", "unlock-keychain", "-p", keychainPsw, keychainName,
	}
	err := RunCommandInDirWithArgsAndWriters("", "security", cargs, BuildLogWriter, BuildLogWriter)
	return err
}

func ExecuteBuildWithParams(buildParams BuildParams) error {
	projActionArg := ""
	if strings.HasSuffix(buildParams.ProjectFile, ".xcodeproj") || strings.HasSuffix(buildParams.ProjectFile, ".xcodeproj/") {
		projActionArg = "-project"
	} else if strings.HasSuffix(buildParams.ProjectFile, ".xcworkspace") || strings.HasSuffix(buildParams.ProjectFile, ".xcworkspace/") {
		projActionArg = "-workspace"
	}
	if projActionArg == "" {
		return errors.New("Invalid project file - can't determine the project file type!")
	}

	if err := WriteLineToBuildLog("[[build-start]]"); err != nil {
		return err
	}

	// unlock keychain
	if err := ExecuteUnlockKeychain(buildParams.KeychainName, buildParams.KeychainPassword); err != nil {
		return err
	}

	buildCmd := fmt.Sprintf(`%s %s "%s" -scheme "%s" -destination "%s" -sdk iphonesimulator clean test CODE_SIGN_IDENTITY="%s" OTHER_CODE_SIGN_FLAGS="--keychain %s"`,
		buildParams.BuildTool,
		projActionArg,
		buildParams.ProjectFile,
		buildParams.SchemeName,
		buildParams.DeviceDestination,
		buildParams.CodeSignIdentity,
		buildParams.KeychainName)
	if buildParams.ProvisioningProfile != "" {
		buildCmd = fmt.Sprintf(`%s PROVISIONING_PROFILE="%s"`,
			buildCmd,
			buildParams.ProvisioningProfile)
	}

	cargs := []string{
		"--login",
		"-c",
		buildCmd,
	}
	buildErr := RunCommandInDirWithArgsAndWriters(buildParams.ProjectFileDir, "bash", cargs, BuildLogWriter, BuildLogWriter)

	WriteLineToBuildLog("[[build-finished]]")
	return buildErr
}
