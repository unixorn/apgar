# Goal

We wanted a quick, simple and standardized way of doing health checks for the various services in our environment. Apgar walks a directory tree (by default `/etc/apgar/healthchecks`), runs the healthCheck scripts it finds there and aggregates the results into a directory (`/var/lib/apgar` by default). That directory is then served by a simple standalone web server so the results can be used as health checks by Amazon Load Balancers and Auto Scaling Groups.

## Details

Apgar consists of two parts, apgar-server which serves the health information, and apgar-probe which collects & aggregates the individual server health checks.


# FAQ

## Why Apgar?

Why not Apgar? Virginia Apgar invented the Apgar score as a method to quickly summarize the health of newborn children. Seemed appropriate for a quick health check system.

## Why not just piggyback on a system's existing web server?

* I don't want to have to maintain configuration files for every webserver out there
* Not every system has a webserver installed
* Using a standalone server helps minimize Apgar's impact on existing services. The apgar-server component is written in golang so that it doesn't pull in any dependencies that might conflict with those needed by the services you actually care about on the system.

## Why Go

I wanted Apgar to have as little impact on the host system as possible. Go gives us static binaries, and it was a good excuse to start learning Go.
