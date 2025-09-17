# Inventory Management System

This project is a Go-based Inventory Management System designed to manage devices, owners, properties, and related logs. It provides a modular structure for handling inventory data, supporting CRUD operations, and maintaining configuration and logging.

## Project Structure

- **main.go**: Entry point of the application.
- **config/**: Contains configuration files (`app.yaml`, `app_example.yaml`).
- **devices/**: Device management (DAO, handlers, services, logs, photos).
- **owners/**: Owner management (DAO, handlers, services).
- **properties/**: Property management (DAO, handlers, services).
- **types/**: Type management (DAO, handlers, services, property types).
- **utils/**: Utility functions (database connection, models).
- **inventory.sql**: SQL schema for database setup.
- **inventory.log**: Log file for application events.
- **go.mod / go.sum**: Go module dependencies.

## Features

- Device, owner, property, and type management
- Modular DAO, service, and handler layers
- YAML-based configuration
- Logging support
- SQL schema for database initialization

## Getting Started

1. **Clone the repository**
2. **Configure the application**
   - Copy `config/app_example.yaml` to `config/app.yaml` and update as needed.
3. **Set up the database**
   - Use `inventory.sql` to initialize your database.
4. **Build and run**
   ```sh
   go build -o inventory
   ./inventory
   ```

## Requirements
- Go 1.18+
- A supported SQL database (see `inventory.sql`)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the MIT License.
