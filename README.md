# teamhex - RESTful API for returning color data on sports teams

This project is the backend RESTful API for [TeamHex](https://teamhex.dev/). The API runs at the following URL: [https://api.teamhex.dev/](https://api.teamhex.dev/).

There is a companion project, [teamhex-vue](https://github.com/weters/teamhex-vue), that provides the frontend for TeamHex.

## Project Setup

```
go mod download
```

## Run the development server

```
go run github.com/weters/teamhex/cmd/teamhexserver
```

Once running, you should be able to hit [localhost:5000](http://localhost:5000/)