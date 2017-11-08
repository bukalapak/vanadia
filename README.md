# Vanadia

A utility to export [API Blueprint](https://github.com/apiaryio/api-blueprint) `.apib` file to a Postman collection.

## Installation

```sh
git clone https://github.com/bukalapak/vanadia.git
cd vanadia
make
```

## Usage

Let's say we have an API Blueprint document, `API.apib` in our working directory. Then we can do:

```sh
./vanadia --input API.apib --output API.postman_collection.json
```

## Configuration

Vanadia is configurable with a `config.yml` in working directory. Please refer to `config.yml` in this repository to see what is configurable.
