FROM golang:latest AS build-container
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
ARG version
RUN go test ./...
RUN GOOS=linux \
    CGO_ENABLED=0 \
    go build \
        -ldflags "-X main.Version=$version" \
        -o teamhexserver github.com/weters/teamhex/cmd/teamhexserver
RUN go get github.com/go-swagger/go-swagger/cmd/swagger
RUN go install github.com/go-swagger/go-swagger/cmd/swagger
RUN swagger generate spec -o swagger.json

FROM alpine:latest AS build-json
WORKDIR /build
RUN apk add jq
COPY configs/teamhex.json teamhex-bloated.json
RUN cat teamhex-bloated.json | jq . -c > teamhex.json

FROM alpine:latest
WORKDIR /app
COPY --from=build-container /build/teamhexserver /bin/teamhexserver
COPY --from=build-container /build/swagger.json .
COPY --from=build-json /build/teamhex.json configs/teamhex.json
ENTRYPOINT [ "/bin/teamhexserver" ]
