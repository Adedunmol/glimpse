# Go Glimpse

Go Glimpse is a monorepo containing the backend services, shared packages, and development tooling required to run the application locally.

The mobile application is maintained in a separate repository:

- Mobile App: [https://github.com/ufuomaisaac/Glimpse]

## Prerequisites

Before getting started, ensure the following tools are installed:

* Bun
* Docker & Docker Compose (required by backend services)
* Go (for backend development)

## Installing Dependencies

Install all JavaScript/TypeScript dependencies for the workspace:

```bash
bun install
```

This command installs dependencies for all packages and applications managed within the monorepo.

## Running the Development Environment

To start all development services, run:

```bash
bun dev
```

This command starts the configured development processes across the workspace, including any frontend applications, backend services, and supporting tooling defined in the project.

## Development Workflow

1. Clone the repository.

2. Install dependencies:

   ```bash
   bun install
   ```

3. Start the development environment:

   ```bash
   bun dev
   ```

4. Make your changes.

5. Run any relevant tests or code quality checks before committing.

## Repository Structure

```text
.
├── backend/        # Go backend application and infrastructure
├── packages/       # Shared packages and tooling
└── package.json    # Workspace configuration and scripts
```

## Useful Commands

### Install Dependencies

```bash
bun install
```

### Start Development Services

```bash
bun dev
```

### Update OpenAPI Documentation

If backend API contracts have changed:

```bash
cd packages/openapi
bun gen
```

This regenerates the OpenAPI artifacts and keeps the API documentation synchronized with the backend implementation.
