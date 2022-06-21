# Eientei

## Kaguya

4chan scraper.

Designed mostly around making as few HTTP requests as possible while shirking other performance considerations. Should not be used to scrape fast boards.

Scrapes the archive but disregards further activity in archived threads, i.e. late post deletions.

Requires S3 access or some equivalent, i.e. minio.

Board partitions should be created manually. Just edit the schema.sql to add them in.

