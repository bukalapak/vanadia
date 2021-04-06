# Vanadia

[![Build Status](https://travis-ci.org/bukalapak/vanadia.svg?branch=master)](https://travis-ci.org/bukalapak/vanadia)

A utility to export [API Blueprint](https://github.com/apiaryio/api-blueprint) `.apib` file to a Postman collection.

## Running locally

If you want to test your APIB files locally without installing vanadia on your machine, follow these steps:
```sh
git clone https://github.com/SharperShape/vanadia.git
cd vanadia
docker build .
docker tag $(docker images -q | head -1) vanadia:latest
```

Add this to your `.profile` or `.bashrc`:
```sh
alias vanatest='docker run -v "$(pwd)":/data vanadia -input documentation.apib -output TEST.postman_collection.json'
```
If you are using fish and podman, run this in the shell
```sh
alias -s vanatest='podman run --privileged -v (pwd):/data vanadia -input documentation.apib -output TEST.postman_collection.json'
```


Now you can run `vanatest` in a repository where you want to test your API documentation.
Load the `TEST.postman_collection.json` file into your Postman desktop app and verify that the endpoints and the documentation is correct.

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
