# syntax=docker/dockerfile:1

FROM node:25.6.0-alpine AS node_builder

WORKDIR /src

ENV PNPM_HOME="/pnpm"
ENV PATH="${PNPM_HOME}:${PATH}"
ENV CI=1

RUN npm install --global --force corepack@latest && corepack enable && pnpm --help

COPY --link ./package.json ./pnpm-lock.yaml ./
COPY --link ./assets ./assets
COPY --link ./internal/web/templates ./internal/web/templates

RUN \
  --mount=type=cache,id=pnpm,target=/pnpm/store \
  pnpm install --frozen-lockfile \
  && pnpm run build


FROM golang:1.25.6 AS go_builder

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY --link . .
COPY --from=node_builder /src/internal/web/assets ./internal/web/assets
RUN CGO_ENABLED=0 go build -o build/void-tool main.go


FROM scratch

COPY --from=go_builder /src/build/void-tool /

EXPOSE 8080

ENTRYPOINT [ "/void-tool", "serve" ]

# TODO: health
