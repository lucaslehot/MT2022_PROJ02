MT2022_PROJ02
====
#### ***Docker & Golang - FaaS implementation***
The aim of this project is to experiment FaaS (Function as a service), by implementing an infrastructure which permits to resize images with "on-demand" functionality.

# Project Status
*This repo is a [HETIC school](https://www.hetic.net/) project and its purpose is purely educational.* 

*Feel free to fork the project, but be aware that development might slow down or stop completely at any time, and that we are not looking for maintainers or owner.*

# Table of Contents
- [Project Status](#project-status)
- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [Deployment](#deployment)
- [Known Issues](#known-issues)
- [Built With](#built-with)
- [Team Members](#team-members)
- [Acknowledgments](#acknowledgments)
- [License](#license)

# Overview
This project should be carried out according to the following guidelines:

![alt text](./cloud.PNG/)
![alt text](./cloudless.PNG/)

We will be focusing on the second architecture. 

Our system is composed of :
- One app container
- One worker container
- One redis container based on redis:6.2 image
- One MySQL database container
- One volume

Our database is composed of a unique User table with the following columns :
- UserId
- CreatedAt
- UpdatedAt
- Username   
- Password   
- AvatarPath 

The typical happy path is the following :
- When an image is uploaded, it is stored in our volume.
- It's path is stored in the AvatarPath column of the current user.
- A 'generate_conversions' task is generated and published to the redis queue.

- The consumer, or worker container, then gets this task from the redis queue.
- It retrieves the user from the UserId found in the task's payload.
- It retrieves the image from the volume using it's path.
- It generates conversions using 	"github.com/nfnt/resize" package
- It stores the newly generated conversions in the volume.
- It notifies the queue that the task has been handled.

# Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) section for notes on how to deploy the project on a live system.

## Requirements
* [Golang](https://golang.org/dl/) (v1.15.8 or higher.)
* [rmq](https://github.com/adjust/rmq) Golang message queue system
* [gorm](https://gorm.io/gorm) Golang message queue system
* [resize](https://github.com/nfnt/resize) Golang image resizing package
* Docker installed and running
* Familiarity with basic Docker functionality and commands

## Deployment

In App:
```
go mod init
go mod tidy
go get
```

In Worker:
```
go mod init
go mod tidy
go get
```

From root:
```
docker-compose build
docker-compose up -d db redis-server
docker-compose up 
```

Web app is now available on port 8080.
Worker is on port 3000.
Redis-server is on port 6379.
Database is on port 3306.

# Known Issues
It is necessary to create a user with id = 1 first for the process to work

# Built With
* [Golang](https://golang.org/) - Open source programming language
* [Docker]() - Application packaging solution

# Team Members
* **Lucas Lehot** - [lucaslehot](https://github.com/lucaslehot)
* **Cyrille Banovsky** - [Ban0vsky](https://github.com/Ban0vsky)
* **Quentin Maillard** - [Tichyus](https://github.com/Tichyus)
* **Corentin Boulanouar** - [Shawnuke](https://github.com/Shawnuke)

# Acknowledgments
* 

# License
This project is licensed under the terms of the [MIT](https://opensource.org/licenses/MIT) license.
