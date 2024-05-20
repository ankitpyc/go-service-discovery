# 🌟 Go Service Discovery 🌟

![Go](https://img.shields.io/badge/Go-1.16-blue.svg)
![License](https://img.shields.io/github/license/ankitpyc/go-service-discovery)
![Build Status](https://img.shields.io/github/actions/workflow/status/ankitpyc/go-service-discovery/go.yml?branch=main)

🔍 A lightweight service discovery mechanism written in Go to manage and monitor distributed systems efficiently.

![Service Discovery](https://user-images.githubusercontent.com/your-image.png)

## 🚀 Features

- 🌐 **Cluster Management**: Easily manage and configure clusters.
- 📡 **Service Discovery**: Real-time service discovery and health checks.
- 🔄 **Fault Tolerance**: Automatic detection and handling of node failures.
- 📊 **Scalability**: Designed to scale horizontally.

## 📖 Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API](#api)
- [Contributing](#contributing)
- [License](#license)

## 🛠️ Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/ankitpyc/go-service-discovery.git
    cd go-service-discovery
    ```

2. **Build the project:**

    ```bash
    go build
    ```

## 🚦 Usage

### Starting the Server

To start the TCP server, use the following command:

```bash
./go-service-discovery --host <HOST> --port <PORT>

```

## Example:

 ```bash
./go-service-discovery --host 127.0.0.1 --port 2212

```

 ```bash

// Sending a Request
You can use netcat to send a request to the server:
echo -n -e '\x00{"ClusterID": "1", "NodeID": "node-123"}' | nc 127.0.0.1 2212
```
📚 API

## Features

- **AddClusterMemberList(member ClusterMember):** Add a new member to the cluster.
- **ListenForBroadcasts():** Listen for broadcast messages within the cluster.
- **Server Methods:**
  - **StartServer() (*Server, error):** Start the TCP server.
  - **ListenAndAccept() error:** Accept incoming TCP connections.
  - **StopServer() error:** Stop the TCP server.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](link_to_your_issues_page).

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature/feature-name`
5. Open a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](link_to_your_license_file) file for details.

## ✨ Credits

**Author:** Ankit

**Contributors:** See the list of contributors who participated in this project.
