FROM golang
ADD . /code
WORKDIR /code
RUN apt-get update
RUN apt-get -y upgrade
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends apt-utils
RUN apt-get install -y net-tools



# Build
# docker build --no-cache -t mchirico/septa .

# Upload to gsutil
# docker push mchirico/septa

# To get this file
# docker pull mchirico/septa


