# WireGuard Helper

[中文](README.md) | [English](README_EN.md)

## WireGuard Helper 是一个用于管理和连接 WireGuard 隧道的工具。它提供了加载配置、管理插件和处理连接的功能。

### 特性
* 加载和管理 WireGuard 隧道配置
* 支持插件以扩展功能
* 基于模板的配置生成
* 网络连接自动检测与切换
* 可自定义等待时间的连接和断开处理

### 配置
配置文件请参考 `bindemo` 目录下的 `config` 。

### 使用
#### 1.前置条件
需安装 [WireGuard](https://www.wireguard.com/install/) 。

#### 2.下载可执行文件
从 [Releases](https://github.com/ace-zhaoy/wireguard-helper/releases) 页面下载适用于您的操作系统的可执行文件。

#### 3.配置
* 将`bindemo`目录下的配置`config` 复制到可执行文件所在目录。
* 修改 `tunnel_manager.yaml`文件中`tunnels`字段的配置，配置文件中的`tunnels`字段是一个列表，每个元素是一个隧道配置，包括隧道名称、配置文件路径、插件路径、插件参数等。
* 在 `config/tpl`下放置隧道配置文件模板

#### 4.运行
需使用管理员权限在命令行运行。

### 许可证
此项目使用 MIT 许可证。有关详细信息，请参阅 [LICENSE](LICENSE) 文件。