# Educational TCP Server Project

## Overview

This project is designed as an educational tool to help learners understand the basics of TCP server operations in Go. It demonstrates how to handle TCP connections, manage concurrent client interactions using goroutines, and implement a basic message distribution system using channels and maps.

## Key Concepts Covered

- **TCP Server Setup**: Establishing a TCP server that listens on a specific port.
- **Connection Handling**: Managing incoming connections and ensuring proper closure of connections.
- **Concurrency**: Using goroutines to handle multiple clients simultaneously.
- **Message Handling**: Receiving, processing, and distributing messages based on client requests.
- **Error Handling**: Properly handling potential errors that may occur during network communication.

## How It Works

- The server listens on port 2000 for incoming TCP connections.
- Each connection is handled in a separate goroutine to allow multiple clients to interact with the server concurrently.
- Clients can join specific queues, leave queues, or send messages to queues.
- Messages sent to a queue are distributed to all members of that queue.

## Usage

- Start the server by running the `main` function.
- Connect to the server using any TCP client on `localhost:2000`.
- Send JSON formatted commands to join queues, leave queues, or send messages.

## Educational Objectives

- Understand the basics of network programming in Go.
- Learn how to use goroutines for handling concurrency.
- Explore the use of channels for inter-goroutine communication.
- Gain insights into practical error handling in network applications.

This project is intended for educational purposes and provides a basic framework that can be expanded with additional features for more advanced learning experiences.
