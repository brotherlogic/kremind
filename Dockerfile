# syntax=docker/dockerfile:1

FROM golang:1.23 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/kremind

COPY go.mod ./
COPY go.sum ./

RUN mkdir proto
COPY proto/*.go ./proto/

RUN mkdir server
COPY server/*.go ./server/

RUN mkdir db
COPY db/*.go ./db/

RUN mkdir runner
COPY runner/*.go ./runner

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /kremind

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /kremind /kremind

EXPOSE 80
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/kremind"]