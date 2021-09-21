

Build for ev3dev and run on EV3:

```
    GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -o toypad-read.arm .
    scp toypad-read.arm robot@ev3dev.local
    ssh -t robot@ev3dev.local sudo ./toypad-read.arm
```