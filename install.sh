#!/bin/sh
set -e

# Flashdoc installer script
# Usage: curl -sSL https://raw.githubusercontent.com/heidene/flashdoc/main/install.sh | sh

REPO="heidene/flashdoc"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Darwin)
        OS="darwin"
        ;;
    Linux)
        OS="linux"
        ;;
    MINGW* | MSYS* | CYGWIN*)
        OS="windows"
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64 | amd64)
        ARCH="x86_64"
        ;;
    arm64 | aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected OS: $OS"
echo "Detected Architecture: $ARCH"

# Get latest release info
RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"
echo "Fetching latest release..."

if command -v curl > /dev/null 2>&1; then
    RELEASE_JSON=$(curl -sL "$RELEASE_URL")
elif command -v wget > /dev/null 2>&1; then
    RELEASE_JSON=$(wget -qO- "$RELEASE_URL")
else
    echo "Error: curl or wget is required"
    exit 1
fi

# Extract version and download URL
VERSION=$(echo "$RELEASE_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
    echo "Error: Could not fetch latest version"
    exit 1
fi

echo "Latest version: $VERSION"

# Construct download URL
ARCHIVE="flashdoc_${VERSION#v}_${OS}_${ARCH}"
if [ "$OS" = "windows" ]; then
    ARCHIVE="${ARCHIVE}.zip"
    BINARY="flashdoc.exe"
else
    ARCHIVE="${ARCHIVE}.tar.gz"
    BINARY="flashdoc"
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE"
echo "Downloading from: $DOWNLOAD_URL"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

cd "$TMP_DIR"

# Download
if command -v curl > /dev/null 2>&1; then
    curl -sLO "$DOWNLOAD_URL"
elif command -v wget > /dev/null 2>&1; then
    wget -q "$DOWNLOAD_URL"
fi

# Extract
echo "Extracting archive..."
if [ "$OS" = "windows" ]; then
    unzip -q "$ARCHIVE"
else
    tar -xzf "$ARCHIVE"
fi

# Install
echo "Installing to $INSTALL_DIR..."

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "$INSTALL_DIR/flashdoc"
    chmod +x "$INSTALL_DIR/flashdoc"
else
    echo "Note: Installing to $INSTALL_DIR requires sudo"
    sudo mv "$BINARY" "$INSTALL_DIR/flashdoc"
    sudo chmod +x "$INSTALL_DIR/flashdoc"
fi

echo ""
echo "✅ Flashdoc installed successfully!"
echo ""
echo "Run 'flashdoc --version' to verify installation"
echo "Run 'flashdoc --help' to see usage"
