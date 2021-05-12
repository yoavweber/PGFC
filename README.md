# PGFS - ParaGliding File System
## Project Description
This project is built as (self)tasked in the course PROG2005 Cloud Technologies for the Spring 2021 Semester. The project is suggested to be available as open source, for further work from external contributors after the project period has ended.
This application is based IPFS (InterPlanetary File System), and is intended to be used by paragliders to share flight-info with each other, completely decentralized. The application in the current state is entirely backend-focused and is written in golang.
Paraglider-flight files (in standardized format ".igc") are intended for sharing on a peer-to-peer basis.

This program utilizes a specialized version of the go-implementation of IPFS - located at [go-pgfs](https://github.com/yoavweber/go-pgfs)<br>
All licenses and dependencies from the original project this repository has been forked from ([go-ipfs](https://github.com/ipfs/go-ipfs)) have been preserved.

Further Project Progress Documentation is located further down in this README.

## Features
The program is directly utilizing a specialized version of IPFS (found at [go-pgfs](https://github.com/yoavweber/go-pgfs)), where all features from the go-implementation of IPFS ([go-ipfs](https://github.com/ipfs/go-ipfs)) are preserved, yet modified to fit the project case. Additional features (as described in the project description) include:
* Uploading and downloading files is exclusive for files of the .igc-format (standardized file-format for paragliding flight-information)
* Running the program through docker-compose initiates 3 PGFS-nodes for testing-purposes
* Mother-node (main server bootstrap-node) is located on the internal NTNU-network. In practice, this means the program only works internally in the NTNU institutional network. To use your own server-node (for further development purposes), the bootstrap-node must be updated to coincide with your own server-node.

## Deployment
This program is initialized by docker-compose. The compose-file initiates three nodes within the pgfs-network. These nodes are of different classifications. One node is initialized as a server node (and is the main bootstrap), and the two other client-nodes. The two client nodes serve different purposes, one of them aims to upload a valid IGC-file, the other client-node fetches this file from the network.

## The Development Process (for assessment purposes)
### Original Project Plan
The project idea consisted mainly of interrogating new technologies related to Cloud Technologies. Taking inspiration from a brief introduction in a lecture, IPFS was a topic that interested us greatly. We sought after getting to know this protocol better, on a deep and advanced level. The project seemed like a golden opportunity to explore more about this new type of technology. Learn about the IPFS-protocol and its functionality stood as a pillar for this project. The main plan consisted of:
1. IPFS - Getting to know the protocol and its functionality
2. IPFS - Interrogating the go-implementation of IPFS and customizing it for optimized functionality.
3. Developing an OpenSource, completely decentralized, application, based on a specialized version of IPFS and its libraries, that will be used for sharing files between Paragliders. The idea of decentralization is forward-thinking, as we believe the future of the internet, and thus file-sharing, will be in a decentralized manner. IPFS and its protocols serve as a decentralized way of sharing files. This application-idea is inspired and pitched by Mariusz Nowostawski by NTNU Gjøvik.
### What went well
The group has had an incredible amount of learning-experience from working with an entirely foreign prospect for our project. The learning-curve has been steep, yet enriching. The program serves as a "proof of concept", but what has been developed has been thorougly thought through and tested. This project does not center itself around a program, but rather the learning experience of using new technologies. A thorough project-report is featured in the Wiki. <br>
The technologies used for the program, have gone smoothly. We have:
* Successfully deployed a working instance of the program on a PGFS-node within the NTNU OpenStack-network. Communication between local PGFS-nodes and the cloud-based server-PGFS-node is working perfectly.
* Completely dockerized the "proof-of-concept"-program by utilizing docker-compose with 3 docker-containers (1 server-node + 2 client-nodes)
* Learned a lot about IPFS, and its go-implementation.
* -- more
### What didn't go as expected
Getting to know IPFS on a deep level was both challenging and rewarding, and it was way more time-consuming than first estimated. The application serves as a "proof of concept" at this stage. The time usage within the proposed timeframe of the project mainly consisted of the group familiarizing themselves with IPFS and its protocols. This led to the development itself to slow down, as developing a program with IPFS requires a deep-level understanding of how it actually works. 
* Some platform issues, we are a diverse group in terms of preferred OS. 1x Linux, 1x Windows, 1x macOS.
### The most challenging aspects
* Here goes aspects of the project the group deemed most challenging.
* Limited Resources - Little Info for IPFS due to the "newness" of the whole concept.
* Platform-limiting factors, especially for Windows. NO info could be found regarding IPFS-development for Windows. All external development (used for help in most instances) is primarily executed in Linux.
### The learning outcome
* What did we learn? IPFS being the obvious answer, but it may be worth adding stuff like interrogating a foreign repository to make it fit specialized needs.
### Total work hours
For this project, a high focus on collaboration has been prioritized. For every week of the project period, the group has had 6-10 hours of collaborative work (usually 2-3 sessions per week). The sessions have mostly happened physically (on campus) and each session usually has had a specific purpose (some sprints, some commune problem solving). In addition to these teamwork-sessions, we also split some work to work on an individual basis, to increase the learning-outcome per individual group member. ----put some more later, but ground work has been laid ;)



<br>Authored by<br>
<b>Yoav Weber</b><br>
<b>Milosz Antoni Wudarczyk</b><br>
<b>Kristian Amundsen Øhman-Norén</b><br>
2021, Norwegian University of Science and Technology
