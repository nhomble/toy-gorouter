toy gorouter
============
![Go Build & Test](https://github.com/nhomble/toy-gorouter/workflows/Go%20Build%20&%20Test/badge.svg)

Does simple routing across statically configured backends. Async checks healthiness of backends and round robins to
different healthy resources.

# usage
```bash
$ cd gorouter; go build
$ ./gorouter -config ..\resources\your-configuration.yml
```

And if you want a simple backend to try it out
```bash
$ cd logserver; go build
$ ./logserver -port <port that matches config>
```

# todo
- http healthchecks
- configure base uris
- dynamically refresh configurations
- auth

# credit
inspired by kasvith