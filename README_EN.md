# WireGuard Helper

[中文](README.md) | [English](README_EN.md)

## WireGuard Helper is a tool for managing and connecting WireGuard tunnels. It provides features for loading configurations, managing plugins, and handling connections.

### Features
* Load and manage WireGuard tunnel configurations
* Support plugins to extend functionality
* Template-based configuration generation
* Automatic network connection detection and switching
* Customizable wait times for connection and disconnection handling

### Configuration
Please refer to the `config` in the `bindemo` directory for configuration files.

### Usage
#### 1. Prerequisites
Install [WireGuard](https://www.wireguard.com/install/).

#### 2. Download the executable
Download the executable for your operating system from the [Releases](https://github.com/ace-zhaoy/wireguard-helper/releases) page.

#### 3. Configuration
* Copy the `config` from the `bindemo` directory to the directory where the executable is located.
* Modify the `tunnel_manager.yaml` file in the `config` directory. The `tunnels` field in the configuration file is a list, and each element is a tunnel configuration, including tunnel name, configuration file path, plugin path, plugin parameters, etc.
* Place the tunnel configuration file templates in the `config/tpl` directory.

#### 4. Run
Run with administrator privileges in the command line.

### License
This project is licensed under the MIT License. For more details, please refer to the [LICENSE](LICENSE) file.