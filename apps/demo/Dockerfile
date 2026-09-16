# syntax=docker/dockerfile:1

# --- build -------------------------------------------------------------------
FROM golang:1.25 AS build
ARG TARGETOS=linux
ARG TARGETARCH
WORKDIR /src

# Cache dependencies first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Cross-compile for the target platform (buildx sets TARGETARCH), so no emulation
# is needed for the Go build itself.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /out/demo .

# --- runtime -----------------------------------------------------------------
# Distroless static, nonroot (uid 65532) — matches the PSS "restricted" profile
# and the pod securityContext in the manifests.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/demo /demo
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/demo"]
