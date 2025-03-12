#!/bin/bash
set -e

# Configuration
APP_NAME="Property Management"
DMG_NAME="PropertyManagement.dmg"


# Clean previous builds
rm -rf build/bin

# Build the application
echo "Building application for macOS..."
wails build -platform darwin/universal

# Create a DMG installer
echo "Creating DMG..."
mkdir -p /tmp/dmg-contents
cp -r "build/bin/$APP_NAME.app" /tmp/dmg-contents/
ln -s /Applications /tmp/dmg-contents/
hdiutil create -volname "$APP_NAME" -srcfolder /tmp/dmg-contents -ov -format UDZO "$DMG_NAME"
rm -rf /tmp/dmg-contents

echo "DMG created: $DMG_NAME"
echo "Done! The DMG file can now be distributed to macOS users."

