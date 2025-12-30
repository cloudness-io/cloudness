# ---------------------------------------------------------#
#                     Build web image                      #
# ---------------------------------------------------------#
FROM node:24 AS web

WORKDIR /app

# Copy package files first for better caching
COPY package.json package-lock.json* ./

RUN npm ci --prefer-offline || npm i

# Copy full web assets for tailwind to parse templ files
COPY ./app/web ./app/web

RUN npx esbuild ./app/web/assets/index.js --bundle --outdir=./dist
RUN npx @tailwindcss/cli -i ./app/web/assets/app.css -o ./dist/styles.css

# ---------------------------------------------------------#
#                     Cert image                           #
# ---------------------------------------------------------#
FROM alpine:latest AS cert-image

RUN apk --update add ca-certificates


# ---------------------------------------------------------#
#                   Build Cloudness image                  #
# ---------------------------------------------------------#
FROM golang:1.25-trixie AS builder

# Build arguments for versioning
ARG GIT_COMMIT
ARG CLOUDNESS_VERSION_MAJOR=0
ARG CLOUDNESS_VERSION_MINOR=0
ARG CLOUDNESS_VERSION_PATCH=0
ARG CLOUDNESS_VERSION_PRE

# Setup workig dir
WORKDIR /app

# Get dependencies first - cached if go.mod/go.sum unchanged
COPY go.mod go.sum ./

ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN go install tool
RUN go mod download

# Copy templ files first for templ generate caching
COPY ./app/web ./app/web
COPY --from=web /app/dist ./app/web/public/assets

# Generate the template code (cached if templ files unchanged)
RUN go tool templ generate

# Copy remaining source code
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    LDFLAGS="-X github.com/cloudness-io/cloudness/version.GitCommit=${GIT_COMMIT} -X github.com/cloudness-io/cloudness/version.major=${CLOUDNESS_VERSION_MAJOR} -X github.com/cloudness-io/cloudness/version.minor=${CLOUDNESS_VERSION_MINOR} -X github.com/cloudness-io/cloudness/version.patch=${CLOUDNESS_VERSION_PATCH} -X github.com/cloudness-io/cloudness/version.pre=${CLOUDNESS_VERSION_PRE} -extldflags '-static'" && \
    CGO_ENABLED=1 \
    CC=$CC go build -ldflags="$LDFLAGS" -o ./cloudness ./cmd/app

# ---------------------------------------------------------#
#                   Create final image                     #
# ---------------------------------------------------------#
FROM scratch AS final

# setup app dir and its content
WORKDIR /app

COPY --from=builder /app/cloudness /app/cloudness
COPY --from=cert-image /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8000
ENTRYPOINT [ "/app/cloudness", "server" ]