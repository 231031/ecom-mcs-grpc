#!/bin/sh

# Replace __FOLDER_NAME__ with the actual environment variable value
sed -i "s|__FOLDER_NAME__|$FOLDER_NAME|g" .air.toml

# Run Air
exec air -c .air.toml
