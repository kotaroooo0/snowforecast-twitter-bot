FROM nginx:1.19.0-alpine

RUN apk update \
    && apk upgrade \
    #openssl追加
    && apk add --no-cache openssl \
    #certbot追加
    && apk add --no-cache certbot \
    && apk add --no-cache bash \
    && certbot-auto certonly --standalone -d kotaroooo0.org -m admin@kotaroooo0.org --agree-tos -n\
    && certbot --nginx -d kotaroooo0.org

COPY ./nginx.conf /etc/nginx/nginx.conf
