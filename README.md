# Server
[![Build Status](https://travis-ci.org/auto-profile/server.svg?branch=master)](https://travis-ci.org/auto-profile/server)[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
                                                                                                                         
  
This project serves as an API server that proxies profiling data requests from an agent and inserts them (via a [Driver](https://github.com/auto-profile/server/tree/master/driver)) to a given database backend.

## Installation
To get started, be sure you have [Go 1.7 or higher]() installed and a local version of [ElasticSearch](https://www.elastic.co/guide/en/elasticsearch/guide/current/running-elasticsearch.html) or [MongoDB](https://docs.mongodb.com/manual/installation/) running.    
#### Step 1: Get source

We first need to go get the source for the project:

```
go get github.com/auto-profile/server
```

#### Step 2: Build  

We can now run a build:

```
cd $GOPATH/src/github.com/auto-profile/server && go build
```

#### Step 3: Create configuration file

In order for the server to connect to one of the backend databases with the drivers we have available, we must specify a configuration file. A sample MongoDB config.json file might look something like:  

```json
{
    "driver": "mongo",
    "hostname": "127.0.0.1",
    "port": "27017"
}
```

#### Step 4: Run the server

We can now run our local server and point it to the configuration file:
```
./server -configPath ./config.json
```
And with that, your server is now able to receive requests from an agent at ```http://127.0.0.1:8080```!  

## Drivers
Because of the variable nature of profiling data, NoSQL backends seem to fit nicely. We currently have supported [ElasticSearch](https://github.com/auto-profile/server/blob/master/driver/elasticsearch.md) and [MongoDB](https://github.com/auto-profile/server/blob/master/driver/mongo.md) drivers.
