# BEP_1920_Q2

The front-end will be developed in Angular with TypeScript, all source code can be found in `front-end/`.

The back-end will be developed in Go, all source code can be found in `back-end/`.

The computer client library will be developed in Python, all source code can be found in `cc-library/`.

Each package has its own readme for further information specific to that package.

## License
The license can be found in `LICENSE.md` in the root of this project.

### Structure

## S.C.I.L.E.R Development back-end
So you want to improve our system but don't what to read the report? 
here is a summery of were you need to be to add the feature or remove the bug you need

### structure
```
SCILER
│   README.md
└─── back-end  
│   │   README.md
│   │   
│   └───resources
│   │   └───manuals
│   │   └───production
│   │   └───testing
│   │   
│   └───src
│       └───sciler
│           │   README.md
│           │   main.go    
│           │
│           └───communication
│           │   │   communication.go
│           │   │   communication_test.go
│           │   
│           └───handler
│           │   │   handler.go
│           │   │   handler_helper.go
│           │   │   handler_test.go
│           │   │   confirmation_test.go
│           │   │   instruction_test.go
│           │   │   status_test.go
│           │   
│           └───config
│               │   configHandler.go
│               │   configHandler.go
│               │   readConfigTypes.go
│               │   workingConfigTypes.go
│               │   workingConfigTypes_test.go
│
└─── front-end
│   │   README.md
│   │
│   └─── src
│
└─── cc_library
│   │   README.md
│   │
│   └───py_scc
│   │   │   README.md
│   │   └───scclib
│   │   
│   └───js_scc
│   │   │   README.md
│   │   │   index.js
│   │   
│   └───example_scripts_py_scc
│   └───example_scripts_js_scc

```

