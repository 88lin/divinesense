#!/bin/bash
# =============================================================================
# DivineSense 阿里云 2C2G 一键安装脚本 v3.0
# =============================================================================
#
# 使用方式:
#   curl -fsSL https://raw.githubusercontent.com/hrygo/divinesense/main/deploy/aliyun/install.sh | bash
#
# 选择部署模式:
#   curl -fsSL ... | bash -s -- --mode=docker   # Docker 模式 (默认)
#   curl -fsSL ... | bash -s -- --mode=binary  # 二进制模式 (Geek Mode 原生支持)
#
# 支持系统:
#   - 阿里云 Linux 2/3
#   - CentOS 7/8
#   - Rocky Linux 8/9
#   - Ubuntu 18.04+
#   - Debian 10+
#
# =============================================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 默认配置
DEPLOY_MODE="${DEPLOY_MODE:-docker}"  # docker | binary
INSTALL_DIR="${INSTALL_DIR:-/opt/divinesense}"
CONFIG_DIR="${CONFIG_DIR:-/etc/divinesense}"
REPO_URL="${REPO_URL:-https://github.com/hrygo/divinesense.git}"
BRANCH="${BRANCH:-main}"
BACKUP_DIR="${INSTALL_DIR}/backups"
DOWNLOAD_URL="${DOWNLOAD_URL:-https://github.com/hrygo/divinesense/releases}"
VERSION="${VERSION:-latest}"

# 系统要求
MIN_RAM_MB=1800
MIN_DISK_MB=4096

# 网络超时设置
CURL_CONNECT_TIMEOUT=30
CURL_MAX_TIME=300

# 日志函数
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_step() { echo -e "${CYAN}[STEP]${NC} $1"; }

# 显示使用帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --mode=docker   Docker 容器部署模式 (默认)"
    echo "  --mode=binary  二进制部署模式 (推荐 Geek Mode)"
    echo "  --version=V    指定安装版本 (仅 binary 模式)"
    echo "  -h, --help     显示此帮助信息"
    echo ""
    echo "部署模式对比:"
    echo "  Docker:  资源占用较高，启动较慢，需要额外配置 Geek Mode"
    echo "  Binary:  资源占用低，启动快，原生支持 Geek Mode"
    echo ""
    exit 0
}

# 解析参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --mode=*)
                DEPLOY_MODE="${1#*=}"
                shift
                ;;
            --mode)
                DEPLOY_MODE="$2"
                shift 2
                ;;
            --version=*)
                VERSION="${1#*=}"
                shift
                ;;
            -h|--help)
                show_help
                ;;
            *)
                shift
                ;;
        esac
    done

    # 验证模式
    if [[ "$DEPLOY_MODE" != "docker" && "$DEPLOY_MODE" != "binary" ]]; then
        log_error "无效的部署模式: $DEPLOY_MODE"
        log_info "支持的模式: docker, binary"
        exit 1
    fi
}

# 打印 Banner
print_banner() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  ${GREEN}DivineSense 一键部署脚本 v3.0${NC}                             ${CYAN}║${NC}"
    echo -e "${CYAN}║${NC}  ${GREEN}适用于阿里云 2C2G 服务器${NC}                                   ${CYAN}║${NC}"
    echo -e "${CYAN}╠════════════════════════════════════════════════════════════╣${NC}"
    echo -e "${CYAN}║${NC}  部署模式: ${YELLOW}${DEPLOY_MODE^^}${NC}                                        ${CYAN}║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# 检查系统资源
check_system_resources() {
    log_step "检查系统资源..."

    # 检查内存
    local total_mem_kb=$(grep MemTotal /proc/meminfo 2>/dev/null | awk '{print $2}')
    local total_mem_mb=$((total_mem_kb / 1024))

    if [ "$total_mem_mb" -lt "$MIN_RAM_MB" ]; then
        log_warn "内存不足: 需要 ${MIN_RAM_MB}MB，当前 ${total_mem_mb}MB"
        log_info "建议升级配置或使用 swap 空间"
    else
        log_success "内存检查通过: ${total_mem_mb}MB"
    fi

    # 检查磁盘空间
    local available_disk_mb=$(df -m /opt | awk 'NR==2 {print $4}')

    if [ "$available_disk_mb" -lt "$MIN_DISK_MB" ]; then
        log_error "磁盘空间不足: 需要 ${MIN_DISK_MB}MB，当前可用 ${available_disk_mb}MB"
        exit 1
    fi
    log_success "磁盘检查通过: ${available_disk_mb}MB 可用"
}

