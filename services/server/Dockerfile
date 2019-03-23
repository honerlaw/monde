FROM node:10.15-alpine

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go git make

WORKDIR /base

# copy everythinng over
COPY scripts scripts
COPY src server/src
COPY .env server/.env
COPY assets server/assets
COPY Makefile server/Makefile
COPY renderer/scss renderer/scss
COPY renderer/typescript/src renderer/typescript/src
COPY renderer/typescript/types renderer/typescript/types
COPY renderer/typescript/tsconfig.json renderer/typescript/tsconfig.json
COPY renderer/package.json renderer/package.json
COPY renderer/package-lock.json renderer/package-lock.json

# build the go binary and remove the sources
WORKDIR /base/server/
ENV GOPATH=/base/server
ENV GIN_MODE=release
RUN make build
RUN rm -rf src/ Makefile

# install all dependencies, build, remove dependencies / sources, install only production dependencies
WORKDIR /base/renderer
ENV NODE_ENV=development
RUN npm install && npm run build
RUN rm -rf /base/renderer/typescript/src && rm -rf /base/renderer/typescript/tsconfig.json && rm -rf /base/renderer/typescript/types && rm -rf /base/renderer/typescript/tsconfig.json  && rm -rf /base/renderer/node_modules
ENV NODE_ENV=production
RUN npm install
RUN mv ../assets/* ../server/assets/

WORKDIR /base/server
RUN chmod +x /base/scripts/run.sh
ENTRYPOINT ["/base/scripts/run.sh"]
