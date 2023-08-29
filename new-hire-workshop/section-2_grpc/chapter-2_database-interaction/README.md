# Chapter 2: Database Interaction

## Introduction to Databases

* [Confluence on Database credentials](https://confluence.eng.nimblestorage.com/display/CDSDEVOPS/How+To+Populate+Postgres+DBs+and+Vault+Policies+and+Onboarding+Secrets+in+a+DSCC+Cluster#HowToPopulatePostgresDBsandVaultPoliciesandOnboardingSecretsinaDSCCCluster-Databases)

## Databases in DSCC

* [Confluence on RDS](https://confluence.eng.nimblestorage.com/display/IFP/RDS)

## Exercise: Start a PostgreSQL Instance

* [Run PostgreSQL in Docker](https://hub.docker.com/_/postgres)

## Exercise: Connect to PostgreSQL in `orderstore`

* [Example from COSM](https://github.hpe.com/cloud/cloud-objectstore-manager/blob/master/internal/drivers/postgres/connection_provider.go)
* [Example from tilde-common](https://github.hpe.com/cloud/tilde-common/tree/master/pkg/pgdb)

## Exercise: Write SQL queries for gRPC endpoints

* Provide SQL schema with a few existing entries.
* Provide command to create, clean (and migrate?) database.
* Write SQL queries for both endpoints. 