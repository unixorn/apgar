# Apgar

[![GitHub stars](https://img.shields.io/github/stars/unixorn/apgar.svg)](https://github.com/unixorn/apgar/stargazers)
[![Code Climate](https://codeclimate.com/github/unixorn/apgar/badges/gpa.svg)](https://codeclimate.com/github/unixorn/apgar)
[![Issue Count](https://codeclimate.com/github/unixorn/apgar/badges/issue_count.svg)](https://codeclimate.com/github/unixorn/apgar)

# Design Goals

We wanted a quick, simple and standardized way of doing health checks for the various services in our environment.

Apgar walks a directory tree (by default `/etc/apgar/healthchecks`), runs the healthCheck scripts it finds there in parallel, to keep the run time as short as possible, and aggregates the results into a directory (`/var/lib/apgar` by default).

The status directory is then served by a simple standalone web server so the results can be used as health checks by Amazon Load Balancers and Auto Scaling Groups.

## Details

Apgar consists of two parts, `apgar-server` which serves the health information, and `apgar-probe` which collects & aggregates the individual server health checks.

`apgar-probe` runs the health checks in parallel, and will report failed status as soon as any of them fail - it does _not_ wait until all checks are complete to report failure.

# Writing Health Check Scripts

An apgar health check must be:

* Executable, with a shebang line.
* It should either write "OK" to console _and_ return 0, or write "NOT OK" to console _and_ return anything other than 0. Note that `apgar` relies on the exit code of the check to determine OK/FAIL, _not_ any text output of the health check script.
* It must be named with a **.healthCheck** suffix
* The check should never print anything else unless _--verbose_ is passed on the command line. If called with _--quiet_, it may print nothing at all, but _must still return zero or non-zero._
* Check scripts _must not assume that they will be run in a particular order_, or that other check scripts will *not* be running simultaneously with them. To minimize check time, Apgar runs the check scripts in parallel as soon as it finds them.
* Checks should be fast - since apgar will run all the checks in parallel, it is better to have 3 separate tests that each run in N seconds than one test that runs in 3N seconds.
* Checks should be idempotent
* Checks must be non-destructive and not change the state of the underlying service - by definition, they will be run while the service they're checking is in production.

# FAQ

## Why Apgar?

Virginia Apgar invented the Apgar score as a method to quickly summarize the health of newborn children. Seemed like an appropriate name for a quick health check system.

## Can I run one test at a time?

No. Apgar is designed to run all the health check scripts it finds in parallel. This allows it to fail the health check as fast as possible - the longest time that Apgar will take to report a failed health check is the time it takes to run the slowest health check script. You should break up your check scripts into small scripts that each check one aspect of your service instead of large scripts that test multiple parts of your services.

## My scripts have to be run in a particular order, how do I specify that?

You don't. The current work around is to have a single large check script that runs multiple tests in the order you want, but that will slow down the apgar run.

## Why do you run your own web server instead of piggybacking on a system's existing web server?

You are of course free to use your own webserver instead of `apgar-server`, but we opted for a stand-alone server for the following reasons:

* Nobody wanted to have to maintain configuration files for every webserver out there
* Not every system already has a webserver installed
* Using a standalone server helps minimize Apgar's impact on existing services. Both `apgar-server` and `apgar-probe` are written in golang so that using Apgar doesn't pull in any dependencies that might conflict with those needed by the services you actually care about on a system.
* Using a standalone server keeps you from having to alter your existing webserver's configuration to use Apgar.

## Why Go?

I wanted Apgar to have as little impact on the host system as possible. Go gives us static binaries, and it was a good excuse to start learning Go.

## How can I package this so I don't have to build it on every server?

If you have bundler installed, `rake deb` will build a deb, and `rake rpm` will build a rpm.
