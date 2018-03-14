MariaDB Tinsmith
===================

A _tinsmith_, in Blacksmith Services parlance, is a Cloud Foundry
Service Broker that runs as a Cloud Foundry Application, consumes
a provisioned, on-demand Blacksmith Service, and makes it
available to the rest of the world as a shared, multi-tenant
service.

A picture may make things clearer:

![Tinsmith Architectural Diagram](docs/tinsmith.png)

This Tinsmith, in particular, takes a dedicated MariaDB
instance (most likely deployed via the [MariaDB Blacksmith
Forge][mariadb-forge]) and carves up the dedicated host / cluster
into individual databases for each service.  Bound applications
will have free run of their own database, but will be unable to
interact with other shared databases on the same installation.

Deploying
---------

To deploy this Tinsmith, you need the code, a Cloud Foundry, and a
MariaDB service from your CF Marketplace.  We heartily recommend a
service deployed via the [Blacksmith][blacksmith] and its [MariaDB
Forge][mariadb-forge], but any service that is tagged `mariadb`,
and provides the `host`, `port` (as a number), `username` and
`password` fields in its `$VCAP_SERVICES` credentials block should
do.

```
git clone https://github.com/blacksmith-community/cf-mariadb-tinsmith
cd cf-mariadb-tinsmith

# push the code...
cf push --no-start
cf bind-service mariadb-tinsmith YOUR-DATABASE-SERVICE

# you may want to set some other environment variables at
# this stage; see "Configuration", below.
cf set-env mariadb-tinsmith SB_BROKER_USERNAME my-broker
cf set-env mariadb-tinsmith SB_BROKER_PASSWORD a-secret

# start the app
cf start mariadb-tinsmith

# register this tinsmith as a service broker in CF
cf create-service-broker mariadb-tinsmith my-broker a-secret \
  https://mariadb-tinsmith.$APP_DOMAIN
cf enable-service-access mariadb

# marvel at your handiwork
cf marketplace
```

Configuring
-----------

This tinsmith is configured entirely through environment
variables.

There are environment variables for governing the presentation of
this brokers service / plan in the marketplace:

- `$SERVICE_ID` - The internal ID of the service that this broker
  provides to the marketplace.
- `$SERVICE_NAME` - The CLI-friendly name of the service.
- `$PLAN_ID` - The internal ID of the plan that this broker
  provides to the marketplace.
- `$DESCRIPTION` - A human-friendly description of the service /
  plan, to be displayed in the marketplace
- `$TAGS` - A comma-separated list of tags to apply to instances
  of the service.

There are environment variables for controlling the security and
authentication parameters of the broker:

- `$SB_BROKER_USERNAME` - The HTTP Basic Auth username that must
  be used to access this broker.  Defaults to `b-mariadb`.
- `$SB_BROKER_PASSWORD` - The HTTP Basic Auth password that must
  be used to access this broker.  Defaults to `mariadb`.

You can also override the service selection logic and force it to
pick a specific, named service by setting the `$USE_SERVICE`
environment variable to its name.  Otherwise, the broker will look
for bound services that are tagged `mariadb`.




[mariadb-forge]: https://github.com/blacksmith-community/mariadb-forge-boshrelease
[blacksmith]: https://github.com/cloudfoundry-community/blacksmith
