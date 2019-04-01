FROM node:10.15-alpine

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go git make

WORKDIR /base

# copy everythinng over
COPY src server/src
COPY .env server/.env
COPY assets server/assets
COPY Makefile server/Makefile
COPY renderer/scss server/renderer/scss
COPY renderer/typescript/src server/renderer/typescript/src
COPY renderer/typescript/types server/renderer/typescript/types
COPY renderer/typescript/tsconfig.json server/renderer/typescript/tsconfig.json
COPY renderer/package.json server/renderer/package.json
COPY renderer/package-lock.json server/renderer/package-lock.json

# build the go binary and remove the sources
WORKDIR /base/server/
ENV GOPATH=/base/server
ENV GIN_MODE=release
RUN make build
RUN rm -rf src/ Makefile

# install all the renderer dependencies, and build the bundles, then remove the sources
WORKDIR /base/server/renderer
ENV NODE_ENV=development
RUN npm install && npm run build
RUN rm -rf /base/server/renderer

WORKDIR /base/server
ENTRYPOINT ["/base/server/bin/server"]
