#
# Use this configuration with WireGuard client
#
[Interface]
PrivateKey = {{.privatekey}}
Address = 10.14.0.2/16
DNS = 127.0.0.1
MTU = 1420
# PreUp = "C:\Program Files\WireGuard\bat\routes-up.bat"
# PostUp = "C:\Program Files\WireGuard\bat\dns-up.bat"
# PreDown = "C:\Program Files\WireGuard\bat\routes-down.bat"
# PostDown = "C:\Program Files\WireGuard\bat\dns-down.bat"

[Peer]
PublicKey = {{.publickey}}
AllowedIPs = {{.allowedips}}
Endpoint = {{.addr}}:51820

