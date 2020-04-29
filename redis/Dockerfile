FROM redis:5.0.9-alpine

LABEL maintainer="https://gihub.com/popper2710"

COPY ./redis.conf /usr/local/etc/redis/redis.conf
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]