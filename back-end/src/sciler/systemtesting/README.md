## System Testing
### deployment for windows:
run as administrator (in powershell)
```
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Windows-Subsystem-Linux
```
install linux subsystem through microsoft store

### development for windows:
When edition the bash scripts in windows, make sure to fix newline characters by
running the bash file through dos2unix:
- run `sudo apt-get install dos2unix`
- run `dos2unix [filename.sh]`
### broker
- edit mosquitto broker `bind address` from `localhost` to the local ip of the Pi
- run script by `sudo ./setupAndRunBroker.sh`
