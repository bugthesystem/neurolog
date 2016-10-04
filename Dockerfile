FROM redis

RUN buildDepsNR='git gcc libc6-dev make' \
    && set -x \ 
    && apt-get update && apt-get install -y  $buildDepsNR --no-install-recommends \
    && git clone https://github.com/antirez/neural-redis.git \
    && cd neural-redis \
    && make \
    && ls \
    && apt-get purge -y --auto-remove $buildDeps

CMD [ "redis-server", "--loadmodule", "neural-redis/neuralredis.so" ]