# PGFS - ParaGliding File System
## Project Description
This project is built as (self)tasked in the course PROG2005 Cloud Technologies for the Spring 2021 Semester. The project is suggested to be available as open source, for further work from external contributors after the project period has ended.
This application is based on the IPFS-protocol (InterPlanetary File System), and is intended to be used by paragliders to share flight-info with each other, in a completely decentralized manner. The application in the current state is entirely backend-focused and is written in golang.
Paraglider-flight files (in standardized format ".igc") are intended for sharing on a peer-to-peer basis.

This program utilizes a specialized version of the go-implementation of IPFS - located at [go-pgfs](https://github.com/yoavweber/go-pgfs) <br>
All licenses and dependencies from the original project this repository has been forked from ([go-ipfs](https://github.com/ipfs/go-ipfs)) have been preserved.

Further Project Progress Documentation is located in the Project Wiki.

## Features
The program is directly utilizing a specialized version of IPFS (found at [go-pgfs](https://github.com/yoavweber/go-pgfs)), where all features from the go-implementation of IPFS ([go-ipfs](https://github.com/ipfs/go-ipfs)) are preserved, yet some modified to fit the project case. Additional features include:
* Uploading and downloading files is exclusive for files of the .igc-format (standardized file-format for paragliding flight-information)
* Running the program through docker-compose initiates 3 PGFS-nodes. One of the nodes is the server-node, and the other two exist as clients. One sender and one reciever. This dockerization is primarily for showing off the proof-of-concept of this application.
* Mother-node (main server bootstrap-node) has been thoroughly tested to be working in a server/cloud-situation. This required a lot of testing and fine-tuning of the program. It now works as intended. The bootstrap node is directly interchangeble in main.go for all nodes.

## Deployment
This program is initialized by docker-compose. The compose-file initiates three containers, each containing a node initialized to the pgfs-network. These nodes are of different classifications. One node is initialized as a server node (and is the main bootstrap), and the two other client-nodes. The two client nodes serve different purposes, one of them aims to upload a valid IGC-file, the other client-node fetches this file from the network.

## Project Assessment
This project is scheduled for OpenSource Development after the project-period period has ended. 
For assessing project-work in the course of Cloud Technologies, please visit the Wiki. In the wiki, you will find important documents from the development-process, as well as comprehensive project reports. The "Project Report" will serve as the main assessment-point for the Spring 2021 "Cloud Technologies" project (due to the major differences in development, contra the more streamlined development of REST-applications in Go). Read the Project Report in the wiki to assert our project. Please defer to the wiki for further information.

<br>Authored by<br>
<b>Yoav Weber</b><br>
<b>Milosz Antoni Wudarczyk</b><br>
<b>Kristian Amundsen Øhman-Norén</b><br>
2021, Norwegian University of Science and Technology
