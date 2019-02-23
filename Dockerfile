FROM node:10.15-alpine

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go git curl

# install go dep manager
ENV GOPATH=/server
RUN mkdir -p /server/bin
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY server /server
WORKDIR /server/src/monde
RUN /server/bin/dep ensure

# build the render server
RUN mkdir /render-server
COPY render-server/typescript /render-server/typescript
COPY render-server/package.json /render-server/package.json
COPY render-server/package-lock.json /render-server/package-lock.json
WORKDIR /render-server
RUN npm install && npm run build

WORKDIR /

COPY deploy/scripts/docker /scripts
RUN chmod +x /scripts/run.sh
