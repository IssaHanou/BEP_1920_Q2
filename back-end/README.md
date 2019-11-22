# Setup gofmt and golint to run automatically:
- install golint 
```
    go get -u golang.org/x/lint/golint
```
- natigate to Settings > Tools > File Watchers 
- press + go fmt
- enable auto saved files to trigger the watcher in advanced settings
- press ok
- select and copy and paste the go fmt watcher
- edit name and program to "golint"
- edit arguments to "-set_exit_status $FilePath$"
- press ok