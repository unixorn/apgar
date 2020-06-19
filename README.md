# Apgar

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/unixorn/apgar.svg)](https://github.com/unixorn/apgar/stargazers)
[![Code Climate](https://codeclimate.com/github/unixorn/apgar/badges/gpa.svg)](https://codeclimate.com/github/unixorn/apgar)
[![Issue Count](https://codeclimate.com/github/unixorn/apgar/badges/issue_count.svg)](https://codeclimate.com/github/unixorn/apgar)

# Design Goals

We wanted a quick, simple and standardized way of doing health checks for the various services in our environment.

Apgar walks a directory tree (by default `/etc/apgar/healthchecks`), runs the healthCheck scripts it finds there in parallel, to keep the run time as short as possible, and aggregates the results into a directory (`/var/lib/apgar` by default).

The status directory is then served by a simple standalone web server so the results can be used as health checks by Amazon Load Balancers and Auto Scaling Groups.

## Details

Apgar consists of two parts, `apgar-server` which serves the health information, and `apgar-probe` which collects & aggregates the individual server health checks.

`apgar-probe` runs the health checks in parallel, and will report failed status as soon as any of them fail - it does _not_ wait until all checks are complete to report failure. Status is scrapable at `hostname:9000/status`, but you can override the port in `config.toml`.

# Writing Health Check Scripts

An apgar health check must be:

* Executable, with a shebang line.
* It should either exit 0 if the check pases, or exit with anything other than 0. Note that `apgar-probe` relies on the exit code of the check to determine OK/FAIL, _not_ any text output of the health check script.
* It must be named with a **.healthCheck** suffix
* The check should never print anything else unless `--verbose` is passed on the command line. If called with `--quiet`, it may print nothing at all, but _must still return zero or non-zero._
* To minimize check time, `apgar-probe` walks the healthchecks directory and immediately runs any check scripts it finds in parallel as soon as it finds them, so check scripts _must not assume that they will be run in a particular order_, or that other check scripts will *not* be running simultaneously with them.
* Checks should be fast - since `apgar-probe` will run all the checks in parallel, it is better to have 3 separate tests that each run in N milliseconds than one test that runs in 3N milliseconds.
* Checks should be idempotent
* Checks _must_ be non-destructive and _must not_ change the state of the underlying service - by definition, they will be run at regular intervals while the service they're checking is in production.

# Packaging

To make it easy to install on both Debian and CentOS based systems, The included Rakefile can build both deb and rpm files - `rake deb` will build a deb file, and `rake rpm` will build a rpm file. This requires `rake` and `bundler`, but only on your build machine, not on machines you're going to install Apgar on.

# FAQ

## Why is it named Apgar?

Virginia Apgar invented the Apgar score as a method to quickly summarize the health of newborn children. This seemed like an appropriate name for a quick health check system.

## Can I run one test at a time?

No. `apgar-probe` is designed to run all the health check scripts it finds in parallel. This allows it to fail the health check as fast as possible - the longest time that `apgar-probe` will take to determine a machine has failed its health check is the time it takes to run the slowest health check script. You should break up your check scripts into small scripts that each check one aspect of your service instead of large scripts that sequentially test multiple parts of your services.

## My scripts have to be run in a particular order, how do I specify that?

You don't. The work around is to have a single large check script that runs multiple tests in the order you want, but that will slow down the `apgar` run.

## Why do you run your own web server instead of piggybacking on a system's existing web server?

You are of course free to use your own webserver instead of `apgar-server`, but we opted to provide a stand-alone server for the following reasons:

* Nobody wants to have to maintain configuration files for every webserver out there
* Not every system already has a webserver installed
* Using a standalone server helps minimize Apgar's impact on existing services. Both `apgar-server` and `apgar-probe` are written in golang so that using Apgar doesn't pull in any dependencies that might conflict with those needed by the services you actually care about on a given system.
* Using a standalone server keeps you from having to alter your existing webserver's configuration to use Apgar.
* Our server is deliberately minimal. It barely uses any resources and all it does is serve up files. It doesn't act on user input, so it can't be exploited by malformed user input.

## Why Golang?

I wanted Apgar to have as little impact on the host system as possible. Go gives us static binaries, so we don't have to worry about dependency conflicts with other services on machines. And it was a good excuse to start learning Go.
