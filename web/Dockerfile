FROM node:12.2.0-alpine as node_cache
WORKDIR /root/app
COPY package.json .

ENV NPM_CONFIG_LOGLEVEL warn
ENV NPM_CONFIG_REGISTRY https://registry.npm.taobao.org

RUN yarn install --loglevel notice

COPY . .
RUN yarn run build

# 
FROM nginx:1.16.0-alpine

LABEL MAINTAINER="Colynn Liu <colynn.liu@gmail.com>"

ADD ./deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=node_cache /root/app/dist /usr/share/nginx/html

