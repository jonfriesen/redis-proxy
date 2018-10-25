#!/bin/sh
# wait-for-redis.sh
set -e
CMD="$@"
HOST=redis
PORT=6379

echo "Waiting for Redis to warm up."
echo "Attempting to connect to Redis."
PONG=`redis-cli -h $HOST -p $PORT ping | grep PONG`
while [ -z "$PONG" ]; do
    sleep 3
    echo "Checking for life on Redis again."
    PONG=`redis-cli -h $HOST -p $PORT ping | grep PONG`
done
echo "Houston, we have Redis!"
exec $CMD
