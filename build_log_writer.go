package main

import (
	"fmt"
	"io"
	"os"
)

var (
	BuildLogWriter io.Writer
	buildLogFile   *os.File
)

func OpenBuildLogWriter(buildParams BuildParams) error {
	if buildParams.BuildOutputFilePath != "" {
		outputfile, err := os.Create(buildParams.BuildOutputFilePath)
		if err != nil {
			return err
		}
		buildLogFile = outputfile
		BuildLogWriter = outputfile
		fmt.Println(" BuildLog writer opened with file: ", buildParams.BuildOutputFilePath)
	} else {
		BuildLogWriter = os.Stdout
		fmt.Println(" (!) No Build log file defined!")
		fmt.Println(" BuildLog writer opened STDOUT")
	}
	return nil
}

func WriteStringToBuildLog(s string) error {
	_, err := io.WriteString(BuildLogWriter, s)
	return err
}

func WriteLineToBuildLog(s string) error {
	return WriteStringToBuildLog(fmt.Sprintf("%s\n", s))
}

func CloseBuildLogWriter() error {
	if buildLogFile != nil {
		return buildLogFile.Close()
	}
	return nil
}
