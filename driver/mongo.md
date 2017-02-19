# Mongo storage driver implementation
To use [MongoDB](https://www.mongodb.com/) as the backend for storing profiling data, you must first install [Docker](https://www.docker.com/).
Once docker is installed, you need to pull the official docker container for mongo:  
```$ docker pull mongo```  
Then to run the docker container:  
```$ docker run --volume /data/db:/data/db -p 27017:27017 mongo -d```  
Voila! You now have a MongoDB instance running on port 27017, with data being written to /data/db on your host.  
### Configuration
The minimal configuration file to use Mongo as the backend:  
```json
{
    "hostname": "127.0.0.1",
    "port": "27017",
    "driver": "mongo"
}
```