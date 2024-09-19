FROM golang:1.23.1 AS builder

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /go/bin/riverboat -a -ldflags '-linkmode external -extldflags "-static"' .

FROM gcr.io/distroless/static:nonroot

# `nonroot` coming from distroless
USER 65532:65532

# Copy the binary
COPY --from=builder /go/bin/riverboat /bin/riverboat

# Run the server on container startup
ENTRYPOINT [ "/bin/riverboat" ]
CMD ["serve"]
