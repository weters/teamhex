# teamhex - API for returning color data on sports teams

This project is the backend API for [Team Hex](https://teamhex.dev/). The API runs at the following URL: [https://api.teamhex.dev/](https://api.teamhex.dev/).

There is a companion project, [teamhex-vue](https://github.com/weters/teamhex-vue), that provides the frontend for Team Hex.

## API Documentation

This API is documented using Swagger 2.0 (aka OpenAPI 2.0). View the documentation here: [https://petstore.swagger.io/?url=https://api.teamhex.dev/swagger.json](https://petstore.swagger.io/?url=https://api.teamhex.dev/swagger.json)

The raw Swagger JSON can be found at the following URL: [https://api.teamhex.dev/swagger.json](https://api.teamhex.dev/swagger.json)

## Development

### Project Setup

```
go mod download
```

### Run the Tests

```
make test
```

### Run the development server

```
go run github.com/weters/teamhex/cmd/teamhexserver
```

Once running, you should be able to hit [localhost:5000](http://localhost:5000/)