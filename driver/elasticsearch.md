# Elasticsearch storage driver implementation
To use [Elasticsearch](https://www.elastic.co/products/elasticsearch) as the backend for storing profiling data, you must first install [Docker](https://www.docker.com/).
Once docker is installed, you need to pull the official docker container:  

```$ docker run -p 9200:9200 -d -v "$PWD/db":/usr/share/elasticsearch/data elasticsearch```  

Voila! You now have a Elasticsearch instance running on port `9200`, with data being written to `$PWD/db` on your host.  

### Configuration
The minimal configuration file to use Elasticsearch as the backend:  
```json
{
    "hostname": "127.0.0.1",
    "port": "9200",
    "driver": "elasticsearch"
}
```
