# inventory-manager

This is the code for the `inventory-manager` component, which is in charge of managing everything that has something to do with assets, edge-controllers, and devices.

## Getting Started

The component resides in the management cluster and handles the information going and coming from all the elements in the inventory (`assets`, `edge-controllers` and `devices`).

### Prerequisites

* [`authx`](https://github.com/nalej/authx)
* [`vpn-server`](https://github.com/nalej/vpn-server)
* [`system-model`](https://github.com/nalej/system-model)
* [`device-manager`](https://github.com/nalej/device-manager)
* [`nalej-bus`](https://github.com/nalej/nalej-bus)
* [`network-manager`](https://github.com/nalej/network-manager)
* [`edge-inventory-proxy`](https://github.com/nalej/edge-inventory-proxy)

### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, downloads the required dependencies, runs existing tests and generates ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

No real tests are being performed in this repository at the moment.

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```


## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the available versions, see the [tags on this repository](https://github.com/nalej/inventory-manager/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/inventory-manager/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
