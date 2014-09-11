#!/bin/bash

#
# Removes the LaunchAgent plist file
# [!] Does NOT tries to stop the server - you have to do it yourself if you need it!
#

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "${THIS_SCRIPT_DIR}"
source _launchctl_common.sh

if [ -f "${server_plist_path}" ]; then
	rm "${server_plist_path}"
	if [ $? -ne 0 ]; then
		echo " [!] Failed to remove LaunchAgent plist file at ${server_plist_path}!"
		exit 1
	fi
else
	echo " (i) No LaunchAgent plist file found at path: ${server_plist_path}"
fi