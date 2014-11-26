# Dancible

## What

In it's simplest form, a program to generate Dockerfiles that:

* Install development Ansible
* Pull down a repository containing an ansisble playbook
* Uses the playbook to do all the heavy lifting and configure the container.

In the long term this will not only generate the Dockerfile but also build the docker image with the option to 
publish/upload to a docker registry.

## How

Clone this repository and build 

    $ go build

Run the executable

    $ ./dancible -help
        Usage of ./dancible:
          -branch="": a branch to be checked out from playbook repo containing 'site.yml'
          -name="": name of docker container to be produced
          -os="ubuntu": the base operating system to use in the container (ubuntu, centos, debian)
          -repo="": git URL to pull an ansible playbook and configure this container
          -version="latest": the version of the docker image to use

The only required parameter is the repository to pull down the ansible-playbook.

By default the Dockerfile will be created from the latest ubuntu image.

### Example

Create a Dockerfile to deploy the (awesome) sovereign on Debian

    $ ./dancible -os debian -repo git@github.com:al3x/sovereign.git

## Why

I love Docker and I love Ansible. 

It seemed silly not to use simple Ansible Playbooks (of which there are loads already created) to 
configure containers, hopefully avoiding large, complex Dockerfiles. 

