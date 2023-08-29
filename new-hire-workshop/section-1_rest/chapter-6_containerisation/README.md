## What is containerization?
Containerization is the process of packaging an application and its dependencies into a single unit, known as a container. Containers are isolated from the host system and other containers running on the same system, which makes them a secure and portable way to deploy applications. Containers typically include the application code, runtime, system tools, libraries, and settings required to run the application.

## What are the benefits of containerization?
Containerization provides several benefits, including:

- Portability: 
Containers can be run on any system that supports the container runtime, making it easy to move applications between different environments.
- Consistency: 
Containers provide a consistent environment for the application to run in, regardless of the host system. This reduces the likelihood of errors caused by differences in the underlying infrastructure.
- Efficiency: 
Containers are lightweight and share the host system's resources, which makes them more efficient than virtual machines (VMs).
- Scalability: 
Containers can be scaled up or down quickly, allowing applications to respond to changing workloads.

### Containerization tools and technologies
There are several tools and technologies commonly used in containerization, including:

* Docker: Docker is a popular containerization platform that allows developers to build, ship, and run applications in containers. Docker provides a command-line interface (CLI) and a set of APIs for managing containers.
* Kubernetes: Kubernetes is an open-source container orchestration platform that automates the deployment, scaling, and management of containerized applications. Kubernetes provides a set of APIs for managing containers and a declarative configuration system for defining application resources.
* Containerd: Containerd is an open-source container runtime that provides a lightweight, high-performance environment for running containers. Containerd is used as the runtime engine for Docker and Kubernetes.
* CRI-O: CRI-O is an alternative container runtime that is designed to work with Kubernetes. CRI-O provides a lightweight, secure environment for running containers and is optimized for Kubernetes workloads.
* OCI: The Open Container Initiative (OCI) is a set of industry standards for container formats and runtime specifications. The OCI defines the format for container images and the runtime specification for running containers.


A Dockerfile is a set of commands that will be run to set up a container. Sometimes, it can help to think of it as the list of commands you’d need to run when setting up a brand new computer for development in a particular technology. In the case of this particular tutorial, we can think of it as the list of commands we’d need to run in order to set up a brand new computer for Go development. A project can have one or more Dockerfiles.

The docker-compose.yml file is a configuration file that will allow us to manage all our different containers. As we mentioned above, a project can have one or more Dockerfiles, which means it can be made up of one or more containers. The docker-compose.yml file can be thought of as a single project manager for all these containers.
## Conclusion
Containerization is a powerful tool for building and deploying applications in a consistent and portable way. By packaging applications in containers, developers can ensure that they run consistently across different environments, from development to production. There are several tools and technologies available for containerization, including Docker, Kubernetes, Containerd, CRI-O, and OCI.

## Commands to run the Docker image

- docker run -p 8080:8080 <image_name>

Replace <image_name> with the name you gave to the container image when you built it. The -p option maps port 8080 on your local machine to port 8080 in the container, so you can access the REST server at http://localhost:8080.

Note that if your Golang REST server listens on a different port than 8080, you'll need to adjust the -p option accordingly.


This docker-compose.yml file defines a single service called golang-rest-server. The build field specifies the path to the directory containing the Dockerfile, and Docker Compose uses the Dockerfile to build the container image for the service.

The ports field maps port 8080 on your local machine to port 8080 in the container, so you can access the REST server at http://localhost:8080.

The environment field defines environment variables that are passed to the Golang application running in the container.

To run the Golang REST server using Docker Compose, save the docker-compose.yml file to a directory on your computer, navigate to that directory in a terminal window, and run the following command:

- docker-compose up

This will start the container and output the logs from the Golang application to your terminal. To stop the container, press CTRL-C in the terminal window.



## Kubernetes

- Why you need Kubernetes and what it can do

Containers are a good way to bundle and run your applications. In a production environment, you need to manage the containers that run the applications and ensure that there is no downtime. For example, if a container goes down, another container needs to start. Wouldn't it be easier if this behavior was handled by a system?

That's how Kubernetes comes to the rescue! Kubernetes provides you with a framework to run distributed systems resiliently. It takes care of scaling and failover for your application, provides deployment patterns, and more. For example: Kubernetes can easily manage a canary deployment for your system.

Kubernetes provides us with:

- Service discovery and load balancing Kubernetes can expose a container using the DNS name or using their own IP address. If traffic to a container is high, Kubernetes is able to load balance and distribute the network traffic so that the deployment is stable.
- Storage orchestration Kubernetes allows you to automatically mount a storage system of your choice, such as local storages, public cloud providers, and more.
- Automated rollouts and rollbacks You can describe the desired state for your deployed containers using Kubernetes, and it can change the actual state to the desired state at a controlled rate. For example, you can automate Kubernetes to create new containers for your deployment, remove existing containers and adopt all their resources to the new container.
- Automatic bin packing You provide Kubernetes with a cluster of nodes that it can use to run containerized tasks. You tell Kubernetes how much CPU and memory (RAM) each container needs. Kubernetes can fit containers onto your nodes to make the best use of your resources.
- Self-healing Kubernetes restarts containers that fail, replaces containers, kills containers that don't respond to your user-defined health check, and doesn't advertise them to clients until they are ready to serve.
- Secret and configuration management Kubernetes lets you store and manage sensitive information, such as passwords, OAuth tokens, and SSH keys. You can deploy and update secrets and application configuration without rebuilding your container images, and without exposing secrets in your stack configuration.


One can refer to the officical documentation for extensive knowledge on kubernetes:

- Kubernetes official documentation: https://kubernetes.io/docs/home/



You can use minikube for easy learning and developing for kubernetes locally on your development machine. Minikube can be used to experiment with Kubernetes features, test application deployments, and develop and debug Kubernetes applications locally. Minikube runs on a virtual machine or a container runtime, such as VirtualBox or Docker, and supports multiple operating systems, such as Windows, macOS, and Linux.

Installation steps are provided in the document below:
https://minikube.sigs.k8s.io/docs/start/

Once the installation is complete create a minikube cluster with the below command.

- minikube start


Follow the link below for various command related to deployment and administration

https://kubernetes.io/docs/tutorials/hello-minikube/


## Helm:



Helm is a package manager for Kubernetes that allows you to install, upgrade, and manage Kubernetes applications as packages, called charts. Helm provides a command-line tool that allows you to create, package, and share charts, and a repository for hosting and sharing charts. Helm charts can include multiple Kubernetes resources, such as deployments, services, and config maps, and can be customized using templates and values files. Helm makes it easier to manage the complexity of deploying and managing Kubernetes applications by providing a standard way to package and distribute applications.

- Helm official documentation: https://helm.sh/docs/
- Getting started with Helm tutorial: https://helm.sh/docs/intro/getting-started/
- Helm in Action book by Ryan Riches: https://www.manning.com/books/helm-in-action
