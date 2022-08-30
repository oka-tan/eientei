# Eientei Stack

## Migration Testing

* Set the replica count for kaguya and moon to zero.
* Set the replica count for mysql and neofuuka to one.
* Call `docker compose up`
* Give it a bit for neofuuka to scrape whatever board is configured.
* Stop docker compose.
* Set the replica count for neofuuka to zero.
* Set the replica count for postgres, s3, lnx, envoy, mokou, reisen and moon to 1.
* Call `docker compose up` again.
* Give it a bit for mokou to migrate data to postgres and S3.
* Go to localhost:9001 and set the ayase bucket as public in the minio console.
* Check reisen.localhost for the current stack frontend (ugly) presenting the migrated data.

## Actual Migration Usage

Essentially, configure mokou.json to point it to your actual asagi database and your actual asagi image/thumbnail folder and then run it.

## Actual Deployment

Not yet. Also, don't deploy docker-compose to production.

## Kubernetes

Half-assed kubernetes deployment can be done with

```sudo sh bootstrap-kubernetes.sh```

