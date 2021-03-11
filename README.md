# Vanadia

[![Build Status](https://travis-ci.org/bukalapak/vanadia.svg?branch=master)](https://travis-ci.org/bukalapak/vanadia)

A utility to export [API Blueprint](https://github.com/apiaryio/api-blueprint) `.apib` file to a Postman collection.

## Installation

The latest executables for Linux and OSX are available from the [release page](https://github.com/SharperShape/vanadia/releases), so it can be executed directly:

```sh
$ wget https://github.com/SharperShape/vanadia/releases/download/${VERSION}/vanadia-${VERSION}.${OS}-amd64.tar.gz
$ tar -xzf vanadia-${VERSION}.${OS}-amd64.tar.gz
$ ./vanadia -h
```

### Manual build

Manual build needs python2 so if your python is python3:

```sh
 $ virtualvenv venv -p python2
 $ source venv/bin/activate
```

If you want to be in bleeding edge, you can manually build from `master`:

```sh
$ cd ~/go/src/github.com/SharperShape/
$ git clone https://github.com/SharperShape/vanadia.git
$ cd vanadia
$ make
```

Make sure you have Go 1.9 and build-essentials(fedora; dnf install @development-tools) as we should compile [Drafter](https://github.com/apiaryio/drafter) as one of its dependency.

## Usage

Let's say we have an API Blueprint document, `API.apib` in our working directory. Then we can do:

```sh
$ ./vanadia --input API.apib --output API.postman_collection.json
```

Vanadia can also read input from standard input and give its output via standard output; just omit the `--input` and `--output` flag.

### Configuration

Vanadia is configurable with a `vanadia.yml` in working directory. Please refer to `vanadia.yml` in this repository to see what is configurable.

You can also configure Vanadia from other location by specifying the config file:

```sh
$ ./vanadia --input api/API.apib --output api/API.postman_collection.json --config api/vanadia.yml
```
