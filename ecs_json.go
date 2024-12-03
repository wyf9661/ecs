package main

var jsonTemplate = `{
    "ociVersion": "1.0.0",
    "platform": {
        "os": "sylixos",
        "arch": "arm64"
    },
    "process": {
        "user": {
            "uid": 0,
            "gid": 0 
        }, 
        "args": [
            ""
        ], 
        "env": [
            "PATH=/usr/bin:/bin:/usr/pkg/sbin:/sbin:/usr/local/bin",
            "LD_LIBRARY_PATH=/usr/lib:/lib:/usr/local/lib"
        ], 
        "cwd": "/"
    },
    "root": {
        "path": "rootfs",
        "readonly": false
    },
    "hostname": "sylixos_ecs",
    "mounts": [
        {
            "destination": "/etc/lic",
            "source": "/etc/lic",
            "options":["ro"]
        }
    ], 
    "sylixos": {
        "devices": [
            {
                "path": "/dev/fb0",
                "access": "rw"
            },
            {
                "path": "/dev/input/xmse",
                "access": "rw"
            },
            {
                "path": "/dev/input/xkbd",
                "access": "rw"
            },
            {
                "path": "/dev/net/vnd",
                "access": "rw"
            }
        ], 
        "resources": {
            "cpu": {
                "highestPrio": 150,
                "lowestPrio": 250,
                "defaultPrio": 200
            },
            "itimer": {
                "defaultPrio": 200
            },
            "affinity": [
            ],
            "memory": {
                "kheapLimit": 536870912,
                "memoryLimitMB": 512
            },
            "kernelObject": {
                "threadLimit": 4096,
                "threadPoolLimit": 32,
                "eventLimit": 32768,
                "eventSetLimit": 500,
                "partitionLimit": 6000,
                "regionLimit": 50,
                "msgQueueLimit": 8192,
                "timerLimit": 64,
                "rmsLimit": 32,
                "threadVarLimit": 16,
                "posixMqueueLimit": 300,
                "dlopenLibraryLimit": 50,
                "xsiipcLimit": 100,
                "socketLimit": 1024,
                "srtpLimit": 30,
                "deviceLimit": 60
            },
            "disk": {
                "limitMB": 2048
            }
        },
        "commands": [
            "exec",
            "top",
            "cpuus",
            "vi",
            "cat",
            "touch",
            "ps",
            "ts",
            "tp",
            "ss",
            "ints",
            "ls",
            "cd",
            "pwd",
            "modules",
            "varload",
            "varsave",
            "shstack",
            "srtp",
            "shfile",
            "help",
            "debug",
            "shell",
            "ll",
            "sync",
            "ln",
            "kill",
            "free",
            "ifconfig",
            "mems",
            "env",
            "rm",
            "exit"
        ], 
        "network": {
            "telnetdEnable": true,
            "ftpdEnable": true,
            "sshdEnable": false
        }
    }
}`
