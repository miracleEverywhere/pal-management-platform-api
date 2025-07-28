#!/bin/bash

###########################################
# 用户自定义设置请修改下方变量，其他变量请不要修改 #
###########################################

# --------------- ↓可修改↓ --------------- #
# pmp暴露端口，即网页打开时所用的端口
PORT=80

# 数据库文件所在目录，例如：./database
CONFIG_DIR="./config"

# 虚拟内存大小，例如 1G 4G等
SWAPSIZE=2G
# --------------- ↑可修改↑ --------------- #

###########################################
#     下方变量请不要修改，否则可能会出现异常     #
###########################################

USER=$(whoami)
ExeFile="$HOME/pmp"

cd "$HOME" || exit

function echo_red() {
    echo -e "\033[0;31m$*\033[0m"
}

function echo_green() {
    echo -e "\033[0;32m$*\033[0m"
}

function echo_yellow() {
    echo -e "\033[0;33m$*\033[0m"
}

function echo_cyan() {
    echo -e "\033[0;36m$*\033[0m"
}

# 检查用户，只能使用root执行
if [[ "${USER}" != "root" ]]; then
    echo_red "请使用root用户执行此脚本 (Please run this script as the root user)"
    exit 1
fi

# 设置全局stderr为红色并添加固定格式
function set_tty() {
    exec 2> >(while read -r line; do echo_red "[$(date +'%F %T')] [ERROR] ${line}" >&2; done)
}

# 恢复stderr颜色
function unset_tty() {
    exec 2> /dev/tty
}

# 定义一个函数来提示用户输入
function prompt_user() {
    clear
    echo_green "帕鲁管理平台(PMP)"
    echo_green "--- https://github.com/miracleEverywhere/pal-management-platform-api ---"
    echo_yellow "————————————————————————————————————————————————————————————"
    echo_green "[0]: 下载并启动服务(Download and start the service)"
    echo_yellow "————————————————————————————————————————————————————————————"
    echo_green "[1]: 启动服务(Start the service)"
    echo_green "[2]: 关闭服务(Stop the service)"
    echo_green "[3]: 重启服务(Restart the service)"
    echo_yellow "————————————————————————————————————————————————————————————"
    echo_green "[4]: 更新管理平台(Update management platform)"
    echo_green "[5]: 强制更新平台(Force update platform)"
    echo_green "[6]: 更新启动脚本(Update startup script)"
    echo_yellow "————————————————————————————————————————————————————————————"
    echo_green "[7]: 设置虚拟内存(Setup swap)"
    echo_green "[8]: 退出脚本(Exit script)"
    echo_yellow "————————————————————————————————————————————————————————————"
    echo_yellow "请输入选择(Please enter your selection) [0-8]: "
}

# 检查jq
function check_jq() {
    echo_cyan "正在检查jq命令(Checking jq command)"
    if ! jq --version >/dev/null 2>&1; then
        OS=$(grep -P "^ID=" /etc/os-release | awk -F'=' '{print($2)}' | sed "s/['\"]//g")
        if [[ ${OS} == "ubuntu" ]]; then
            apt install -y jq
        else
            if grep -P "^ID_LIKE=" /etc/os-release | awk -F'=' '{print($2)}' | sed "s/['\"]//g" | grep rhel; then
                yum install -y jq
            fi
        fi
    fi
}

function check_curl() {
    echo_cyan "正在检查curl命令(Checking curl command)"
    if ! curl --version >/dev/null 2>&1; then
        OS=$(grep -P "^ID=" /etc/os-release | awk -F'=' '{print($2)}' | sed "s/['\"]//g")
        if [[ ${OS} == "ubuntu" ]]; then
            apt install -y curl
        else
            if grep -P "^ID_LIKE=" /etc/os-release | awk -F'=' '{print($2)}' | sed "s/['\"]//g" | grep rhel; then
                yum install -y curl
            fi
        fi
    fi
}

function check_strings() {
    echo_cyan "正在检查strings命令(Checking strings command)"
    if ! strings --version >/dev/null 2>&1; then
        OS=$(grep -P "^ID=" /etc/os-release | awk -F'=' '{print($2)}' | sed "s/['\"]//g")
        if [[ ${OS} == "ubuntu" ]]; then
            apt install -y binutils
        else
            if grep -P "^ID_LIKE=" /etc/os-release | awk -F'=' '{print($2)}' | sed "s/['\"]//g" | grep rhel; then
                yum install -y binutils
            fi
        fi
    fi

}