# 检测系统
detect_os() {
    log_step "检测操作系统..."

    # 检查阿里云 Linux
    if [ -f /etc/aliyun-release ]; then
        OS="aliyun"
        . /etc/aliyun-release 2>/dev/null || true
        OS_VERSION="${VERSION_ID:-unknown}"
        PKG_MANAGER="yum"
        log_info "检测到阿里云 Linux: $OS_VERSION"
        return 0
    fi

    # 检查标准 os-release
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        OS_VERSION=$VERSION_ID
    else
        # 兼容老系统
        if [ -f /etc/redhat-release ]; then
            OS="centos"
            OS_VERSION=$(rpm -qf /etc/redhat-release --queryformat '%{VERSION}' | cut -d. -f1)
        elif [ -f /etc/debian_version ]; then
            OS="debian"
            OS_VERSION=$(cat /etc/debian_version)
        else
            log_error "无法检测操作系统"
            exit 1
        fi
    fi

    case "$OS" in
        alpine|arch|manjaro)
            PKG_MANAGER="apk"
            ;;
        debian|ubuntu|linuxmint)
            PKG_MANAGER="apt"
            ;;
        centos|rhel|fedora|rocky|almalinux|aliyun)
            PKG_MANAGER="yum"
            ;;
        *)
            log_error "不支持的操作系统: $OS"
            log_info "支持的系统: 阿里云 Linux, CentOS, Rocky, Debian, Ubuntu"
            exit 1
            ;;
    esac

    log_success "系统: $OS $OS_VERSION | 包管理器: $PKG_MANAGER"
}

# 检测架构
detect_arch() {
    log_step "检测系统架构..."

    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)
            BINARY_ARCH="amd64"
            ;;
        aarch64|arm64)
            BINARY_ARCH="arm64"
            ;;
        *)
            log_error "不支持的架构: $ARCH"
            exit 1
            ;;
    esac

    log_success "架构: $ARCH (二进制: $BINARY_ARCH)"
}

# 检查是否为 root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "此脚本需要 root 权限运行"
        log_info "请使用: sudo $0"
        exit 1
    fi
}

# 安装基础工具
install_base_tools() {
    log_step "安装基础工具..."

    local tools="curl wget openssl git"

    case "$PKG_MANAGER" in
        apt)
            export DEBIAN_FRONTEND=noninteractive
            apt-get update -qq
            apt-get install -y -qq $tools 2>/dev/null || true
            ;;
        yum)
            yum install -y -q $tools 2>/dev/null || true
            ;;
        apk)
            apk add --no-cache $tools
            ;;
    esac

    log_success "基础工具已安装"
}

generate_password() {
    if command -v openssl &>/dev/null; then
        openssl rand -hex 16 | head -c 20
    else
        tr -dc A-Za-z0-9 </dev/urandom 2>/dev/null | head -c 20
    fi
}

