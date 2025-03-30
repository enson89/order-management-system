
# Order Management System

A microservices-based Order Management System built with Golang, PostgreSQL, Redis, Kafka, Docker Compose, Minikube, and Swagger. This project demonstrates a layered architecture, best practices for containerization, and automated deployments.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Folder Structure](#folder-structure)
- [Prerequisites](#prerequisites)
- [Setup and Installation](#setup-and-installation)
    - [Local Development with Docker Compose](#local-development-with-docker-compose)
    - [Kubernetes Deployment with Minikube](#kubernetes-deployment-with-minikube)
- [Makefile Usage](#makefile-usage)
- [API Documentation (Swagger)](#api-documentation-swagger)
- [Running Tests](#running-tests)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

---

## Overview

The Order Management System is a RESTful API service that allows clients to create and retrieve orders. It uses a layered architecture for clear separation of concerns and integrates with:

- **PostgreSQL** for data persistence
- **Redis** for caching
- **Kafka** for event streaming
- **Docker Compose** for local multi-container development
- **Kubernetes (Minikube)** for orchestration and deployment
- **Swagger** for API documentation

---

## Features

- **CRUD Operations:** Create and retrieve orders.
- **Database Integration:** Uses PostgreSQL for reliable data storage.
- **Caching:** Redis to cache frequently accessed data.
- **Event Streaming:** Kafka to publish order events asynchronously.
- **Containerization:** Docker Compose for local development and testing.
- **Kubernetes Deployment:** Manifests provided for deploying on Minikube.
- **API Documentation:** Automatically generated Swagger docs.
- **Automated CI/CD:** Makefile for building, testing, and deploying.

---

## Architecture

The application follows a clean, layered architecture:

- **API Layer (Handlers):** Manages HTTP requests and routes.
- **Service Layer:** Contains business logic.
- **Repository Layer:** Handles data persistence with PostgreSQL.
- **Cache Layer:** Integrates with Redis for caching.
- **Event Streaming:** Uses Kafka to send order events.
- **Configuration:** Uses a centralized YAML configuration file.

---

## Folder Structure

```
order-management/
├── api/                    # Swagger documentation files (generated)
├── cmd/
│   └── app/
│       └── main.go       # Application entry point
├── config/
│   └── config.yaml       # Application configuration file
├── deployments/          # Kubernetes manifests for deployments, services, and ingress
│   ├── app.yaml
│   ├── redis.yaml
│   ├── postgres.yaml
│   └── ingress.yaml
├── docker/
│   └── Dockerfile        # Dockerfile for building a multi-arch image
├── internal/
│   ├── cache/            # Redis integration
│   ├── config/           # Configuration loading logic
│   ├── db/               # Database models and PostgreSQL integration
│   ├── handlers/         # HTTP handlers (API endpoints)
│   ├── kafka/            # Kafka producer and consumer logic
│   ├── repository/       # Data access (CRUD operations)
│   ├── routes/           # API route definitions
│   └── service/          # Business logic and order services
├── scripts/              # Helper scripts (e.g., SQL initialization)
├── go.mod                # Go module file
├── go.sum                # Go module checksums
├── Makefile              # Automation for building, testing, and deploying
└── README.md             # Documentation
```

---

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go** (version 1.20+ recommended)
- **Docker** (with Buildx enabled for multi-architecture builds)
- **Minikube** (for local Kubernetes deployment)
- **kubectl** (Kubernetes command-line tool)
- **PostgreSQL**, **Redis**, and **Kafka** (via Docker images)
- **Swagger CLI** (for generating API documentation)

---

## Setup and Installation

### Local Development with Docker Compose

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/order-management.git
   cd order-management
   ```

2. **Configure the Application:**

   Edit `config/config.yaml` to set your PostgreSQL, Redis, and Kafka connection details.

3. **Build and Run with Docker Compose:**

   ```bash
   docker-compose up --build
   ```

   This command builds the Docker images for the application and supporting services. Your application should be accessible (typically on port 8080).

4. **Test the Application:**

   You can use cURL or your browser to test endpoints (e.g., `http://localhost:8080/health`).

---

### Kubernetes Deployment with Minikube

1. **Start Minikube:**

   ```bash
   make minikube-start
   ```

2. **Deploy the Application:**

   Deploy application, Redis, PostgreSQL, and Ingress manifests using:
   ```bash
   make deploy-minikube
   ```

3. **Open Minikube Tunnel:**

   If using a LoadBalancer or exposing services, run:
   ```bash
   sudo make minikube-tunnel
   ```

4. **Access the Application:**

   Open your browser and navigate to:

   ```
   http://order.local
   ```

---

## Makefile Usage

The project includes a `Makefile` to automate various tasks such as building, testing, and deploying the application.

### 🚀 **Common Targets**

- **Run All Checks and Execute Locally:**
```bash
make all
```

- **Build and Run Locally:**
```bash
make run
```

- **Check Code Quality:**
```bash
make lint
```

- **Check Code Format:**
```bash
make format
```

- **Check for Vulnerabilities:**
```bash
make vulncheck
```

- **Run Tests:**
```bash
make test
```

- **Build Docker Image:**
```bash
make docker-build
```

- **Push Docker Image:**
```bash
make docker-push
```

- **Start Minikube and Enable Ingress:**
```bash
make minikube-start
```

- **Deploy to Minikube:**
```bash
make deploy-minikube
```

- **Open Minikube Tunnel:**
```bash
sudo make minikube-tunnel
```

- **Delete Minikube Cluster:**
```bash
make delete-minikube
```

- **Clean Up:**
```bash
make clean
```

---

## API Documentation (Swagger)

1. **Generate Swagger Documentation:**

   Install the Swagger CLI:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

   Then run:

   ```bash
   swag init -g cmd/app/main.go
   ```

2. **Access the Swagger UI:**

   The Swagger UI is available at:

   ```
   http://order.local/swagger/index.html
   ```

---

## Running Tests

- **Unit Tests:**

  Run the tests from the root directory:

  ```bash
  go test ./internal/... -v
  ```

- **Integration Tests:**

  Integration tests can be run using Docker Compose or Minikube to spin up the environment, then executing the tests against that environment.

---

## Troubleshooting

- **Service Endpoints Not Found:**
    - Verify that your Deployment’s pod template has the label `app: order-management`.
    - Check endpoints with:
      ```bash
      kubectl get endpoints order-management
      ```

- **Ingress Not Routing Traffic:**
    - Confirm that your Ingress has an ADDRESS by running:
      ```bash
      kubectl get ingress order-management-ingress
      ```
    - Ensure your `/etc/hosts` file maps `order.local` to your Minikube IP.

- **Pod Readiness Issues:**
    - Inspect your pods:
      ```bash
      kubectl describe pod <pod-name>
      ```
    - Verify that readiness probes are passing.

- **Logs and Events:**
    - Check Ingress controller logs for errors:
      ```bash
      kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
      ```

---

## Contributing

Contributions are welcome! Please fork this repository and submit a pull request with your improvements or bug fixes. For major changes, open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For any questions or issues, please open an issue on GitHub or contact [your.email@example.com](mailto:your.email@example.com).

---

Happy coding!
