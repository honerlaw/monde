FROM node:10.15-alpine

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go git

WORKDIR /base

# copy everythinng over
COPY deploy/scripts/docker scripts
COPY server server
COPY render-server/scss render-server/scss
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
RUN go build -o /base/server/server
RUN rm -rf /base/server/pkg /base/server/src

# install all dependencies, build, remove dependencies / sources, install only production dependencies
WORKDIR /base/render-server
ENV NODE_ENV=development
RUN npm install && npm run build
RUN rm -rf /base/render-server/typescript/src && rm -rf /base/render-server/typescript/tsconfig.json && rm -rf /base/render-server/typescript/types && rm -rf /base/render-server/typescript/tsconfig.json  && rm -rf /base/render-server/node_modules
ENV NODE_ENV=production
RUN npm install

WORKDIR /base/server
RUN chmod +x /base/scripts/run.sh
ENTRYPOINT ["/base/scripts/run.sh"]
