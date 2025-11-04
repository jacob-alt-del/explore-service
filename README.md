# explore-service

## Run

### MySQL DB
```shell
cd _mysql
docker-compose up -d
```

### explore-service - local

```shell
go run ./cmd/server/main.go
```

### explore-service - docker
```shell
docker build . -t server
docker run -p 50051:50051 server
```

## Design

### SQL

```sql
USE explore;
CREATE TABLE users (
  id CHAR(36) NOT NULL,
  username VARCHAR(50) NOT NULL UNIQUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

Table for storing user ids as UUIDs. While not strictly necessary to implement the 4 gRPC endpoints, it was created to allow the use of foreign keys in the decisions talbe to ensure data consistency.

```sql
USE explore;
CREATE TABLE decisions (
  actor_id CHAR(36) NOT NULL,
  recipient_id CHAR(36) NOT NULL,
  liked BOOLEAN NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (actor_id, recipient_id),
  INDEX idx_recipient_liked (recipient_id, liked, updated_at DESC),
  INDEX idx_pair_recipient_actor (recipient_id, actor_id, liked),


  CONSTRAINT fk_actor FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_recipient FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

The main table required to support the functionality of the 4 gRPC endpoints.
- Composite PK to prevent duplicates
- Index idx_recipient_liked for ListLikedYou and ListNewLikedYou
- Index idx_pair_recipient_actor for JOIN on ListNewLikedYou
- Timestamps for pagination
- Foreign keys for data consistency

### Service

The service architecture was kept simple, consisting of two layers, a repository (dataaccess) and a service (service) layer.

The repositry layer encaptulated all of the database queries.

The service layer handled all of the business logic for each of the gRPC endpoints. 

Custom request validation was added to each endpoint as a layer of protection. Allows checking for things such as matching recipient_user_id and actor_user_id in the PutDecision endpoint which would result in invalid data.

## Testing

### Unit tests

Unit testing was done separately in each of the two layers. Testing was done only on the PutDecision endpoint due to time constraints but servs as an example of how it would be done at each layer.

Repository layer: using go-sqlmock to mock the MySQL connection.

Service layer: manualy stubbing out the function calls the the repository layer to allow purer unit tests.

### Smoketests

Smoketesting was done via the /cmd/client/main.go application which attempted to hit each endpoint using valid data from the MySQL database.

## Notes

- Potential flaw with ORDER BY updated_at DESC in ListLikedYou, a user could pass then relike to put them at the front of the other persons ListLikedYou list.

- When testing, noticed that the pagination can skip over rows with the same updated_at value. Potential fix could be to add an id column to the decisions table and using a combination of id and updated_at to form the pagination token.

```shell
protoc -I=. --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  explore/explore-service.proto
```