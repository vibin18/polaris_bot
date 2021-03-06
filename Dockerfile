FROM golang:1.18-alpine as build
RUN apk upgrade --no-cache --force
RUN apk add --update build-base make git
WORKDIR /go/src/github.com/vibin18/polaris_bot

# Compile
COPY ./ /go/src/github.com/vibin18/polaris_bot
RUN make dependencies
RUN make build
RUN ./polaris_bot --help

# Final Image
FROM gcr.io/distroless/static AS export-stage
COPY --from=build /go/src/github.com/vibin18/polaris_bot/polaris_bot /