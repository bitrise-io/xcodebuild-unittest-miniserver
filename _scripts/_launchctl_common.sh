#!/bin/bash

server_label_id="bitrise.tools.osx.xcuserver"
curr_user_lib_launch_agents_dir="$HOME/Library/LaunchAgents"
server_plist_path="${curr_user_lib_launch_agents_dir}/${server_label_id}.plist"
server_full_path="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd)/bin/osx/xcuserver"
server_logs_dir_path="${HOME}/logs"
server_log_file_path="${server_logs_dir_path}/${server_label_id}.log"