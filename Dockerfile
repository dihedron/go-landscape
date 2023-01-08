FROM ubuntu

LABEL Maintainer="Andrea Funtò <dihedron.dev@gmail.com>" Description="Landscape API container" Version="1.0"

# Add pgRouting launchpad repository
RUN \
    apt-get install -y software-properties-common && \
    add-apt-repository --update ppa:landscape/landscape-api && \
    apt-get update && \
    apt-get install landscape-api