# 获取服务器 IP
get_server_ip() {
    # 阿里云元数据服务（最快）
    local ip=$(curl -s --connect-timeout 1 http://100.100.100.200/latest/meta-data/network/interfaces/macs/ 2>/dev/null | head -1 | \
              xargs -I {} curl -s http://100.100.100.200/latest/meta-data/network/interfaces/{}/ipv4/primary-ip-address 2>/dev/null)

    if [ -z "$ip" ]; then
        ip=$(curl -s --connect-timeout 3 -4 ifconfig.me 2>/dev/null)
    fi
    if [ -z "$ip" ]; then
        ip=$(curl -s --connect-timeout 3 -4 icanhazip.com 2>/dev/null)
    fi
    if [ -z "$ip" ]; then
        ip=$(hostname -I | awk '{print $1}')
    fi

    echo "$ip"
}

# =============================================================================
# Docker 模式安装函数
# =============================================================================

# 安装 Docker
install_docker() {
    log_step "安装 Docker..."

    if command -v docker &>/dev/null; then
        local installed_version=$(docker --version | grep -oP '\d+\.\d+\.\d+' || echo "unknown")
        log_success "Docker 已安装: $installed_version"
        return 0
    fi

    case "$PKG_MANAGER" in
        apt)
            # Ubuntu/Debian
            if [ ! -f /usr/share/keyrings/docker-archive-keyring.gpg ]; then
                install -m 0755 -d /etc/apt/keyrings
                curl -fsSL --connect-timeout $CURL_CONNECT_TIMEOUT --max-time $CURL_MAX_TIME \
                    https://download.docker.com/linux/${OS}/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
                chmod a+r /etc/apt/keyrings/docker.gpg

                echo \
                  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/${OS} \
                  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
                  tee /etc/apt/sources.list.d/docker.list > /dev/null

                apt-get update -qq
            fi
            apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-buildx-plugin
            ;;
        yum)
            # CentOS/RHEL/Aliyun Linux
            yum install -y -q yum-utils
            yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
            yum install -y -q docker-ce docker-ce-cli containerd.io docker-buildx-plugin
            ;;
        apk)
            apk add docker docker-cli-compose
            ;;
    esac

    # 启动 Docker
    systemctl enable docker 2>/dev/null || true
    systemctl start docker

    log_success "Docker 安装完成"
}

# 安装 Docker Compose
install_docker_compose() {
    log_step "安装 Docker Compose..."

    if docker compose version &>/dev/null; then
        local compose_version=$(docker compose version --short)
        log_success "Docker Compose 已安装: $compose_version"
        return 0
    fi

    # Docker Compose v2 随 Docker 已安装
    if ! docker compose version &>/dev/null; then
        log_info "安装独立 docker-compose..."
        curl -fsSL --connect-timeout $CURL_CONNECT_TIMEOUT --max-time $CURL_MAX_TIME \
            -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
            -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
    fi

    log_success "Docker Compose 已就绪"
}

# 配置镜像加速
setup_docker_mirror() {
    log_step "配置 Docker 镜像加速..."

    local docker_config_dir="/etc/docker"
    local daemon_config="$docker_config_dir/daemon.json"

    mkdir -p "$docker_config_dir"

    # 备份现有配置
    if [ -f "$daemon_config" ]; then
        cp "$daemon_config" "${daemon_config}.backup.$(date +%Y%m%d%H%M%S)"
    fi

    cat > "$daemon_config" << 'EOF'
{
  "registry-mirrors": [
    "https://docker.m.daocloud.io",
    "https://dockerproxy.com",
    "https://docker.mirrors.ustc.edu.cn",
    "https://docker.nju.edu.cn"
  ],
  "max-concurrent-downloads": 10,
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  }
}
EOF

    # 重启 Docker
    systemctl restart docker 2>/dev/null || service docker restart 2>/dev/null || true

    log_success "镜像加速已配置"
}

# 克隆仓库 (Docker 模式)
clone_repo() {
    log_step "下载 DivineSense 部署文件..."

    cd "$INSTALL_DIR"

    if [ -d ".git" ]; then
        log_info "更新现有仓库..."
        git pull origin "$BRANCH" 2>/dev/null || true
    else
        log_info "克隆仓库..."
        git clone -b "$BRANCH" --depth 1 "$REPO_URL" "$INSTALL_DIR" 2>/dev/null || {
            log_error "Git 克隆失败，尝试下载发布包..."
            wget -O /tmp/divinesense.tar.gz \
                --connect-timeout=$CURL_CONNECT_TIMEOUT \
                --timeout=$CURL_MAX_TIME \
                "https://github.com/hrygo/divinesense/archive/refs/heads/main.tar.gz" 2>/dev/null || {
                log_error "下载失败，请检查网络连接"
                exit 1
            }
            tar -xzf /tmp/divinesense.tar.gz -C "$INSTALL_DIR" --strip-components=1
            rm -f /tmp/divinesense.tar.gz
        }
    fi

    log_success "部署文件已下载"
}

