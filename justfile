# Build the Docker image
build:
    docker build -t ghcr.io/simmsb/neat .

push:
    docker push ghcr.io/simmsb/neat
