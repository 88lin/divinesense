#!/bin/bash
# =============================================================================
# DivineSense Release Package Script
# =============================================================================
#
# Packages built binaries into distributable archives.
#
# Usage:
#   ./scripts/release/package-release.sh [version]
#
# Output:
#   releases/divinesense-<version>-<platform>.tar.gz
#
# =============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
DIST_DIR="${PROJECT_ROOT}/dist"
RELEASE_DIR="${PROJECT_ROOT}/releases"
VERSION="${1:-dev}"

# Logging
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Print banner
print_banner() {
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC}  ${GREEN}DivineSense Release Package${NC}                                ${BLUE}║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    log_info "Version: ${VERSION}"
    echo ""
}

# Create release package for a platform
create_package() {
    local binary=$1
    local platform=$2
    local package_name="divinesense-${VERSION}-${platform}.tar.gz"
    local package_path="${RELEASE_DIR}/${package_name}"
    local staging_dir="${RELEASE_DIR}/.staging/${platform}"

    log_info "Packaging ${platform}..."

    # Create staging directory
    rm -rf "${staging_dir}"
    mkdir -p "${staging_dir}"

    # Copy binary
    cp "${DIST_DIR}/${binary}" "${staging_dir}/divinesense"
    chmod +x "${staging_dir}/divinesense"

    # Copy service file
    cp "${SCRIPT_DIR}/divinesense.service" "${staging_dir}/"

    # Create directory structure script
    cat > "${staging_dir}/install.sh" << 'INSTALL_EOF'
#!/bin/bash
# Quick install script for extracted release

set -e

INSTALL_DIR="/opt/divinesense"
CONFIG_DIR="/etc/divinesense"
SERVICE_FILE="/etc/systemd/system/divinesense.service"

echo "Installing DivineSense..."

# Create user if not exists
if ! id -u divinesense &>/dev/null; then
    sudo useradd -r -s /bin/false -d /opt/divinesense divinesense
fi

# Create directories
sudo mkdir -p "${INSTALL_DIR}"/{bin,data,logs,backups}
sudo mkdir -p "${CONFIG_DIR}"

# Copy binary
sudo cp -f divinesense "${INSTALL_DIR}/bin/"
sudo chmod +x "${INSTALL_DIR}/bin/divinesense"

# Copy service file
sudo cp -f divinesense.service "${SERVICE_FILE}"

# Set ownership
sudo chown -R divinesense:divinesense "${INSTALL_DIR}"

echo "Installation complete!"
echo ""
echo "Next steps:"
echo "  1. Configure: sudo vi ${CONFIG_DIR}/config"
echo "  2. Start service: sudo systemctl enable --now divinesense"
echo "  3. Check status: sudo systemctl status divinesense"
INSTALL_EOF

    chmod +x "${staging_dir}/install.sh"

    # Create tarball
    cd "${staging_dir}"
    tar -czf "${package_path}" .

    # Cleanup
    rm -rf "${staging_dir}"

    log_success "Created ${package_name}"
}

# Main packaging process
main() {
    print_banner

    # Check if dist directory exists
    if [ ! -d "${DIST_DIR}" ]; then
        log_error "Distribution directory not found: ${DIST_DIR}"
        log_info "Run: ./scripts/release/build-release.sh ${VERSION}"
        exit 1
    fi

    # Create release directory
    rm -rf "${RELEASE_DIR}"
    mkdir -p "${RELEASE_DIR}"

    # Package each binary
    for binary in "${DIST_DIR}"/divinesense-*; do
        if [ -f "${binary}" ]; then
            local basename=$(basename "${binary}")
            # Extract platform from filename
            local platform=$(echo "${basename}" | sed "s/divinesense-${VERSION}-//")
            create_package "${basename}" "${platform}"
        fi
    done

    # Copy checksums
    if [ -f "${DIST_DIR}/checksums.txt" ]; then
        cp "${DIST_DIR}/checksums.txt" "${RELEASE_DIR}/"
    fi

    echo ""
    log_success "=========================================="
    log_success "Release packages created!"
    log_success "=========================================="
    echo ""
    log_info "Output directory: ${RELEASE_DIR}"
    echo ""
    log_info "Packages:"
    ls -lh "${RELEASE_DIR}"/*.tar.gz 2>/dev/null | awk '{printf "  %-50s %s\n", $9, $5}'
    echo ""
}

main "$@"
