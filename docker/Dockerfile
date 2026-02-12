FROM --platform=$BUILDPLATFORM golang:1.26.0 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG BUILDTAGS

# Used when fetching trust center jobs, otherwise ignored
ENV GOPRIVATE=github.com/theopenlane/corejobs
# Ensure git uses SSH instead of HTTPS for private repos

RUN --mount=type=secret,id=id_ed25519 \
    if [ -n "$BUILDTAGS" ]; then \
    git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"; \
    mkdir -p -m 0700 /root/.ssh; \
    ssh-keyscan github.com >> /root/.ssh/known_hosts; \
    cp /run/secrets/id_ed25519 /root/.ssh/id_ed25519; \
    chmod 600 /root/.ssh/id_ed25519; \
    fi

WORKDIR /go/src/app
COPY . .

# Ensure modules are up to date before building, ignore errors when there are no build tags to possible unused dependencies
RUN if [ -n "$BUILDTAGS" ]; then \
    go mod tidy; \
    else \
    go mod tidy -e; \
    fi

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -tags=${BUILDTAGS} -o /go/bin/riverboat .

FROM cgr.dev/chainguard/bash:latest

# `nonroot` coming from distroless
USER 65532:65532

# Copy the binary
COPY --from=builder /go/bin/riverboat /bin/riverboat

# Run the server on container startup
ENTRYPOINT [ "/bin/riverboat" ]
CMD ["serve"]
