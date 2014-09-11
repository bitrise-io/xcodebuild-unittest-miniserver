xcodebuild-unittest-miniserver (XCUServer)
==============================

Mini server (written in Go) to perform Xcode Unit Tests through the command line / terminal.


# The problem

Xcode Unit Tests require the iPhone/iOS Simulator to run, and the xcodebuild tests have to be run through a logged in user in a GUI context.

It's a real pain to integrate into automation systems which perform commands through SSH.


# The solution

You have to start this simple server from a logged in user from the GUI (for example: from the Terminal app). You can automate this through *launchctl* so when the user logs in the server will start automatically.

When the server is running you can issue Xcode unit test commands through it's web interface, for example with *curl*.


# Run and Install

You can run the server without installing it (useful for testing):

* You can simply build and run the server
	* You can use the *_scripts/build.sh* script to test and build a binary, then move it to *bin/osx/xcuserver*
* You can run the prepared binary from *bin/osx/xcuserver*

You can install (register) the server for OS X LaunchControl for the current user with the *_scripts/install_launchctl_plist_for_current_user.sh* script. This will start the server every time the user logs in (on the GUI).

Once the server is running you can access it at port 8081 (you can configure the port in *main.go*).
A curl example:
	
	$ curl http://localhost:8081/unittest?configfile=[build_config_file_absolute_path]

The result of the call will be a JSON response with the format:

	{"status": "ok OR error", "msg": "Error or success description"}

The end of the *build output log* file will be marked with a result as well.
For a successful unit test:
	
	XCODEBUILDUNITTESTFINISHED: ok

For a failed unit test:

	XCODEBUILDUNITTESTFINISHED: error


# Params

Parameters can be defined:

* through URL query parameters (code signing parameters not yet supported) - you'll have to properly URL encode your parameters!
* a config file, which's path is passed as a URL query parameter

## Supported URL Query Parameters

* buildtool : xcodebuild or xctool
* projectdir : root directory path of the project
* projectfile : relative path of the project (or workspace) file (relative to the projectdir foler)
* scheme : the shared (!) scheme in the project you want to build
* devicedestination : value for xcodebuild's destination parameter. An example: platform=iOS Simulator,name=iPad
* outputlogpath : path to a log file. The build's output will be redirected into this file.
* configfile : path to a config file which contains the build parameters (see the *Supported Config File Parameters* section below)

## Supported Config File Parameters

A config file can be specified with the *configfile* query parameter.
Example:

	$ curl http://localhost:8081/unittest?configfile=[build_config_file_absolute_path]

The config file contains the parameters in the following format:

* one line, one key-value pair
* key-value pairs separated with an equal (=) sign

Example:

	buildtool=xcodebuild
	devicedestination=platform=iOS Simulator,name=iPad
	scheme=YouSharedScheme

The config file supports all the *URL Query Parameters* except the *configfile* option (of course) and additionaly it supports the following parameters:

* code_sign_identity : the identity from the certificate file you want to use, for example: iPhone Developer: Bitrise Dev (C9NJZ996T5)
	* it have to be already loaded, into the keychain you want to use for signing.
* provisioning_profile : the provisioning profile ID you want to use. Example: 323BFB95-3D5D-4BCC-9288-CEF43694AA5D
* keychain_name : the loaded keychain's name. Example: bitrise.keychain
* keychain_password : the specified keychain's password.

> For more information on Keychain handling check the **Keychain handling** section*


## Parameter priority

The parameters you specify as *URL Query Parameters* are handled with higher priority than the ones you specify in your *configfile*.

If you specify a *configfile* and also specify some *URL Query Parameters* the parameters you specify as *URL Query Parameters* will be used if the same parameter is defined in the *configfile*.


# Dependencies

The server itself doesn't have any dependency, the binary you can find at *bin/osx/xcuserver* can be started without any configuration.

For performing actual builds you have to install the following dependencies:

* Xcode and Xcode Command Line Tools
* if you want to use xctool as the *buildtool* you have to install it


# Keychain handling

Keychain handling not (yet?) included: for code signing you'll have to prepare the keychain for XCUServer. For an example you can check the [Bitrise Xcode Builder Step's code](https://github.com/bitrise-io/steps-xcode-builder) - specifically the *keychain.sh* utility script.


# Example usage in a real project

You can check how we at [Bitrise](http://www.bitrise.io) use XCUServer.

XCUServer is included in the Bitrise OS X VMs, part of the [OS X Box Bootstrap repo](https://github.com/bitrise-io/osx-box-bootstrap).
The install script can be found [here](https://github.com/bitrise-io/osx-box-bootstrap/blob/master/installers/install_xcuserver.sh)

XCUServer is used in the [Bitrise Xcode Builder Step](https://github.com/bitrise-io/steps-xcode-builder) - this Bitrise Step runs it's unit tests with XCUServer.


# TODO

* support all the Build Parameters through *URL Query Parameters* (the ones which are currently only supported through the config file)
* print version number if executed with -version parameter.
* replace the *install_launchctl_plist_for_current_user* script with a built in method, when executed with -install parameter.
* handle Keychain (creation, profile imports, etc.)
