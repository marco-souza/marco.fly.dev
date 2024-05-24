#! /bin/sh
file=$(cat .env | grep =)
for line in $file; do
  echo "Applying secret to fly.io"
  fly secrets set $line
done
fly secrets deploy
fly secrets list
