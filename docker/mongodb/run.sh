#!/bin/bash

docker run -it --name mongodb \
-e MONGO_INITDB_ROOT_USERNAME=root \
-e MONGO_INITDB_ROOT_PASSWORD=root \
-v mongo_volume:/data/db -d \
-p 27017:27017 mongo