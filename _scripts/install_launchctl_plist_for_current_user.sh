#!/bin/bash

#
# Run this script from a logged in user - with the user you want to run the Xcode Unit Tests with!
#
#  For launchctr related configs check the _launchctl_common.sh file
#

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "${THIS_SCRIPT_DIR}"
source _launchctl_common.sh


echo " (i) curr_user_lib_launch_agents_dir: ${curr_user_lib_launch_agents_dir}"
mkdir -p "${curr_user_lib_launch_agents_dir}"
if [ $? -ne 0 ]; then
	echo " [!] Failed to create the required LaunchAgents dir at ${curr_user_lib_launch_agents_dir}!"
	exit 1
fi

echo " (i) server_full_path: ${server_full_path}"
if [ ! -f "${server_full_path}" ]; then
	echo " [!] Server full path is invalid - server not found at path: ${server_full_path}"
	exit 1
fi

echo " (i) server_logs_dir_path: ${server_logs_dir_path}"
echo " (i) server_log_file_path: ${server_log_file_path}"
mkdir -p "${server_logs_dir_path}"
if [ $? -ne 0 ]; then
	echo " [!] Failed to create the required 'logs' dir at ${server_logs_dir_path}!"
	exit 1
fi

echo " (i) server_plist_path: ${server_plist_path}"

cat >"${server_plist_path}" <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>${server_label_id}</string>
    <key>ProgramArguments</key>
    <array>
        <string>${server_full_path}</string>
    </array>
    <key>StandardOutPath</key>
    <string>${server_log_file_path}</string>
    <key>StandardErrorPath</key>
    <string>${server_log_file_path}</string>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOL
if [ $? -ne 0 ]; then
	echo " [!] Failed to write LaunchAgent plist to path: ${server_plist_path}"
	exit 1
fi

echo " (i) LaunchAgent plist content:"
cat "${server_plist_path}"
if [ $? -ne 0 ]; then
	echo " [!] Failed to read LaunchAgent plist from path: ${server_plist_path}"
	exit 1
fi


echo
echo "==> INSTALL SUCCESS"
echo " * LaunchAgent plist saved to path: ${server_plist_path}"
echo " * You can start (or restart) the server with the reload_server_with_launchctl.sh script"
echo " * You can start the server with: launchctl load "${server_plist_path}""
echo " * You can stop the server with: launchctl unload "${server_plist_path}""

