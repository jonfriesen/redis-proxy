# HTTP Redis-Proxy

## Overview
This is a simple HTTP Redis proxy to remove pressure from the central Redis cache. The proxy handles time based key expiry and least recently used (LRU) eviction to only keep the most current key-value pairs available.

## Assumptions
- API endpoint/path formatting was not specified, assuming any custom format is acceptable
- There doesn't appear to be a HTTP proxy requirement for write operations, nor is it requested in the spec so a write endpoint is not included. This is justified as this project is designed to reduce pressure from lookups on the Redis instance

## Requirements
- Proxy to a single "backing" Redis server caching values on the proxy
- Proxy needs to protect Redis from excessive calls
- Cached records need to expire after a time duration
- Cache needs to be size restricted by record count
- Least recently used records need to be evicted when the cache size is exceed when adding new records
- Cache requests need to be accepted sequentially
- Single-click build and testing
- End-to-end tests for all requirements

## Architecture
This package is built to be extendable and easy to understand and consume. It is composed of three main groupings of code with a structure that looks like:

```
redis-proxy
    |
    |-- api
    |-- cache
    |   `-- lrucache
    |-- cmd
    |   `-- redis-proxy-http
    |-- scripts
    |-- storage
    |   |-- inmem
    |   `-- redis
    `-- tests
```

In order:

- api contains the HTTP endpoints, middleware, and composition of dependencies which are injected by main.
- cache contains a logic layer that handles getting and automatically populating from the storage interface that is created on app startup.
      - lrucache is cache implementation, it has two methods, Get and Push which are self-explanatory. To make the lrucache safe in concurrent situations it requires manually locking and unlocking.
- cmd holds different main classes, currently it's just `redis-proxy-http` but in the future could hold other projects, for example a gRPC implementation.
- storage contains our datasource interface and implementations.
    - inmem is an basic in memory implementation used during development.
    - redis is our production implementation of our data storage.
- tests holds end to end integration tests designed to be run using the `make test` command.

### API
The API is a REST HTTP with a single endpoint:

```
/v1/get/<lookup key>
```

cURL sample:
```
curl -X GET http://localhost:4000/v1/get/cachedKey
```

### Cache
The cache holds the records in a doubly linked list alongside a hashmap for key to record lookups.

Expanded, the linked list is used to manage order of recently used records. Built as a queue, new and recently used records are pushed on the end. Older records make their way to the head where they will be the first ones to face eviction when the cache is full. Full is defined as the amount of keys, a configuration parameter, set at run time.

The lrucache has two actions

#### Get
Retrieving a record from the cache occurs in O(1) time. As we are using a map to associate our key to records which have a constant lookup time we do not have any iterations or searching for values.

Get also handles expiry of nodes. After a set amount of time, records become expired. This check is done during the get operation. In the event that a record is expired, a not found error is returned so current data can be pulled from the Redis data source.

#### Push
Adding new records to the cache uses O(1) time, similarly to the Get action. When adding new records we append our record node to the end of the queue which a pointer is maintained and add it to the map; both operations occur in constant time. In the event that the cache is full we remove the head of the queue and map the now nil head pointer to the next least used node.

### Testing
As per the requirements there are end-to-end (e2e) integration tests that run against the full server + Redis instance. The e2e test suite uses Go httptest libraries to spin up a server and connect to said servers. This means that the code is running in an official capacity with the ability to debug tests through the application and test code.

It should be noted that the main class is not tested with this strategy. This is justified as the main class is merely wiring up the dependencies and watching for an exit, neither of which are conducive to the testing.

## Usage
Requirements:
- Docker
- Docker-compose
- make
- bash
- Internet connection

### Running the server
Running `make run` will start a containerized service available at localhost:4000 with it's own Redis container.

### Executing tests
Running `make test` will run the entire suite of tests including the e2e integration tests.

### Building a reusable Container
Running `make build` will create a container image called `redis-proxy` that can be used to interface with remote or custom Redis instances.

```
Usage of redis-proxy-http:
  -cachesize int
        Represents the amount of keys that can be help in the cache (default 10)
  -expiry int
        Represents the expiration limit in milliseconds (default 60000)
  -host string
        Host IP address (default "0.0.0.0")
  -port string
        Host port (default "4000")
  -redis-host string
        Redis IP address (default "0.0.0.0")
  -redis-port string
        Redis port (default "6379")
```

## Time Breakdown

| Task | ~Time (Minutes) |
| --- |:---:| 
| Design | 45 |
| Cache | 65 |
| HTTP/API | 35 |
| Storage | 60 |
| Tests | 75 |
| Docs | 20 |
| Make | 15 |
| **Total** | **5h 15m** |


## Missing Requirements
The core requirements are fully met with this implementation. Time constraints prevented the bonus requirements from being implemented.

## Roadmap
### Phase 2
- Add more e2e tests (better concurrency in particular)
- Add unit tests
- Revise docker files and build system to for simplicity and clarity
- Enhance input validation

### Phase 3
- Add logging middleware
- Add metrics endpoint for reporting
- Add Parallel concurrent processing to stop specific locking of the cache
- Add Redis interface for access via HTTP

