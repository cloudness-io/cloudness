# Cloudness

An open-source & self-hostable alternative to Heroku / Netlify / Vercel for Kubernetes.

<!-- [![License](https://img.shields.io/github/license/cloudness-io/cloudness)](LICENSE)
[![Build Status](https://github.com/cloudness-io/cloudness/actions/workflows/pr-validation.yml/badge.svg)](https://github.com/cloudness-io/cloudness/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudness-io/cloudness)](https://goreportcard.com/report/github.com/cloudness-io/cloudness) -->

## 📖 About the Project

Cloudness is an open-source & self-hostable platform for deploying and managing applications on Kubernetes.

It helps you manage your servers, applications, and databases on your own infrastructure; you only need a Kubernetes cluster. You can manage any Kubernetes cluster - cloud providers, on-premises, Raspberry PIs, and anything else.

Imagine having the ease of a cloud platform like Heroku, but with your own infrastructure. That is Cloudness.

**No vendor lock-in** - all configurations for your applications, databases, and services are stored as Kubernetes manifests. If you decide to stop using Cloudness, you can still manage your running resources. You lose the automations and all the magic. 🪄

## 📥 Installation

```bash
curl -fsSL https://get.cloudness.io/install.sh | bash
```

> **Note:** Please refer to the [documentation](https://docs.cloudness.io) for more information about installation and configuration.

## ✨ Features

- 🚀 **Easy Deployment** - Deploy applications to Kubernetes with minimal configuration
- 🔄 **CI/CD Integration** - Built-in pipeline support for automated builds and deployments  
- 📦 **Template System** - Pre-configured templates for common services (PostgreSQL, Redis, MySQL, Valkey)
- 🔐 **Authentication & Authorization** - Secure access control with multi-tenant support
- 📊 **Project Management** - Organize applications, environments, and deployments
- 📝 **Real-time Logs** - Stream application logs in real-time
- 🎯 **Multi-tenant** - Support for multiple organizations and projects

## 💬 Support

- 📖 [Documentation](https://docs.cloudness.io)
- 💬 [GitHub Discussions](https://github.com/cloudness-io/cloudness/discussions)
- 🐛 [Issue Tracker](https://github.com/cloudness-io/cloudness/issues)

## 🛠️ Development

### Pre-Requisites

Install the latest stable version of Node and Go version 1.21 or higher. Ensure the GOPATH [bin directory](https://go.dev/doc/gopath_code#GOPATH) is added to your PATH.

### Clone the repository

```bash
git clone https://github.com/cloudness-io/cloudness.git
cd cloudness
```

### Install required Go tools

```bash
make dep
make tools
```

### Build

Build the Cloudness binary:

```bash
make build
```

### Run

This project supports all operating systems and architectures supported by Go. This means you can build and run the system on your machine; docker containers are not required for local development and testing.

To start the server at localhost:8000, simply run the following command:

```bash
./cloudness server .local.env
```

The application will start at `http://localhost:8000`. The database schemas will be auto-migrated on startup.

## 💻 CLI

This project includes command line tools for development and running the service. For a full list of supported operations, please see:

```bash
./cloudness --help
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CLOUDNESS_DATABASE_DRIVER` | Database driver (postgres/sqlite) | `postgres` |
| `CLOUDNESS_DATABASE_HOST` | Database host | `localhost` |
| `CLOUDNESS_DATABASE_PORT` | Database port | `5432` |
| `CLOUDNESS_DATABASE_NAME` | Database name | `cloudness` |
| `CLOUDNESS_DATABASE_USER` | Database username | - |
| `CLOUDNESS_DATABASE_PASSWORD` | Database password | - |
| `CLOUDNESS_PUBSUB_PROVIDER` | Pub/Sub provider (redis/inmem) | `inmem` |
| `CLOUDNESS_REDIS_ENDPOINT` | Redis endpoint (if using) | - |
| `CLOUDNESS_REDIS_PASSWORD` | Redis password | - |
| `CLOUDNESS_DEBUG` | Enable debug logging | `false` |
| `CLOUDNESS_TRACE` | Enable trace logging | `false` |

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## 🏗️ Architecture

Cloudness is built with:

- **Backend:** Go with Gin framework
- **Frontend:** HTML/Templ/Alpine.js/Tailwindcss
- **Database:** PostgreSQL (primary), MySQL (supported)
- **Cache/Pub-Sub:** Redis or in-memory
- **Orchestration:** Kubernetes

### Project Structure

```
cloudness/
├── app/              # Application core
│   ├── auth/         # Authentication logic
│   ├── controller/   # HTTP controllers
│   ├── middleware/   # HTTP middleware
│   ├── router/       # Route definitions
│   ├── services/     # Business logic services
│   ├── store/        # Application data stores
│   └── web/          # HTML Templ frontend
├── blob/             # Blob storage interface
├── cli/              # CLI implementation
├── cmd/              # Application entrypoint
├── errors/           # Error types and handling
├── helpers/          # Utility functions
├── http/             # HTTP client utilities
├── job/              # Background job scheduler
├── k8s/              # Kubernetes manifests
├── lock/             # Distributed locking
├── logging/          # Logging configuration
├── logstream/        # Real-time log streaming
├── plugins/          # Builder and deployer plugins
├── profiler/         # Performance profiling
├── pubsub/           # Pub/Sub implementation
├── schema/           # JSON schemas
├── scripts/          # Installation and ops scripts
├── store/            # Database layer
├── templates/        # Application templates
├── types/            # Type definitions
└── version/          # Version information
```
<!-- 
## Donations

To stay completely free and open-source, with no features behind a paywall, we need your help. If you like Cloudness, please consider donating to support the project's development.

[Become a sponsor](https://github.com/sponsors/cloudness-io)

Thank you so much! -->

## 🗺️ Roadmap

- [ ] Enhanced monitoring and observability
- [ ] Multi-cloud support
- [ ] Advanced deployment strategies (canary, blue-green)
- [ ] Marketplace for community templates
- [ ] GitOps integration
- [ ] Cost optimization features

## 📄 License

This project is licensed under the Apache License 2.0, see [LICENSE](LICENSE).

## 🙏 Acknowledgments

Built using:

- [Go](https://golang.org/)
- [Kubernetes](https://kubernetes.io/)
- [Templ](https://templ.guide/)
- [Alpine.js](https://alpinejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)

---

Made with ❤️ by the Cloudness team
