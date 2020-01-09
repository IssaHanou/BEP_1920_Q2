### License

GNU GENERAL PUBLIC LICENSE Version 3, [see `LICENSE.md`](LICENSE.md)

### run prettier

`prettier --write "**/*.js"`

### Using this library
- npm install with ```npm install js-scc```
- import lib with ```const SccLib = require("js-scc");```
- in case of:
    - angular: add to dependency in `package.json`
    - browser javascript: (example nodejs servering webpage with javascript which includes this library), [Browserify](http://browserify.org/) your javascript which includes this library.


### tip reading in config.json in Angular:
When reading in a json file
```
import * as data from './data.json';
```

make sure `tsconfig.json` hase `resolveJsonModule: true`:
```
{
  "compilerOptions": {
    ...
    "resolveJsonModule": true,
    ...
}
```

then you can use this object in for example `ngOnInit`
```
  ngOnInit(): void {
       this.jsonData = (data  as  any).default
  }
```
