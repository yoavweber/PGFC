FROM golang:1.15-alpine AS build

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src/

# We want to populate the module cache based on the go.{mod,sum} files.
COPY ../go.mod go.sum /src/

#This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download && go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get -v
COPY ../.. /src/

RUN go build ./reciever/main.go


FROM alpine
COPY --from=build /src/. .

CMD ["./main"] --v
# CMD ["go run","/src/main.go"]
# ENTRYPOINT ["/usr/local/bin/main"]