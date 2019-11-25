## Back-end
The back-end is written in Go. Tools used are go fmt, golint and Testify. 

External imports that need to be retrieved should be added in `.travis.yml` in the `before_script`.
 
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
