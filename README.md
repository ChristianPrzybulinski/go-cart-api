# go-cart-api

[![Go Pipeline](https://github.com/ChristianPrzybulinski/go-cart-api/actions/workflows/go.yml/badge.svg)](https://github.com/ChristianPrzybulinski/go-cart-api/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/ChristianPrzybulinski/go-cart-api/branch/develop/graph/badge.svg?token=Qkz6YqhTPl)](https://codecov.io/gh/ChristianPrzybulinski/go-cart-api)

GOLang Project that provides a simple Endpoint API (**/cart**) that works as a checkout service. 

- It receives a JSON containing a list of products with *ID* and *Quantity*; 
- Returning a JSON that contains the checkout list with Its values;
- It uses an Internal Database that is populated using a provided JSON file;
- It also has a flag that checks If It's Black Friday day, so It will add a Gift to the final checkout list;
- Last but not least, there's an integration with a gRPC provided that can apply a discout percentage in the product list. (The .proto file is inside the project in the discount package).

## Sample uses

Given the following Database:

```
[{
  "id": 1,
  "title": "Ergonomic Wooden Pants",
  "description": "Deleniti beatae.",
  "amount": 15157,
  "is_gift": false
},
{
  "id": 2,
  "title": "Ergonomic Cotton Keyboard",
  "description": "Iste est.",
  "amount": 93811,
  "is_gift": true
}]
```

`/cart`: Using the request below:

   ```
    {
        "products": [
            {
                "id": 1,
                "quantity": 1
            }
        ]
    }
   ```



1. The `Normal` response: 
   ```
   {
    "total_amount": 15157,
    "total_amount_with_discount": 15157,
    "total_discount": 0,
    "products": [
        {
            "id": 1,
            "quantity": 1,
            "unit_amount": 15157, 
            "total_amount": 15157,
            "discount": 0, 
            "is_gift": false 
        }
    ]
   }
   ```
2. The `Black Friday` response:
   ```
   {
    "total_amount": 15157,
    "total_amount_with_discount": 15157,
    "total_discount": 0,
    "products": [
        {
            "id": 1,
            "quantity": 1,
            "unit_amount": 15157, 
            "total_amount": 15157,
            "discount": 0, 
            "is_gift": false 
        },
        {
            "id": 2,
            "quantity": 1,
            "unit_amount": 0, 
            "total_amount": 0,
            "discount": 0, 
            "is_gift": true 
        }
    ]
   }
   ```
3. The `with Discount(5%)` response:
   ```
   {
    "total_amount": 15157,
    "total_amount_with_discount": 14399,
    "total_discount": 758,
    "products": [
        {
            "id": 1,
            "quantity": 1,
            "unit_amount": 15157, 
            "total_amount": 15157,
            "discount": 758, 
            "is_gift": false 
        }
    ]
   }
   ```

# How to run

This can either be done through the `docker-compose`, docker image or by cloning the repo.

## Pre-requisites
The API uses some enviromnent variables, they are all provided in the **.env** file in the project repository. But in case you don't want to use the file, you can manually setup them or use their *default* values.

### Default values:
```
1. APP_PATH=/github.com/ChristianPrzybulinski #Used in Docker-Compose

2. DATABASE_PATH=./database #Used to define the database path for the JSON file

3. LOG_LEVEL=info #Used to define the log level

4. BLACK_FRIDAY= #Used to define the black friday day, empty means no Black Friday

5. API_PORT=8080 #The Port that the API will listen to

6. API_HOST= #The Host the API will listen to

7. DISCOUNT_SERVICE_PORT=50051 #The Port to connect in the Discount gRPC Service

8. DISCOUNT_SERVICE_HOST= #The Host to connect in the Discount gRPC Service

9. DISCOUNT_SERVICE_TIMEOUT=1 #The timeout time in seconds for each Discount request
```

**I highly recommend to use *docker-compose* and the *.env* file provided.**

### Through Docker-Compose:

I suggest to clone the repo, so you can have all files and project. But you only need the `docker-compose.yml` and the `.env` provided to start up the service. In case you setted the enviromnent variables manually, you won't need the `.env` file.

**Note:** In case you're using the default values for the enviromnent variables, you will need to edit the `docker-compose.yml` filling up the {APP_PATH} to any path you want, {API_PORT} to 8080 and {DISCOUNT_SERVICE_PORT} to 50051.

Starting the service: 

```
docker-compose pull go-cart-api

docker-compose up go-cart-api
```

You can also start the Discount Service with the same `docker-compose.yml` provided.

```
docker-compose pull discount-service

docker-compose up discount-service
```

### Through Docker:

```
docker run --env-file=.env christianprzybulinski/go-cart-api

docker run -p 50051:50051 hashorg/hash-mock-discount-service
```
You can skip the `--env-file` parameter in case you manually configured the enviromnent variables or want to use their default values.

### Through Repo:

**Note:** You won't have the Discount Service. Unless you start a Docker container with it. More importantly, you will need to have GO installed.

1. Clone the repo;
2. Inside the go-cart-api/src folder;
3. You can provide three arguments to run the application (overwriting env vars);
```
$1 = hostname:port, hostname, port
Example: localhost:8080

$2 = info, debug, warn, error
Example: debug

$3 = the database path
Example: /database
```

4. Run the application: 
```
go run . $1 $2 $3
Example: go run . 8080 debug
```

*Note: $GOPATH might need to be configured.*

# Sample request:

You can test the API Endpoint with any software that provides you a HTPP POST request to be made. I'll be showing the easiest one since you probably already have everything installed.

**Note:** The endpoint will be available at the `$API_HOST:$API_PORT/api/v1/cart`. In case you used the `.env` file or the default settings, It is going to be as the example below.

## cURL

I recommend using the sample files provided inside the *samples* folder.

```
curl -X POST -d "@samples/1.json" "localhost:8080/api/v1/cart"
```

**Final Note:** The database JSON name must be `products.json` in case you want to use your own database file.

## That's all folks. Thank you!