# 生成 Docker 模式配置文件
generate_docker_env_file() {
    log_step "生成配置文件..."

    local env_file="$INSTALL_DIR/.env.prod"
    local template_file="$INSTALL_DIR/deploy/aliyun/.env.prod.example"
    local db_password=$(generate_password)
    local server_ip=$(get_server_ip)

    if [ -z "$server_ip" ]; then
        log_warn "无法获取公网 IP，请手动配置 INSTANCE_URL"
        server_ip="your-server-ip"
    fi

    # 1. 复制模板文件
    if [ -f "$template_file" ]; then
        cp "$template_file" "$env_file"
    else
        log_warn "模板文件未找到: $template_file"
        log_info "尝试下载模板..."
        curl -fsSL "${REPO_URL%.git}/raw/${BRANCH}/deploy/aliyun/.env.prod.example" -o "$env_file" || {
             log_error "下载配置文件模板失败"
             exit 1
        }
    fi

    # 2. 替换关键变量 (使用兼容的 sed 写法)
    # 替换 IP
    sed -i "s|DIVINESENSE_INSTANCE_URL=.*|DIVINESENSE_INSTANCE_URL=http://${server_ip}:5230|g" "$env_file"
    
    # 替换数据库密码 (注意转义)
    # 为了安全和避免 sed 分隔符冲突，这里简化处理，假设密码没有特殊字符，或者使用 | 分隔
    sed -i "s|POSTGRES_PASSWORD=.*|POSTGRES_PASSWORD=${db_password}|g" "$env_file"

    # 3. 补充 Docker 特有配置 (如果模板中没有包含)
    # 模板中已经包含了大部分公共配置
    
    log_success "配置文件已生成: $env_file"

    # 保存密码到单独文件 (不输出到日志)
    echo "$db_password" > "$INSTALL_DIR/.db_password"
    chmod 600 "$INSTALL_DIR/.db_password"
    log_info "数据库密码已保存到: $INSTALL_DIR/.db_password"
}
    log_success "配置文件已生成: $env_file"

    # 保存密码到单独文件 (不输出到日志)
    echo "$db_password" > "$INSTALL_DIR/.db_password"
    chmod 600 "$INSTALL_DIR/.db_password"
    log_info "数据库密码已保存到: $INSTALL_DIR/.db_password"
}

# 拉取镜像
pull_images() {
    log_step "拉取 Docker 镜像..."

    # 拉取 PostgreSQL
    log_info "拉取 PostgreSQL + pgvector..."
    docker pull pgvector/pgvector:pg16 || {
        log_error "镜像拉取失败"
        log_info "尝试配置镜像加速..."
        setup_docker_mirror
        docker pull pgvector/pgvector:pg16
    }

    # 拉取 DivineSense
    log_info "拉取 DivineSense..."
    docker pull ghcr.io/hrygo/divinesense:latest || {
        log_warn "官方镜像可能不存在，跳过..."
    }

    # 标记为本地镜像名
    if docker images | grep -q "ghcr.io/hrygo/divinesense"; then
        docker tag ghcr.io/hrygo/divinesense:latest divinesense:latest 2>/dev/null || true
    fi

    log_success "镜像准备完成"
}

# 部署 Docker 服务
deploy_docker_services() {
    log_step "部署服务..."

    cd "$INSTALL_DIR"

    # 确保存在必要的目录结构
    mkdir -p "$INSTALL_DIR/docker/compose"

    # 检查 compose 文件
    if [ ! -f "$INSTALL_DIR/docker/compose/prod.yml" ]; then
        log_error "缺少部署文件，请重新运行"
        exit 1
    fi

    # 使用 docker compose 启动服务
    if docker compose version &>/dev/null; then
        docker compose -f docker/compose/prod.yml --env-file .env.prod up -d
    else
        docker-compose -f docker/compose/prod.yml --env-file .env.prod up -d
    fi

    log_success "服务启动完成"
}

# =============================================================================
# 二进制模式安装函数
# =============================================================================

