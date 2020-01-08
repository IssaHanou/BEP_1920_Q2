### License

GNU GENERAL PUBLIC LICENSE Version 3, [see `LICENSE.md`](LICENSE.md)

### Setup

install mqtt: `npm install mqtt --save`
install prettier: `npm install --global prettier`

### run prettier

`prettier --write "**/*.js"`

### Setup development script using this module

run `npm link` in js-scc
run `npm link js-scc` in folder of development script (where the module will be used)

### When using in angular:

- make sure `polyfill.ts` includes the following lines:
    ```
    (window as any).global = window;
       global.Buffer = global.Buffer || require("buffer").Buffer;
       (window as any).process = {
         version: ""
       };
    ```
