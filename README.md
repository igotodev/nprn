# NPRN is REST API with JWT authorization

## Authorization

All the methods in the API protected by JWT.

We must send `Authorization: Bearer <token>` in Header.

### POST

To create user and get authorization token we should use this request

```
/auth/sign-up
```

and send username, password and email, like that:

```
{
    "username": "name",
    "password": "pass",
    "email": "name@examle.com"
}
```

If all fields are correct we will get 200 OK and response with new token:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Mon, 31 Jan 2022 22:29:06 GMT
Content-Length: 189

{
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDM3MTE0NzYsImlhdCI6MTY0MzY2ODI3NiwidXNlcl9pZCI6IjYxZjNhZjI4NjViNWIzMjIyNDNhMDljNyJ9.ysaHAki74J0dQuTScU2YjChfXibXstkju3VQ-uu4dVQ"
}
```

If something went wrong we will get 401 Unauthorized

### GET

To get authorization token we should use this request

```
/auth/sign-in
```

and send username and password, like that:

```
{
    "username": "name",
    "password": "pass"
}
```

If user and password are correct we will get 200 OK and response with new token:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Mon, 31 Jan 2022 22:31:16 GMT
Content-Length: 189

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDM3MTE0NzYsImlhdCI6MTY0MzY2ODI3NiwidXNlcl9pZCI6IjYxZjNhZjI4NjViNWIzMjIyNDNhMDljNyJ9.ysaHAki74J0dQuTScU2YjChfXibXstkju3VQ-uu4dVQ"
}
```

If username or password is not correct we will get 401 Unauthorized


## Sales

### GET

`/api/v1/sale/` - get all sales

Response:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Mon, 31 Jan 2022 23:00:46 GMT
Content-Length: 347

[
  {
    "id": "61f867172c75ef87b9f4d040",
    "article": "12-223-41-33",
    "price_for_one": 222,
    "number_of_units": 1,
    "amount": 222,
    "date": "01-02-2022",
    "seller_id": "61f3af2865b5b322243a09c7"
  },
  {
    "id": "61f869ca2c75ef87b9f4d041",
    "article": "13-222-21-21",
    "price_for_one": 240.8,
    "number_of_units": 2,
    "amount": 240.8,
    "date": "01-02-2022",
    "seller_id": "61f3af2865b5b322243a09c7"
  }
]
```

If token is not correct we will get 401 Unauthorized with the error:

```
{
  "message": "unauthorized: not valid token"
}
```

or

```
{
  "message": "unauthorized: header is empty"
}
```

`/api/v1/sale/{id}` - get a sale

Response:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Mon, 31 Jan 2022 22:55:34 GMT
Content-Length: 174

{
  "id": "61f867172c75ef87b9f4d040",
  "article": "12-223-41-33",
  "price_for_one": 222.2,
  "number_of_units": 1,
  "amount": 222.2,
  "date": "01-02-2022",
  "seller_id": "61f3af2865b5b322243a09c7"
}
```

### POST

`POST /api/v1/sale/` - to add new sale

We need to send:

```
{
  "article":"13-222-21-21",
  "price_for_one": 240.8,
  "number_of_units": 1,
  "amount": 240.8,
  "date": "01-02-2022",
  "seller_id": "61f3af2865b5b322243a09c7"
}
```
Response:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Mon, 31 Jan 2022 22:57:22 GMT
Content-Length: 33

{
"id": "61f867172c75ef87b9f4d040"
}
```

### PUT

`PUT /api/v1/sale/{id}` - to update a sale

```
{
  "id": "61f867172c75ef87b9f4d040",
  "article": "12-223-41-33",
  "price_for_one": 222,
  "number_of_units": 1,
  "amount": 222,
  "date": "01-02-2022",
  "seller_id": "61f3af2865b5b322243a09c7"
}
```

Response:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Mon, 31 Jan 2022 22:57:22 GMT
Content-Length: 33

{
"id": "61f867172c75ef87b9f4d040"
}
```

### DELETE

`DELETE /api/v1/sale/{id}` - to delete a sale

Response:
```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Mon, 31 Jan 2022 23:03:15 GMT
Content-Length: 33

{
"id": "61f869ca2c75ef87b9f4d041"
}
```