# 下载并验证二进制
download_binary() {
    log_step "下载 DivineSense 二进制..."

    local binary_name="divinesense-${VERSION}-linux-${BINARY_ARCH}"
    local download_url="${DOWNLOAD_URL}/download/${VERSION}/${binary_name}"
    local tmp_file="/tmp/divinesense.tmp"
    local checksum_file="/tmp/divinesense.sha256"

    if [ "$VERSION" = "latest" ]; then
        download_url="${DOWNLOAD_URL}/latest/download/${binary_name}"
    fi

    log_info "下载地址: $download_url"

    # 下载二进制
    if ! curl -fsSL \
        --connect-timeout $CURL_CONNECT_TIMEOUT \
        --max-time $CURL_MAX_TIME \
        "$download_url" -o "$tmp_file"; then
        log_error "二进制下载失败"
        log_info "请检查网络或手动下载: $DOWNLOAD_URL"
        exit 1
    fi

    # 下载校验和
    if ! curl -fsSL \
        --connect-timeout $CURL_CONNECT_TIMEOUT \
        --max-time $CURL_MAX_TIME \
        "${download_url}.sha256" -o "$checksum_file"; then
        log_warn "校验文件下载失败，跳过验证"
    else
        # 验证校验和
        cd /tmp
        local expected_checksum=$(cat "$checksum_file" | cut -d' ' -f1)
        local actual_checksum=$(sha256sum "$tmp_file" | cut -d' ' -f1)

        if [ "$expected_checksum" != "$actual_checksum" ]; then
            log_error "校验和不匹配!"
            log_error "预期: $expected_checksum"
            log_error "实际: $actual_checksum"
            rm -f "$tmp_file" "$checksum_file"
            exit 1
        fi
        log_success "校验和验证通过"
    fi

    # 安装二进制
    mv "$tmp_file" "${INSTALL_DIR}/bin/divinesense"
    chmod +x "${INSTALL_DIR}/bin/divinesense"
    rm -f "$checksum_file"

    log_success "二进制已安装"
}

# 创建用户和目录 (二进制模式)
setup_binary_user_dirs() {
    log_step "创建用户和目录..."

    # 创建用户
    if ! id -u divinesense &>/dev/null; then
        useradd -r -s /bin/false -d "$INSTALL_DIR" -c "DivineSense Service" divinesense
        log_success "创建用户: divinesense"
    else
        log_info "用户已存在: divinesense"
    fi

    # 创建目录
    mkdir -p "${INSTALL_DIR}"/{bin,data,logs,backups,docker}
    mkdir -p "$CONFIG_DIR"

    log_success "目录已创建"
}

# 生成二进制模式配置
generate_binary_config() {
    log_step "生成配置文件..."

    local config_file="${CONFIG_DIR}/config"
    local db_password=$(generate_password)
    local server_ip=$(get_server_ip)

    if [ -z "$server_ip" ]; then
        server_ip="your-server-ip"
    fi
    
    # 1. 下载模板文件
    local template_url="${REPO_URL%.git}/raw/${BRANCH}/deploy/aliyun/config.binary.example"
    # 处理 github.com -> raw.githubusercontent.com 转换 (如果 REPO_URL 是常规 github 链接)
    if [[ "$REPO_URL" == *"github.com"* ]]; then
        template_url="${REPO_URL/github.com/raw.githubusercontent.com}"
        template_url="${template_url%.git}/${BRANCH}/deploy/aliyun/config.binary.example"
    fi

    log_info "下载配置模板: $template_url"
    if ! curl -fsSL --connect-timeout $CURL_CONNECT_TIMEOUT --max-time $CURL_MAX_TIME "$template_url" -o "$config_file"; then
         log_error "下载配置模板失败"
         log_info "尝试生成默认配置..."
         # 回退逻辑: 写入最小化配置 (或者直接退出，取决于策略。为了健壮性，这里可以写入一个极简版，但为了 SSOT，也许失败更好？)
         # 这里选择退出，因为没有模板就无法保证配置的正确性
         exit 1
    fi

    # 2. 替换变量
    # 替换 IP
    sed -i "s|DIVINESENSE_INSTANCE_URL=.*|DIVINESENSE_INSTANCE_URL=http://${server_ip}:5230|g" "$config_file"
    
    # 替换数据库密码 (DSN 中)
    sed -i "s|postgres://divinesense:your_secure_password|postgres://divinesense:${db_password}|g" "$config_file"
    
    # 确保安装目录正确 (如果模板中默认值不对)
    sed -i "s|DIVINESENSE_DATA=.*|DIVINESENSE_DATA=${INSTALL_DIR}/data|g" "$config_file"
    sed -i "s|DIVINESENSE_CLAUDE_CODE_WORKDIR=.*|DIVINESENSE_CLAUDE_CODE_WORKDIR=${INSTALL_DIR}/data|g" "$config_file"

    chmod 640 "$config_file"

    log_success "配置已生成: $config_file"

    # 保存密码到单独文件
    echo "$db_password" > "${CONFIG_DIR}/.db_password"
    chmod 600 "${CONFIG_DIR}/.db_password"
    log_info "数据库密码已保存到: ${CONFIG_DIR}/.db_password"
}

