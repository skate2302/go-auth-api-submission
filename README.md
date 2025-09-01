# Go User Authentication API

This is a simple and reliable user authentication service built with Go, Gin, and MongoDB. 

### Features

-   **Core:** User Sign-Up (`/signup`) and User Login (`/login`).
-   **Authentication:** Secure session management using JSON Web Tokens (JWT).
-   **Database:** MongoDB for flexible and scalable data storage.
-   **Security:** Password hashing using `bcrypt` and abuse prevention via IP-based rate limiting.
-   **Containerization:** Full Docker and Docker Compose setup for one-command environment startup.
-   **API Documentation:** Interactive API documentation powered by Swagger/OpenAPI.
-   **Code Quality:** Built with a clean, scalable, and professional project structure.

---

### Technical Stack

-   **Language:** Go
-   **Web Framework:** Gin
-   **Database:** MongoDB (via official Go driver)
-   **Containerization:** Docker & Docker Compose
-   **Security:** `bcrypt` for hashing, `golang.org/x/time/rate` for rate limiting
-   **API Documentation:** `swaggo/swag`

---

### Prerequisites

-   [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

---

### 🚀 Getting Started (One-Command Setup)

This project is fully containerized, making setup incredibly simple. You **do not** need to install Go on your machine.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/skate2302/go-auth-api-submission.git>
    cd go-auth-api
    ```

2.  **Create an environment file:**
    Copy the example `.env.example` file to a new file named `.env`. The default values are already configured for the Docker setup.
    ```bash
    cp .env.example .env
    ```
    *(Note: You can create a `.env.example` file and add your `PORT`, `MONGO_URI`, and `JWT_SECRET` variables there for others to use as a template).*

3.  **Build and run the services:**
    This single command will build the Go application's Docker image and start both the API and database containers.
    ```bash
    docker-compose up --build -d
    ```
    -   `--build`: Forces Docker to build the Go application image from the `Dockerfile`.
    -   `-d`: Runs the containers in the background (detached mode).

The API will now be running and available at `http://localhost:8080`.

---

### API Usage

Once the application is running, you can interact with it in two ways:

#### 1. Interactive API Documentation (Swagger)

The easiest way to explore and test the API is through the built-in Swagger UI.

-   **Open your browser and navigate to:** [**http://localhost:8080/swagger/index.html**](http://localhost:8080/swagger/index.html)

Here you can view all endpoints, see their request/response formats, and execute them directly from the browser.

#### 2. Using Postman or `curl`

-   **Sign Up:** `POST /api/v1/signup`
-   **Login:** `POST /api/v1/login`

---

### Stopping the Application

To stop all running containers, navigate to the project directory and run:
```bash
docker-compose down
```
