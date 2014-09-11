package main

import (
	"bufio"
	"errors"
	"log"
	"net/url"
	"os"
	"strings"
)

var (
	defaultBuildTool   = "xcodebuild"
	defaultDestination = "platform=iOS Simulator,name=iPad"
)

type BuildParams struct {
	ProjectFileDir      string
	BuildTool           string
	ProjectFile         string
	SchemeName          string
	DeviceDestination   string
	BuildOutputFilePath string
	// TODO: ATM supported only through config file - support from Query params too!
	CodeSignIdentity    string
	ProvisioningProfile string
	KeychainName        string
	KeychainPassword    string
}

func NewBuildParams() BuildParams {
	buildParams := BuildParams{}
	buildParams.BuildTool = defaultBuildTool
	buildParams.DeviceDestination = defaultDestination
	return buildParams
}

func FirstNotEmptyString(strs []string) string {
	for _, aStr := range strs {
		if aStr != "" {
			return aStr
		}
	}
	return ""
}

func MergeBuildParams(bp1, bp2 BuildParams) BuildParams {
	mbparams := NewBuildParams()

	mbparams.ProjectFileDir = FirstNotEmptyString([]string{
		bp1.ProjectFileDir, bp2.ProjectFileDir, mbparams.ProjectFileDir})
	mbparams.BuildTool = FirstNotEmptyString([]string{
		bp1.BuildTool, bp2.BuildTool, mbparams.BuildTool})
	mbparams.ProjectFile = FirstNotEmptyString([]string{
		bp1.ProjectFile, bp2.ProjectFile, mbparams.ProjectFile})
	mbparams.SchemeName = FirstNotEmptyString([]string{
		bp1.SchemeName, bp2.SchemeName, mbparams.SchemeName})
	mbparams.DeviceDestination = FirstNotEmptyString([]string{
		bp1.DeviceDestination, bp2.DeviceDestination, mbparams.DeviceDestination})
	mbparams.BuildOutputFilePath = FirstNotEmptyString([]string{
		bp1.BuildOutputFilePath, bp2.BuildOutputFilePath, mbparams.BuildOutputFilePath})
	//
	mbparams.CodeSignIdentity = FirstNotEmptyString([]string{
		bp1.CodeSignIdentity, bp2.CodeSignIdentity, mbparams.CodeSignIdentity})
	mbparams.ProvisioningProfile = FirstNotEmptyString([]string{
		bp1.ProvisioningProfile, bp2.ProvisioningProfile, mbparams.ProvisioningProfile})
	mbparams.KeychainName = FirstNotEmptyString([]string{
		bp1.KeychainName, bp2.KeychainName, mbparams.KeychainName})
	mbparams.KeychainPassword = FirstNotEmptyString([]string{
		bp1.KeychainPassword, bp2.KeychainPassword, mbparams.KeychainPassword})

	return mbparams
}

func (bparam BuildParams) Validate() error {
	if bparam.ProjectFileDir == "" {
		return errors.New("Parameter is not defined: projectdir")
	}
	if bparam.ProjectFile == "" {
		return errors.New("Parameter is not defined: projectfile")
	}
	if bparam.SchemeName == "" {
		return errors.New("Parameter is not defined: scheme")
	}
	if bparam.CodeSignIdentity == "" {
		return errors.New("Parameter is not defined: CodeSignIdentity")
	}
	if bparam.ProvisioningProfile == "" {
		return errors.New("Parameter is not defined: ProvisioningProfile")
	}
	if bparam.KeychainName == "" {
		return errors.New("Parameter is not defined: KeychainName")
	}
	if bparam.KeychainPassword == "" {
		return errors.New("Parameter is not defined: KeychainPassword")
	}
	return nil
}

func FirstString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func BuildParamsFromQueryValues(queryValues url.Values) BuildParams {
	bParams := NewBuildParams()
	bParams.BuildTool = FirstString(queryValues["buildtool"])
	bParams.ProjectFileDir = FirstString(queryValues["projectdir"])
	bParams.ProjectFile = FirstString(queryValues["projectfile"])
	bParams.SchemeName = FirstString(queryValues["scheme"])
	bParams.DeviceDestination = FirstString(queryValues["devicedestination"])
	bParams.BuildOutputFilePath = FirstString(queryValues["outputlogpath"])

	return bParams
}

func ReadBuildParamsFromConfigFile(confFilePth string) (BuildParams, error) {
	buildParams := NewBuildParams()
	file, err := os.Open(confFilePth)
	if err != nil {
		return buildParams, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, "=")
		if len(splits) >= 2 {
			aValue := strings.Join(splits[1:], "=")
			switch splits[0] {
			case "buildtool":
				buildParams.BuildTool = aValue
			case "projectdir":
				buildParams.ProjectFileDir = aValue
			case "projectfile":
				buildParams.ProjectFile = aValue
			case "scheme":
				buildParams.SchemeName = aValue
			case "devicedestination":
				buildParams.DeviceDestination = aValue
			case "outputlogpath":
				buildParams.BuildOutputFilePath = aValue
			case "code_sign_identity":
				buildParams.CodeSignIdentity = aValue
			case "provisioning_profile":
				buildParams.ProvisioningProfile = aValue
			case "keychain_name":
				buildParams.KeychainName = aValue
			case "keychain_password":
				buildParams.KeychainPassword = aValue
			default:
				log.Println("Invalid key - skipping line: ", line)
			}
		} else {
			log.Println("Skipping line: ", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return buildParams, err
	}

	return buildParams, nil
}
