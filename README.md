# apib-to-postman

A utility to export [API Blueprint](https://github.com/apiaryio/api-blueprint) `.apib` file to a Postman collection.

## Installation

```sh
git clone https://github.com/bukalapak/apib-to-postman.git
cd apib-to-postman
make
```

## Usage

Let's say we have an API Blueprint document, `API.apib` in our working directory. Then we can do:

```sh
./apib-to-postman --input API.apib --output API.postman_collection.json
```

## Configuration

The utility is configurable with a `config.yml` in working directory. Please refer to `config.yml` in this repository to see what is configurable.