# 安装 systemd 服务 (二进制模式)
install_binary_service() {
    log_step "安装 systemd 服务..."

    cat > "/etc/systemd/system/divinesense.service" << 'EOF'
[Unit]
Description=DivineSense AI-Powered Personal Second Brain
Documentation=https://github.com/hrygo/divinesense
After=network-online.target
Wants=network-online.target

[Service]
Type=exec
User=divinesense
Group=divinesense
WorkingDirectory=/opt/divinesense/data
EnvironmentFile=-/etc/divinesense/config
ExecStart=/opt/divinesense/bin/divinesense
Restart=always
RestartSec=10s
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/divinesense/data /opt/divinesense/logs /var/log
LimitNOFILE=65536
StandardOutput=journal
StandardError=journal
SyslogIdentifier=divinesense

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    log_success "服务已安装"
}

# 设置 PostgreSQL Docker (二进制模式)
setup_binary_postgres() {
    log_step "设置 PostgreSQL..."

    if ! command -v docker &>/dev/null || ! docker info &>/dev/null; then
        log_warn "Docker 不可用，PostgreSQL 需要手动配置"
        log_info "选项:"
        log_info "  1. 安装 Docker: curl -fsSL https://get.docker.com | sh"
        log_info "  2. 使用系统 PostgreSQL"
        log_info "  3. 使用 SQLite (无 AI 功能)"
        return 0
    fi

    # 读取密码
    local db_password=$(cat "${CONFIG_DIR}/.db_password" 2>/dev/null)

    # 创建 docker-compose 文件
    cat > "${INSTALL_DIR}/docker/postgres.yml" << EOF
version: '3.8'
services:
  postgres:
    image: pgvector/pgvector:pg16
    container_name: divinesense-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: divinesense
      POSTGRES_USER: divinesense
      POSTGRES_PASSWORD: \${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "25432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U divinesense"]
      interval: 10s
      timeout: 5s
      retries: 5
    logs:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  postgres_data:
EOF

    # 创建环境文件
    cat > "${INSTALL_DIR}/docker/.env" << EOF
POSTGRES_PASSWORD=${db_password}
EOF
    chmod 600 "${INSTALL_DIR}/docker/.env"

    # 启动 PostgreSQL
    cd "${INSTALL_DIR}/docker"
    if docker compose version &>/dev/null; then
        docker compose -f postgres.yml up -d
    else
        docker-compose -f postgres.yml up -d
    fi

    # 等待 PostgreSQL 就绪
    local max_wait=60
    local waited=0
    while [ $waited -lt $max_wait ]; do
        if docker exec divinesense-postgres pg_isready -U divinesense &>/dev/null; then
            # 启用 pgvector
            docker exec divinesense-postgres psql -U divinesense -d divinesense \
                -c "CREATE EXTENSION IF NOT EXISTS vector;" &>/dev/null || true
            log_success "PostgreSQL 已就绪"
            return 0
        fi
        sleep 2
        waited=$((waited + 2))
    done

    log_warn "PostgreSQL 启动可能需要更长时间"
}

