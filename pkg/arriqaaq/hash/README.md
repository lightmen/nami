# hash

hash is an in-memory Redis like [hash](https://redis.io/commands#hash) datastructure

| Property | Type | Description |
|---|---|---|
| `key` | `string` | The identifier of the hash key |
| `member` | `string` | The identifier of the member. It must be unique within the hash. |
| `value` | `interface {}` | value associated with this node |


Getting Started
===============

## Installing

To start using hash, install Go and run `go get`:

```sh
$ go get -u github.com/arriqaaq/hash
```

This will retrieve the library.

## Usage

```go
package main

import "github.com/arriqaaq/hash"

func main() {
	hash := hash.New()

	hash.HSet(key, "a", []byte("hash_1"))
	hash.HSet(key, "b", []byte("hash_2"))
	hash.HSet(key, "c", []byte("hash_3"))
}
```

## Supported Commands

```go
HDEL
HEXISTS
HGET
HGETALL
HKEYS
HLEN
HSCAN
HSET
HSETNX
HVALS

```
