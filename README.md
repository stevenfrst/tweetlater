# tweetlater

Project Architecture

```bash
.
├── appUtils
│   ├── appCookieStore
│   │   └── appCookie.go
│   ├── appHttpParser
│   │   └── appJsonParser.go
│   ├── appHttpResponse
│   │   ├── appResponse.go
│   │   └── jsonResponse.go
│   └── appStatus
│       └── appStatus.go
├── delivery
│   ├── authDelivery.go
│   ├── IDelivery.go
│   ├── tweetDelivery.go
│   └── userDelivery.go
├── go.mod
├── go.sum
├── infra
│   └── infra.go
├── main
│   ├── app.go
│   └── appRouter.go
├── manager
│   ├── repoManager.go
│   └── serviceManager.go
├── middleware
│   ├── logRequestMiddleware.go
│   └── tokenValidationMiddleware.go
├── models
│   ├── cred.go
│   ├── tweet.go
│   └── user.go
├── README.md
├── repository
│   ├── tweetRepository.go
│   ├── userAuthRepository.go
│   └── userRepository.go
└── usecases
    ├── tweetUseCase.go
    ├── userAuthUseCase.go
    └── userUseCase.go
```

### Usage

Run App

```go
go run tweetlater/main
```

### Handler

| Handler | Method | Usage|
|---|---|---|
| /login |  POST| Login To Pass MiddleWare |
| /logout| GET | Logout |
| /user | GET | Get user info with id |
|  |  POST| Register User |
|  | DELETE | Delete User |
|/user/upgrade  | POST | Upgrade User To Premium |
|  /app| GET | Get all Draft Tweets |
|  | POST | Post the draft tweets to the DB |
|  | DELETE |Delete draft tweets From DB  |
| /app/post | POST |Post draft Tweet to the DB(With Timer)  |
| /app/tweet | POST | Post tweet |

### .ENV

```bash
DBUSER=  
DBPASSWORD=  
DBHOST=  
DBPORT=  
DBSCHEMA=  
DBENGINE=  
HTTPHOST=  
HTTPPORT=  
CONSUMER_KEY=  
CONSUMER_SECRET=  
ACCESS_TOKEN=  
ACCESS_TOKEN_SECRET=
```

### TODO
