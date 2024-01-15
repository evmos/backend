#!/bin/bash
if [[ -z "${REDIS_HOST}" ]]; then
  MY_HOST="redis"
else
  MY_HOST="${REDIS_HOST}"
fi

if [[ -z "${REDIS_PORT}" ]]; then
  MY_PORT="6379"
else
  MY_PORT="${REDIS_PORT}"
fi

while [ 1 == 1 ]
do
  nc -zv "${MY_HOST}" "${MY_PORT}"
  if [ $? == 0 ]; then
    break
  fi
  sleep 2
done

python -u $1.py
