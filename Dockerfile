FROM node:10.15-alpine

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go git

WORKDIR /base

# copy everythinng over
COPY deploy/scripts/docker scripts
COPY server server
COPY render-server/typescript/src render-server/typescript/src
COPY render-server/typescript/types render-server/typescript/types
COPY render-server/typescript/tsconfig.json render-server/typescript/tsconfig.json
COPY render-server/package.json render-server/package.json
COPY render-server/package-lock.json render-server/package-lock.json

# build the go binary and remove the sources
WORKDIR /base/server/src/package
ENV GOPATH=/base/server
ENV GIN_MODE=release
RUN go get -d -v
RUN go build -o /base/bin/server
RUN rm -rf /base/server

# build the render server
WORKDIR /base/render-server
ENV NODE_ENV=production
RUN npm install && npm run build

WORKDIR /base
RUN chmod +x /base/scripts/run.sh
ENTRYPOINT ["/base/scripts/run.sh"]