# Ubuntu检查GLIBC, rhel需要下载文件手动安装
function check_glibc() {
    check_strings
    echo_cyan "正在检查GLIBC版本(Checking GLIBC version)"
    OS=$(grep -P "^ID=" /etc/os-release | awk -F'=' '{print($2)}' | sed "s/['\"]//g")
    if [[ ${OS} == "ubuntu" ]]; then
        if ! strings /lib/x86_64-linux-gnu/libc.so.6 | grep GLIBC_2.34 >/dev/null 2>&1; then
            apt update
            apt install -y libc6
        fi
    else
        echo_red "非Ubuntu系统，如GLIBC小于2.34，请手动升级(For systems other than Ubuntu, if the GLIBC version is less than 2.34, please upgrade manually)"
    fi
}

# 下载函数:下载链接,尝试次数,超时时间(s)
function download() {
    local url="$1"
    local output="$2"
    local timeout="$3"

    unset_tty
    curl -L --connect-timeout "${timeout}" --progress-bar -o "${output}" "${url}"
    set_tty

    return $? # 返回 wget 的退出状态
}

# 检查进程状态
function check_pmp() {
    sleep 1
    if pgrep pmp >/dev/null; then
        echo_green "启动成功 (Startup Success)"
    else
        echo_red "启动失败 (Startup Fail)"
        exit 1
    fi
}

# 启动主程序
function start_pmp() {
    check_glibc
    if [ -e "$ExeFile" ]; then
        nohup "$ExeFile" -l ${PORT} -s ${CONFIG_DIR} >pmp.log 2>&1 &
    else
        install_pmp
        nohup "$ExeFile" -l ${PORT} -s ${CONFIG_DIR} >pmp.log 2>&1 &
    fi
}

# 关闭主程序
function stop_pmp() {
    pkill -9 pmp
    echo_green "关闭成功 (Shutdown Success)"
    sleep 1
}

# 删除主程序、请求日志、运行日志、遗漏的压缩包
function clear_pmp() {
    echo_cyan "正在执行清理 (Cleaning Files)"
    rm -f pmp
    rm -rf logs
}

# 设置虚拟内存
function set_swap() {
    SWAPFILE=/swapfile

    # 检查是否已经存在交换文件
    if [ -f $SWAPFILE ]; then
        echo_green "交换文件已存在，跳过创建步骤"
    else
        echo_cyan "创建交换文件..."
        sudo fallocate -l $SWAPSIZE $SWAPFILE
        sudo chmod 600 $SWAPFILE
        sudo mkswap $SWAPFILE
        sudo swapon $SWAPFILE
        echo_green "交换文件创建并启用成功"
    fi

    # 添加到 /etc/fstab 以便开机启动
    if ! grep -q "$SWAPFILE" /etc/fstab; then
        echo_cyan "将交换文件添加到 /etc/fstab "
        echo "$SWAPFILE none swap sw 0 0" | sudo tee -a /etc/fstab
        echo_green "交换文件已添加到开机启动"
    else
        echo_green "交换文件已在 /etc/fstab 中，跳过添加步骤"
    fi

    # 更改swap配置并持久化
    sysctl -w vm.swappiness=20
    sysctl -w vm.min_free_kbytes=100000
    echo -e 'vm.swappiness = 20\nvm.min_free_kbytes = 100000\n' > /etc/sysctl.d/pmp_swap.conf

    echo_green "系统swap设置成功 (System swap setting completed)"
}

# 使用无限循环让用户输入命令
while true; do
    # 提示用户输入
    prompt_user
    # 读取用户输入
    read -r command
    # 使用 case 语句判断输入的命令
    case $command in
    0)
        set_tty
        clear_pmp
        install_pmp
        start_pmp
        check_pmp
        unset_tty
        break
        ;;
    1)
        set_tty
        start_pmp
        check_pmp
        unset_tty
        break
        ;;
    2)
        set_tty
        stop_pmp
        unset_tty
        break
        ;;
    3)
        set_tty
        stop_pmp
        start_pmp
        check_pmp
        echo_green "重启成功 (Restart Success)"
        unset_tty
        break
        ;;
    4)
        set_tty
        unset_tty
        break
        ;;
    5)
        set_tty
        unset_tty
        break
        ;;
    6)
        set_tty
        unset_tty
        break
        ;;
    7)
        set_tty
        set_swap
        unset_tty
        break
        ;;
    8)
        exit 0
        break
        ;;
    *)
        echo_red "请输入正确的数字 [0-8](Please enter the correct number [0-8])"
        continue
        ;;
    esac
done