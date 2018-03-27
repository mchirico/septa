#!/bin/bash
# Current directory only
# git archive --format=tar --prefix=septa/ HEAD ../| gzip > septa.tar.gz
docker build --no-cache -t mchirico/septa .
if [ "$(whoami)" != "mchirico" ]; then
        echo "You need permissions to push to repo"
        exit 1
fi
docker run --rm -it mchirico/septa septa
docker push mchirico/septa

