FROM golang:1.17.1 as build

WORKDIR /go/src/wallet

COPY . .

RUN go mod download

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -o /go/bin/wallet ./

FROM alpine:3.13

RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/wallet /bin/
EXPOSE 8080

ENTRYPOINT ["wallet"]
