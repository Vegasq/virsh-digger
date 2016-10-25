# Purpose

In my experience sometime you need to check every node on KVM host due to determinate some problem. Like some of nodes stuck during booting. With this tool it's possible to just go VNC to nodes one-by-one.

# Usage
```
./qemu_vnc_discovery -host 192.168.0.100 -port 22 -user root -password PaSsWoRd
```

# Requirements

```
brew install gtk-vnc
```
