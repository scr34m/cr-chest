# A simple Clash Royale chest opener

Use [cr-proxy](https://github.com/royale-proxy/cr-patcher) to create modified libg.so or you can make it by hand and install with ADB. Offsets are in the [cr-proxy](https://github.com/royale-proxy/cr-proxy) wiki to use.

Example pre / post host name values:
```
game.clashroyaleapp.com
67 61 6d 65 2e 63 6c 61 73 68 72 6f 79 61 6c 65 61 70 70 2e 63 6f 6d

locl.clashroyaleapp.com
6c 6f 63 6c 2e 63 6c 61 73 68 72 6f 79 61 6c 65 61 70 70 2e 63 6f 6d
```

Example usage for ADB (rooted device required): 
```
adb push libg.so /sdcard/libg.so
adb shell
su
cp /sdcard/libg.so /data/data/com.supercell.clashroyale/lib/
chown system:system /data/data/com.supercell.clashroyale/lib/libg.so
chmod 0755 /data/data/com.supercell.clashroyale/lib/libg.so
```

Resources used:
 - https://github.com/royale-proxy/cr-proxy
 - https://github.com/royale-proxy/cr-patcher
 - https://github.com/clugh/coc-proxy
