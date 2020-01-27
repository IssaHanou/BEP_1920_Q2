# Sciler back-end
The back-end server contains and manages all logic in the escape room.

## Deployment on Raspberry Pi
#### Make config:
- in directory `back-end` (copy and) rename [resources/`room_config.json.default`](resources/production/room_config.json.default) to `room_config.json`
- (optional), change this file as desired following the [room_manual.md](resources/manuals/room_manual.md)

#### If needed install go1.13.4:
- run  `pi@raspberrypi:~ $ wget https://storage.googleapis.com/golang/go1.13.4.linux-armv6l.tar.gz`
- run  `pi@raspberrypi:~ $ sudo tar -C /usr/local -xzf go1.13.4.linux-armv6l.tar.gz`
- run ` pi@raspberrypi:~ $ export PATH=$PATH:/usr/local/go/bin` 
- run ` pi@raspberrypi:~ $ go version` to check if succeeded 

#### get ready to run
- add to .profile (with `pi@raspberrypi:~ $ nano .profile` for example): \
`export GOROOT=/usr/local/go` \
`export PATH=$PATH:$GOROOT/bin`\
`export GOPATH=~/go/src/BEP_1920_Q2/back-end`\
`export PATH=$PATH:$GOPATH/bin`
- run `pi@raspberrypi:~/go/src $ git clone https://github.com/IssaHanou/BEP_1920_Q2.git`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2/back-end/src/sciler $ go get ./...`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2 $ go install sciler`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2 $ sciler`

##### To update to latest version
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2 $ git pull`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2 $ go install sciler`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2 $ sciler`

##### To run on boot
- use tool like supervisord
- command `$GOPATH/bin/sciler` from  directory ~/go/src/BEP_1920_Q2 with $GOPATH=/home/pi/go/src/BEP_1920_Q2/back-end

## Development
The back-end is written in Go. Tools used are go fmt, golint and Testify. 

External imports that need to be retrieved should be added in `.travis.yml` in the `before_script`.
Configure (project) GOPATH to /back-end folder in /BEP_1920_Q2.
When using Goland uncheck `Use GOPATH that is defined in system environment` and check `Index entire GOPATH`

run `go get ./...` in the `sciler` folder to go get all dependencies

