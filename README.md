# Basic Streamer Donation using Go
Simulates streaming platforms' (ex. Twitch or YouTube's superchat) donation system using Go. The application contains two files, **client.go** and **server.go**.
* **server.go:** Runs the server, listening for requests from clients (viewers) using TCP and UDP, and hosts a local site for the streamer to start streaming and waiting for donations sent by the viewers using Websocket.
* **client.go:** Runs the viewer client, dials the server after entering the viewer's username, and sends donation to a currently active streamer. The viewer needs to have enough on their balance before sending a donation.

# Client
* The client connects to the TCP server after entering their username. It sends username information, then starts a selection menu where the client acts as a viewer sending their donations.
* **Check Balance:** The client connects to the UDP server and requests for the user balance. The client receives a response from the server, then the client continues as usual.
* **Top-Up Balance:** The client inputs the amount to top-up their balance. The input has to be a positive number to proceed. If the input is a valid number, the client connects to the UDP server and awaits for a response.
* **Send Donation:** The client inputs the streamer to donate to and the amount to donate. The amount has to be a positive number to proceed. The client then sends a TCP request then awaits a response.

# Showcase
YouTube link: https://youtu.be/NYC5OFg_7c4
