# MongoDB with Docker

## Run

```bash
docker run -it --name mongodb \
-e MONGO_INITDB_ROOT_USERNAME=root \
-e MONGO_INITDB_ROOT_PASSWORD=root \
-v mongo_volume:/data/db -d \
-p 27017:27017 mongo
```

## Login

```bash
docker exec -it mongodb /bin/bash

mongo -u root -p root --authenticationDatabase admin
```

## Create Database

```bash
use mydb
```

## Create User

```bash
db.createUser({user: 'myuser', pwd: 'mypass', roles:[{'role': 'readWrite', 'db': 'mydb'}]});
```

## Login with other User

```bash
mongo -u myuser -p mypass --authenticationDatabase mydb

use mydb

show collections
```