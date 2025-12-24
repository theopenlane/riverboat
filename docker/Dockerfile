FROM --platform=$BUILDPLATFORM golang:1.25.5 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG BUILDTAGS

# Used when fetching trust center jobs, otherwise ignored
ENV GOPRIVATE=github.com/theopenlane/corejobs
# Ensure git uses SSH instead of HTTPS for private repos
RUN git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"

WORKDIR /go/src/app
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -tags=${BUILDTAGS} -o /go/bin/riverboat .

FROM cgr.dev/chainguard/bash:latest

# `nonroot` coming from distroless
USER 65532:65532

# Copy the binary
COPY --from=builder /go/bin/riverboat /bin/riverboat

# Run the server on container startup
ENTRYPOINT [ "/bin/riverboat" ]
CMD ["serve"]