# 设置二进制权限
set_binary_permissions() {
    log_step "设置权限..."

    chown -R divinesense:divinesense "${INSTALL_DIR}"
    chmod 755 "${INSTALL_DIR}/bin/divinesense"

    log_success "权限已设置"
}

# =============================================================================
# 通用函数
# =============================================================================

# 配置防火墙
configure_firewall() {
    log_step "配置防火墙..."

    local configured=false

    # UFW (Ubuntu/Debian)
    if command -v ufw &>/dev/null; then
        ufw allow 5230/tcp 2>/dev/null || true
        log_success "UFW 防火墙规则已添加"
        configured=true
    fi

    # firewalld (CentOS/RHEL/Aliyun Linux)
    if command -v firewall-cmd &>/dev/null; then
        if systemctl is-active firewalld &>/dev/null; then
            firewall-cmd --permanent --add-port=5230/tcp 2>/dev/null || true
            firewall-cmd --reload 2>/dev/null || true
            log_success "firewalld 防火墙规则已添加"
            configured=true
        fi
    fi

    # iptables (通用)
    if [ "$configured" = false ]; then
        if command -v iptables &>/dev/null; then
            if ! iptables -C INPUT -p tcp --dport 5230 -j ACCEPT &>/dev/null 2>&1; then
                iptables -I INPUT -p tcp --dport 5230 -j ACCEPT
                if command -v iptables-save &>/dev/null; then
                    iptables-save > /etc/iptables.rules 2>/dev/null || true
                fi
                log_success "iptables 防火墙规则已添加"
            fi
        else
            log_warn "未检测到防火墙，请手动开放 5230 端口"
        fi
    fi
}

# 配置定时备份
setup_cron_backup() {
    log_step "配置定时备份..."

    local cron_file="/etc/cron.d/divinesense-backup"

    if [ "$DEPLOY_MODE" = "docker" ]; then
        cat > "$cron_file" << EOF
# DivineSense 每日自动备份
0 2 * * * root cd ${INSTALL_DIR} && ./deploy.sh backup && ./deploy.sh cleanup > /dev/null 2>&1
EOF
    else
        cat > "$cron_file" << EOF
# DivineSense 每日自动备份
0 2 * * * root ${INSTALL_DIR}/deploy-binary.sh backup > /dev/null 2>&1
EOF
    fi

    chmod 644 "$cron_file"

    # 确保 cron 服务运行
    systemctl enable crond 2>/dev/null || systemctl enable cron 2>/dev/null || true

    log_success "定时备份已配置 (每天凌晨 2 点)"
}

# 复制管理脚本
copy_management_scripts() {
    log_step "安装管理脚本..."

    # 二进制模式需要复制管理脚本
    if [ "$DEPLOY_MODE" = "binary" ]; then
        # 从仓库复制脚本
        if [ -f "${INSTALL_DIR}/deploy/aliyun/deploy-binary.sh" ]; then
            cp "${INSTALL_DIR}/deploy/aliyun/deploy-binary.sh" "${INSTALL_DIR}/deploy-binary.sh"
            chmod +x "${INSTALL_DIR}/deploy-binary.sh"
            log_success "管理脚本已安装: ${INSTALL_DIR}/deploy-binary.sh"
        fi
    fi
}

