# MongoDB

The ACME Serverless Fitness Shop leverages [MongoDB](https://www.mongodb.com/) to store data.

## Create the MongoDB instance

You can run your own container with [MongoDB](https://hub.docker.com/_/mongo) installed, using the [Docker](https://hub.docker.com/_/mongo) image.

```bash
docker pull mongo
docker run --rm -it -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=mongoadmin -e MONGO_INITDB_ROOT_PASSWORD=mongoadmin -e MONGO_INITDB_DATABASE=acmefitness mongo 
```

## Seed the table

To seed MongoDB with random data, you can use the Go app in the [seed](./seed) directory. The app has three required flags and one optional one:

* `username`: The username to connect to MongoDB
* `password`: The password to connect to MongoDB
* `hostname`: The hostname of the MongoDB server
* `port`: The port number of the MongoDB server (optional)

As an example, using the default settings, you can run

```bash
cd seed
go run main.go -username=mongoadmin -password=mongoadmin -hostname=localhost -port=27017
```

To generate your own data, you can use [Mockaroo](https://www.mockaroo.com/) and import the `schema.json` files to start off.
