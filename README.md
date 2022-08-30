# msg-scheduler

### v1.0.1 implementation details:
When a user registers, the API already schedule every stored email using the messaging service (SendGrid). This way, the API doesn't need to handle anything besides registering the user, selecting the messages to be sent and calling SendGrid with the payload.

Pros: it is fairly simple and the scheduling works just fine.
Cons: if we create a new message, there is no register of what was already sent to the user.

### v2.0.0 implementation details:
Now we register everything that is sent to each user, and we won't simply rely on SendGrid. There are two services, one API to handle the CRUD operations and one responsible for querying pending users, messages and sending emails. There is a simple retry mechanism in place, so if we miss one message, it'll be retried every minute. If we register a new message, all users that already received the rest and were marked as inactive will become active again.

Pros: more robust.
Cons: a few parts are not optimized for large scale usage on purpose, due to the time restriction to complete this task.
For example, sent message ids should be stored optimally. The approach used won't age well if the number of messages/users grows because it may overuse the DB and memory.

### How to run v1.0.1
### How to generate the specification

#### Without Docker
1. Generate specification by running `go install github.com/swaggo/swag/cmd/swag@v1.8.4` and
   `swag init --parseDependency --parseInternal --parseDepth 3`
2. Set `env=local` and `DBURL=something.sqlite` so you don't have to install PostgreSQL
3. Set `MSG_SERVICE=sendgrid` and place your SendGrid at `MSG_KEY`
4. Execute `go run main.go`

#### With Docker
1. Set `MSG_SERVICE=sendgrid` and place your SendGrid at `MSG_KEY` at `dev-postgres.env` file.
2. Run `docker-compose up`

### How to run v2.0.0
#### Without Docker
1. configure `dev.env` file by setting `env=local` and `DBURL=something.sqlite` so you don't have to install PostgreSQL
2. Set `MSG_SERVICE=sendgrid` and place your SendGrid at `MSG_KEY`
3. Run `api/main.go` and `scheduler/main.go`

#### With docker
1. Set `MSG_SERVICE=sendgrid` and place your SendGrid at `MSG_KEY` at `dev-postgres.env` file.
2. Run `docker-compose up`

