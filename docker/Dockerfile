FROM --platform=$BUILDPLATFORM golang:1.24.4 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/app
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/riverboat .

FROM cgr.dev/chainguard/bash:latest

# `nonroot` coming from distroless
USER 65532:65532

# Copy the binary
COPY --from=builder /go/bin/riverboat /bin/riverboat

# Run the server on container startup
ENTRYPOINT [ "/bin/riverboat" ]
CMD ["serve"]
