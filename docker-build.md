# The instruction is for the onomy aurora-relayer images

* Build database:
```
docker build -t onomy/aurora-relayer-database --progress=plain --no-cache -f .docker/Dockerfile.database .
```

* Build endpoint:
```
docker build -t onomy/aurora-relayer-endpoint --progress=plain --no-cache -f .docker/Dockerfile.endpoint .
```