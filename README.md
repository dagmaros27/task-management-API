# Task Management REST API Documentation with Clean Architecture

## Overview

This version of the Task Management REST API has been refactored using Clean Architecture principles. The codebase is organized into distinct layers with clear separation of concerns, enhancing maintainability, testability, and scalability. The API uses MongoDB for persistent data storage, and JWT tokens for authentication and authorization, with role-based access control.

## Folder Structure

The codebase is organized as follows:

```
task-manager/
├── Delivery/
│   ├── main.go
│   ├── controllers/
│   │   └── controller.go
│   └── routers/
│       └── router.go
├── Domain/
│   └── domain.go
├── Infrastructure/
│   ├── auth_middleWare.go
│   ├── jwt_service.go
│   └── password_service.go
├── Repositories/
│   ├── task_repository.go
│   └── user_repository.go
└── Usecases/
  ├── task_usecases.go
  └── user_usecases.go
└── env.go
```

## Layer Descriptions

**Delivery**: Handles incoming requests and responses, sets up the HTTP server, and defines the routing configuration.

- `main.go`: Initializes dependencies and starts the server.
- `controllers/controller.go`: Handles HTTP requests and interacts with the use cases layer.
- `routers/router.go`: Configures and initializes the routes using the Gin framework.

**Domain**: Contains core business entities and logic, decoupled from external frameworks.

- `domain.go`: Defines the Task and User structs representing core entities.

**Infrastructure**: Implements external services and dependencies.

- `auth_middleWare.go`: Middleware for handling JWT-based authentication and authorization.
- `jwt_service.go`: Functions to generate and validate JWT tokens.
- `password_service.go`: Functions for securely hashing and comparing passwords.

**Repositories**: Abstracts data access logic using interfaces.

- `task_repository.go`: Interface and implementation for task-related data operations.
- `user_repository.go`: Interface and implementation for user-related data operations.

**Usecases**: Encapsulates the application's business logic.

- `task_usecases.go`: Implements use cases for creating, updating, retrieving, and deleting tasks.
- `user_usecases.go`: Implements use cases for user registration, login, and promotion to admin.

## MongoDB Integration

### Configuration

- Connection String: Configured in the `.env` file.
- Database: `task_manager`
- Collections: `tasks` and `users`

### MongoDB Installation

- Install MongoDB locally or use a cloud provider like MongoDB Atlas.
- Install the MongoDB Go Driver:

```
go get go.mongodb.org/mongo-driver/mongo
```

## Running the API

1. Clone the repository:

```
git clone [https://github.com/dagmaros27/task-management-API.git](https://github.com/dagmaros27/task-management-API.git)
```

2. Navigate to the project directory:

```
cd task-managment-API
```

3. Install dependencies:

```
go mod download
```

4. Start the API server:

```
go run delivery/main.go
```

The server will be running at `http://localhost:8080`.

## API Endpoints

### User Management

#### Register a New User

- Endpoint: `POST /register`
- Description: Creates a new user account.
- Request Body:

```json
{
  "username": "your_username",
  "password": "your_password"
}
```

- Responses:
  - `201 Created`: Successful registration.
  - `400 Bad Request`: Invalid input or username already exists.

#### Login

- Endpoint: `POST /login`
- Description: Authenticates the user and generates a JWT token.
- Request Body:

```json
{
  "username": "your_username",
  "password": "your_password"
}
```

- Responses:
  - `200 OK`: Successful login, returns a JWT token.
  - `401 Unauthorized`: Invalid username or password.

#### Promote User to Admin (Admin Only)

- Endpoint: `POST /promote`
- Description: Promotes a user to an admin role.
- Headers: `Authorization: Bearer <JWT token>`
- Request Body:

```json
{
  "username": "user_to_promote"
}
```

- Responses:
  - `200 OK`: Successful promotion.
  - `403 Forbidden`: Unauthorized access.

### Task Management

#### Create a Task (Admin Only)

- Endpoint: `POST /tasks`
- Description: Creates a new task.
- Headers: `Authorization: Bearer <JWT token>`
- Request Body:

```json
{
  "title": "Task title",
  "description": "Task description"
}
```

- Responses:
  - `201 Created`: Task created successfully.
  - `403 Forbidden`: Unauthorized access.

#### Update a Task (Admin Only)

- Endpoint: `PUT /tasks/:id`
- Description: Updates an existing task.
- Headers: `Authorization: Bearer <JWT token>`
- Request Body:

```json
{
  "title": "Updated task title",
  "description": "Updated task description"
}
```

- Responses:
  - `200 OK`: Task updated successfully.
  - `403 Forbidden`: Unauthorized access.

#### Delete a Task (Admin Only)

- Endpoint: `DELETE /tasks/:id`
- Description: Deletes an existing task.
- Headers: `Authorization: Bearer <JWT token>`
- Responses:
  - `200 OK`: Task deleted successfully.
  - `403 Forbidden`: Unauthorized access.

#### Retrieve All Tasks

- Endpoint: `GET /tasks`
- Description: Retrieves a list of all tasks.
- Headers: `Authorization: Bearer <JWT token>`
- Responses:
  - `200 OK`: Returns the task list.

#### Retrieve a Task by ID

- Endpoint: `GET /tasks/:id`
- Description: Retrieves a task by its ID.
- Headers: `Authorization: Bearer <JWT token>`
- Responses:
  - `200 OK`: Returns task details.
  - `404 Not Found`: Task not found.

## Authentication & Authorization

- JWT Token: After a successful login, the server generates a JWT token, which must be included in the Authorization header for protected routes.
- Format: `Authorization: Bearer <JWT token>`
- User Roles:
  - Admin: Full access to all endpoints.
  - Regular User: Can only retrieve tasks.
- Middleware:
  - Authentication: Validates JWT tokens before granting access.
  - Authorization: Checks user roles for admin-specific routes.

## Security

- Password Storage: Passwords are hashed using bcrypt.
- Token Security: JWT tokens are signed with a secret key.

## Testing

Use Postman or curl to test the API endpoints. Ensure to test both authenticated and unauthenticated scenarios.

Example:

- Register a New User:

```
curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{"username": "newuser", "password": "password123"}'
```

- Login:

```
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"username": "newuser", "password": "password123"}'
```

- Create a Task (Admin Only):

```
curl -X POST http://localhost:8080/tasks \
-H "Authorization: Bearer <JWT token>" \
-H "Content-Type: application/json" \
-d '{"title": "New Task", "description": "Task description"}'
```

## Environment Variables

The API uses a `.env` file for configuration. Ensure the following environment variables should be set in your `.env` file:

- `APP_ENV`: The environment in which the app is running (e.g., development, production).
- `DB_URI`: The URI for connecting to MongoDB.
- `DB_NAME`: The name of the MongoDB database.
- `DB_TASK_COLLECTION`: The collection name for tasks.
- `DB_USER_COLLECTION`: The collection name for users.
- `ACCESS_TOKEN_SECRET`: The secret key used for signing JWT tokens.

## Loading Environment Variables

The environment variables are loaded using the Viper library in the `main.go` file. Ensure that you created `.env` file in the root directory of the project.

## Postman Documentation

The complete postman documentation can be found [here](https://documenter.getpostman.com/view/25928149/2sA3s3HB55)
