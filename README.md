# Purpose

In my experience sometime you need to check every node on KVM host due to determinate some problem. Like some of nodes stuck during booting. With this tool it's possible to just go VNC to nodes one-by-one.

# Usage
Collect all nodes with "admin" in name. And connect to it thru vnc with gvncviewer.
```
for i in $(./virsh-digger -host 192.168.0.100 -port 22 -user root -password r00t -node admin --vncaddr);
  do gvncviewer $i;
done
```

# Requirements

```
brew install gtk-vnc
```
