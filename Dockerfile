# syntax=docker/dockerfile:1

FROM golang:1.25.6 AS builder

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o build/void-tool main.go


FROM scratch

COPY --from=builder /src/build/void-tool /

EXPOSE 8080

ENTRYPOINT [ "/void-tool", "serve" ]

# TODO: health
