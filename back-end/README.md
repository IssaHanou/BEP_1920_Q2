## Back-end (development)
The back-end is written in Go. Tools used are go fmt, golint and Testify. 

External imports that need to be retrieved should be added in `.travis.yml` in the `before_script`.
Configure (project) GOPATH to /back-end folder in /BEP_1920_Q2.
When using Goland uncheck `Use GOPATH that is defined in system environment` and check `Index entire GOPATH`

run `go get ./...` in the `sciler` folder to go get all dependencies

## Back-end on Raspberry Pi
setup config:
- in directory `back-end` (copy and) rename [resources/room_config.json.default](./resources/room_config.json.default) to `room_config.json`
- (optional), change this file as desired following the [room_manual.md](./resources/room_manual.md)

if needed install go1.13.4:
- run  `pi@raspberrypi:~ $ wget https://storage.googleapis.com/golang/go1.13.4.linux-armv6l.tar.gz`
- run  `pi@raspberrypi:~ $ sudo tar -C /usr/local -xzf go1.7.3.linux-armv6l.tar.gz`
- run ` pi@raspberrypi:~ $ export PATH=$PATH:/usr/local/go/bin` 
- run ` pi@raspberrypi:~ $ go version` to check if succeeded 

then:
- add to .profile (with `pi@raspberrypi:~ $ nano .profile` for example): \
`export GOROOT=/usr/local/go` \
`export PATH=$PATH:$GOROOT/bin`\
`export GOPATH=~/go/src/BEP_1920_Q2/back-end`\
`export PATH=$PATH:$GOPATH/bin`
- run `pi@raspberrypi:~/go/src $ git clone https://github.com/IssaHanou/BEP_1920_Q2.git`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2/back-end/src/sciler $ go get ./...`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2 $ go install sciler`
- run `pi@raspberrypi:~/go/src/BEP_1920_Q2 $ sciler`

### Setup gofmt and golint to run automatically:
- install golint: `go get -u golang.org/x/lint/golint`
- natigate to Settings > Tools > File Watchers 
- press `+` and select `go fmt`
- enable auto saved files to trigger the watcher in advanced settings
- press `ok`
- select and copy & paste the go fmt watcher
- edit name and program to `golint`
- edit arguments to `-set_exit_status $FilePath$`
- press `ok`
