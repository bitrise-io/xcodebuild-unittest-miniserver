package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"syscall"
)

var (
	serverPort          = "8081"
	okStatusMsg         = "ok"
	errorStatusMsg      = "error"
	endOfBuildLogMarker = "XCODEBUILDUNITTESTFINISHED"
)

type ResponseModel struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func unittestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request: ", r.URL)

	w.Header().Set("Content Type", "application/json")
	w.WriteHeader(http.StatusOK)

	queryValues := r.URL.Query()

	buildParamsFromConfig := BuildParams{}
	if configFilePath := FirstString(queryValues["configfile"]); configFilePath != "" {
		var confErr error
		buildParamsFromConfig, confErr = ReadBuildParamsFromConfigFile(configFilePath)
		if confErr != nil {
			log.Println("Failed to read config from file: ", configFilePath)
		}
	}

	queryBuildParams := BuildParamsFromQueryValues(queryValues)
	buildParams := MergeBuildParams(queryBuildParams, buildParamsFromConfig)
	fmt.Printf("Merged buildParams: %#v\n", buildParams)

	err := OpenBuildLogWriter(buildParams)
	if err == nil {
		defer CloseBuildLogWriter()

		WriteLineToBuildLog(fmt.Sprintf(" (i) Using Build Params: %#v", buildParams))
		err = buildParams.Validate()
		if err == nil {
			err = ExecuteBuildWithParams(buildParams)
			if err != nil {
				// Did the command fail because of an unsuccessful exit code
				var waitStatus syscall.WaitStatus
				if exitError, ok := err.(*exec.ExitError); ok {
					waitStatus = exitError.Sys().(syscall.WaitStatus)
					exCode := waitStatus.ExitStatus()
					fmt.Println("Exit status: ", exCode)
					if exCode == 65 {
						WriteLineToBuildLog("=> [!] Error with code 65 (Unable to run app in Simulator) - retrying")
						err = ExecuteBuildWithParams(buildParams)
					}
				}
			}
		}
	}

	//
	// Response
	statusMsg := okStatusMsg
	respMsg := "Test finished with success"
	if err != nil {
		log.Println("Error: ", err)
		WriteLineToBuildLog(fmt.Sprintf("[!] Error: %s", err))
		statusMsg = errorStatusMsg
		respMsg = fmt.Sprintf("%s", err)
	}
	//
	respModel := ResponseModel{
		Status: statusMsg,
		Msg:    respMsg,
	}

	WriteLineToBuildLog(fmt.Sprintf("%s: %s", endOfBuildLogMarker, statusMsg))
	WriteLineToBuildLog("-> Build Finished")
	if err := json.NewEncoder(w).Encode(&respModel); err != nil {
		log.Println("Error: ", err)
	}
}

func main() {
	http.HandleFunc("/unittest", unittestHandler)
	fmt.Println("Ready to serve on port:", serverPort)
	fmt.Println()
	http.ListenAndServe(":"+serverPort, nil)
}