# 显示部署结果
show_result() {
    local server_ip=$(get_server_ip)

    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║${NC}  ${GREEN}部署完成！${NC}                                                  ${GREEN}║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${CYAN}部署模式:${NC} ${YELLOW}${DEPLOY_MODE^^}${NC}"
    echo ""
    echo -e "${CYAN}访问信息:${NC}"
    echo -e "  URL:      ${YELLOW}http://${server_ip}:5230${NC}"
    echo ""
    echo -e "${CYAN}重要文件:${NC}"
    if [ "$DEPLOY_MODE" = "docker" ]; then
        echo -e "  配置文件: ${INSTALL_DIR}/.env.prod"
        echo -e "  数据库密码: ${INSTALL_DIR}/.db_password"
        echo -e "  备份目录: ${BACKUP_DIR}"
    else
        echo -e "  配置文件: ${CONFIG_DIR}/config"
        echo -e "  数据库密码: ${CONFIG_DIR}/.db_password"
        echo -e "  管理脚本: ${INSTALL_DIR}/deploy-binary.sh"
    fi
    echo ""
    echo -e "${CYAN}常用命令:${NC}"
    if [ "$DEPLOY_MODE" = "docker" ]; then
        echo -e "  查看状态: ${YELLOW}cd ${INSTALL_DIR} && ./deploy.sh status${NC}"
        echo -e "  查看日志: ${YELLOW}cd ${INSTALL_DIR} && ./deploy.sh logs${NC}"
        echo -e "  重启服务: ${YELLOW}cd ${INSTALL_DIR} && ./deploy.sh restart${NC}"
    else
        echo -e "  查看状态: ${YELLOW}systemctl status divinesense${NC}"
        echo -e "  查看日志: ${YELLOW}journalctl -u divinesense -f${NC}"
        echo -e "  重启服务: ${YELLOW}systemctl restart divinesense${NC}"
        echo -e "  管理工具: ${YELLOW}${INSTALL_DIR}/deploy-binary.sh${NC}"
    fi
    echo ""
    echo -e "${YELLOW}⚠️  下一步:${NC}"
    if [ "$DEPLOY_MODE" = "docker" ]; then
        echo -e "  1. 配置 AI API Keys: ${YELLOW}vi ${INSTALL_DIR}/.env.prod${NC}"
        echo -e "  2. 重启服务: ${YELLOW}cd ${INSTALL_DIR} && ./deploy.sh restart${NC}"
    else
        echo -e "  1. 启用服务: ${YELLOW}sudo systemctl enable --now divinesense${NC}"
        echo -e "  2. 配置 API Keys: ${YELLOW}sudo vi ${CONFIG_DIR}/config${NC}"
        echo -e "  3. 重启服务: ${YELLOW}sudo systemctl restart divinesense${NC}"
    fi
    echo ""
    if [ "$DEPLOY_MODE" = "binary" ]; then
        echo -e "${CYAN}Geek Mode 提示:${NC}"
        echo -e "  1. 安装 Claude Code: ${YELLOW}npm install -g @anthropic-ai/claude-code${NC}"
        echo -e "  2. 启用 Geek Mode: 在 ${CONFIG_DIR}/config 中设置 DIVINESENSE_CLAUDE_CODE_ENABLED=true"
        echo ""
    fi
}

# =============================================================================
# 主函数
# =============================================================================

main_docker_mode() {
    install_docker
    install_docker_compose
    setup_docker_mirror
    create_install_dir
    clone_repo
    generate_docker_env_file
    pull_images
    deploy_docker_services
    wait_for_docker_service
}

main_binary_mode() {
    setup_binary_user_dirs
    download_binary
    generate_binary_config
    install_binary_service
    setup_binary_postgres
    set_binary_permissions
    copy_management_scripts
}

wait_for_docker_service() {
    log_step "等待服务启动..."

    local max_wait=90
    local waited=0

    while [ $waited -lt $max_wait ]; do
        # 检查 PostgreSQL 是否就绪
        if docker exec divinesense-postgres pg_isready -U divinesense &>/dev/null; then
            log_success "PostgreSQL 已就绪"
        fi

        # 检查 DivineSense 是否就绪
        if docker exec divinesense sh -c "cat < /dev/null > /dev/tcp/127.0.0.1/5230" 2>/dev/null; then
            log_success "DivineSense 已就绪"
            return 0
        fi

        sleep 3
        waited=$((waited + 3))
        echo -n "."
    done

    echo ""
    log_warn "服务可能需要更长时间启动，请检查日志"
    return 0
}

# 主入口
main() {
    # 解析参数
    parse_args "$@"

    print_banner

    # 系统检查
    check_root
    detect_os
    detect_arch
    check_system_resources
    install_base_tools

    # 创建安装目录
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$BACKUP_DIR"

    # 根据模式执行安装
    if [ "$DEPLOY_MODE" = "docker" ]; then
        main_docker_mode
    else
        main_binary_mode
    fi

    # 通用配置
    configure_firewall
    setup_cron_backup

    show_result
}

# 运行主函数
main "$@"
