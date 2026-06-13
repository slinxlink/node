#!/bin/bash

set -e

CUR_DIR=$(pwd)

RED='\033[0;31m'
PINK='\033[38;2;240;128;128m'
PLAIN='\033[0m'

step() {
    echo ""
    echo -e "${PINK}>>> $1${PLAIN}"
    echo -e "${PINK}————————————————————————————————————————${PLAIN}"
}

# 检查 root
if [[ $EUID -ne 0 ]]; then
    echo -e "${RED}请以 root 用户运行此脚本${PLAIN}"
    exit 1
fi

step "更新系统"
apt update

step "安装依赖"
apt install -y curl chrony mtr gzip sqlite3

step "设置系统时间同步"
timedatectl set-timezone Asia/Shanghai
systemctl enable chrony
systemctl restart chrony

step "检测 CPU 架构"
ARCH=$(uname -m)
if [[ $ARCH == "x86_64" ]]; then
    SLINX_ARCH="amd64"
elif [[ $ARCH == "aarch64" ]]; then
    SLINX_ARCH="arm64"
else
    echo -e "${RED}不支持的 CPU 架构: $ARCH${PLAIN}"
    exit 1
fi

step "检测最新版本"
RELEASE=$(curl -s "https://api.github.com/repos/slinxlink/node/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [[ -z "$RELEASE" ]]; then
    echo -e "${RED}获取版本号失败${PLAIN}"
    exit 1
fi
echo -e "${PINK}最新版本: ${RELEASE}${PLAIN}"

step "创建目录"
SLINX_DIR="/etc/slinx"
BIN_DIR="$SLINX_DIR/bin"
DATA_DIR="$SLINX_DIR/data"
CERT_DIR="$SLINX_DIR/cert"
mkdir -p $BIN_DIR $DATA_DIR $CERT_DIR

step "下载 slinx"
systemctl stop slinx.service 2>/dev/null || true
SLINX_URL="https://github.com/slinxlink/node/releases/download/${RELEASE}/slinx_linux_${SLINX_ARCH}"
curl -fLo $SLINX_DIR/slinx $SLINX_URL
chmod +x $SLINX_DIR/slinx

step "下载 sing-box"
SINGBOX_URL="https://github.com/slinxlink/node/releases/download/${RELEASE}/sing-box_linux_${SLINX_ARCH}.gz"
curl -fLo /tmp/sing-box.gz $SINGBOX_URL
gunzip /tmp/sing-box.gz
mv /tmp/sing-box $BIN_DIR/sing-box
chmod +x $BIN_DIR/sing-box

step "开启 BBR 加速"
cat <<EOF > /etc/sysctl.d/99-bbr.conf
net.core.default_qdisc=fq
net.ipv4.tcp_congestion_control=bbr
EOF
sysctl --system
if sysctl net.ipv4.tcp_congestion_control | grep -q bbr; then
    echo -e "${PINK}BBR 启动成功${PLAIN}"
else
    echo -e "${RED}BBR 启用失败${PLAIN}"
fi

step "创建 systemd 服务"
cat <<EOF > /etc/systemd/system/slinx.service
[Unit]
Description=SLINX Service
After=network.target nss-lookup.target
Wants=network.target

[Service]
User=root
Group=root
Type=simple
LimitNOFILE=999999
WorkingDirectory=$SLINX_DIR
ExecStart=$SLINX_DIR/slinx
Restart=on-failure
RestartSec=5
StartLimitInterval=100s
StartLimitBurst=3

[Install]
WantedBy=multi-user.target
EOF

step "启动服务"
systemctl daemon-reload
systemctl enable slinx.service
systemctl start slinx.service

step "注册管理命令"
cat <<EOF > /usr/local/bin/slinx
#!/bin/bash
/etc/slinx/slinx cli
EOF
chmod +x /usr/local/bin/slinx

echo ""
echo -e "${PINK}>>> 安装完成${PLAIN}"
echo -e "${PINK}————————————————————————————————————————${PLAIN}"
echo -e "${PINK}管理脚本命令: slinx${PLAIN}"
echo ""
sleep 2
slinx