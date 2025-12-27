FROM golang:1.25 AS build

WORKDIR /src
COPY go.* /src/
RUN go mod download

COPY ./ /src/
ARG DAP_VERSION
RUN CGO_ENABLED=0 go build -o dap -ldflags "-X main.version=${DAP_VERSION#v}" ./cmd/dap

FROM gcr.io/distroless/static

COPY --from=build /src/dap /

ENTRYPOINT ["/dap"]
