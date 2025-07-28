# Project: Che URL Shortener

A full-stack URL shortener service that creates memorable, human-readable links (e.g., `adjective-noun`) and is specifically designed to run within an Eclipse Che workspace.

-----

## 1. App Summary

A full-stack URL shortener service that creates memorable, human-readable links (e.g., `adjective-noun`) and is specifically designed to run within an Eclipse Che workspace.

-----

## 2. User Personas

  - Developer: A software developer using Eclipse Che who needs to share a temporary, easy-to-remember URL for a web application they are hosting inside their workspace.

-----

## 3. User Stories

  - US001: As a User, I can submit a long URL through a simple web interface, so that I can receive a short, memorable equivalent.
  - US002: As a User, I can share the short URL, so that my colleagues can be redirected to the original long URL.
  - US003: As a User, I can see a history of all created URLs and their usage counts, so that I can track how often they are used.
  - US004: As a User, I can run a single script to build the frontend and launch the backend server, so that setup is simple and fast.

-----

## 4. UI/UX Design Brief

  - Overall Feel: A clean, minimalist, single-page application. The design should be functional and straightforward.
  - Layout:
      - A clear title: "Che URL Shortener".
      - A large text input field for the long URL.
      - A prominent "Shorten URL" button.
      - A result area to display the generated short link or any error messages.
      - A simple table below the result area to display the history of all shortened URLs. The table should have columns for "Short URL", "Original URL", and "Clicks".

-----

## 5. Functional Requirements

### Backend (Go)

  - FR001: Word Generation: The server must, upon startup, initialize two hard-coded, in-memory string slices: one with 50 unique adjectives and one with 50 unique nouns.
  - FR002: Data Persistence: The system must use a JSON file named `urls.json` located in the `/backend` directory. If the file does not exist, it should be created with an empty JSON array `[]`.
  - FR003: Concurrency Control: All read and write operations on `urls.json` must be synchronized using a `sync.Mutex` to prevent race conditions.
  - FR004: API Endpoint `POST /api/urls`: Accepts `{ "url": "..." }`, validates the URL, generates a unique `short_code`, saves the record, and returns a `201 Created` status with `{ "short_code": "..." }`.
  - FR005: API Endpoint `GET /api/urls`: Returns a JSON array of all URL records.
  - FR006: Redirect Endpoint `GET /{short_code}`: Looks up the `short_code`, increments `usage_count`, saves the data, and performs a `302 Found` redirect. Returns `404 Not Found` if the code doesn't exist.
  - FR007: Static File Server: The server's root (`/`) must serve `frontend/build/index.html` and other static assets from `frontend/build/`.
  - FR008: Error Handling: The server must return a `500 Internal Server Error` with a JSON error message if it fails to read or write `urls.json`.

### Frontend (React)

  - FR009: API Communication: All API calls must use relative paths (e.g., `fetch('/api/urls')`).
  - FR010: Result Display: The UI must display a full, clickable short link constructed using `window.location.origin` and the received `short_code`.
  - FR011: URL History Display: On page load, the UI must fetch data from `GET /api/urls` and display the list of URLs in a table.

-----

## 6. Non-Functional Requirements

  - Performance: API responses and redirects should complete in under 100ms.
  - Reliability: The backend server must not crash due to concurrent requests.
  - Security: The server must not serve directory listings.

-----

## 7. Tech Stack & Project Structure

  - Backend: Go
  - Frontend: React
  - Project Structure:
    ```
    /
    ├── backend/
    │   ├── go.mod
    │   ├── go.sum
    │   ├── main.go
    │   ├── urls.json
    │   └── README.md
    ├── frontend/
    │   ├── public/
    │   ├── src/
    │   ├── package.json
    │   └── README.md
    ├── start.sh
    ├── .gitignore
    └── README.md
    ```

-----

## 8. Data Model

  - File: `urls.json`
  - Structure: A JSON array of objects.
  - Object Schema:
      - `short_code` (string): e.g., `clever-fox`
      - `long_url` (string): e.g., `https://example.com`
      - `created_at` (string): ISO 8601 format (e.g., `2025-07-17T11:00:00Z`)
      - `usage_count` (integer)

-----

## 9. Backend Unit Tests

  - Framework: Use Go's standard `testing` and `net/http/httptest` packages.
  - Test Cases: Implement unit tests for the following scenarios:
    1.  Successful Creation: A `POST` to `/api/urls` with a valid URL returns a `201` status and a valid JSON response.
    2.  Successful Redirect: A `GET` to a valid `/{short_code}` results in a `302` redirect and increments the `usage_count` in the data store.
    3.  Invalid URL Failure: A `POST` to `/api/urls` with a malformed URL returns a `400` status.
    4.  Not Found Failure: A `GET` to a `/{short_code}` that does not exist returns a `404` status.

-----

## 10. Startup Script (`start.sh`)

  - The agent shall create a shell script named `start.sh` in the project's root directory. The script must be designed to fully automate the build and launch process. Instead of copying a predefined text, the agent should generate a script that fulfills the following functional requirements:

    1.  Strict Error Handling: The script must exit immediately if any command fails. The standard way to achieve this is by including `set -e` at the beginning.
    2.  Frontend Build: The script must compile the React frontend. This requires:
          - Changing the current directory to `/frontend`.
          - Installing all necessary Node.js dependencies.
          - Executing the `build` script defined in `package.json`.
    3.  Backend Compilation: The script must compile the Go backend into a single executable binary. This requires:
          - Changing the current directory to `/backend`.
          - Ensuring the Go module is initialized and its dependencies are tidy.
          - Compiling the application. The resulting binary must be placed in the project's root directory (e.g., `go build -o ../che-url-shortener-server .`).
    4.  Server Execution: The script must launch the compiled backend server from the project's root directory.
    5.  Non-Blocking Operation: The server must be executed as a background process. This is critical as it allows the `start.sh` script to terminate and signal completion to the calling agent. (e.g., `./che-url-shortener-server &`).
    6.  User-Friendly Output: The script should print clear, informative messages to the console indicating which stage is currently running (e.g., "Building frontend...", "Starting backend...").

-----

## 11. Documentation Requirements

This section outlines the goal-oriented requirements for the three separate `README.md` files. The agent should generate each file to meet the specifications below.

### 11.1 Root `README.md`

This file should provide a high-level overview of the entire project.

  - Title: The main heading should be the project name, `Che URL Shortener`.
  - Project Summary: Include a concise, one-to-two sentence description of the application's purpose.
  - Technology Stack: List the primary technologies used (e.g., Backend: Go, Frontend: React).
  - How to Run: This section is crucial. It must clearly explain how to start the entire application using the `start.sh` script. It should instruct the user to:
    1.  Make the script executable (e.g., `chmod +x start.sh`).
    2.  Execute the script (e.g., `./start.sh`).

### 11.2 Backend `README.md` (`/backend/README.md`)

This file should detail the specifics of the Go backend service.

  - Title: The main heading should be "Backend Service".
  - Core Technology: State that the service is written in Go.
  - API Endpoints: This section must document all available API routes. For each endpoint, provide:
      - The HTTP Method and Path (e.g., `POST /api/urls`).
      - A brief description of what the endpoint does.
      - A clear example of how to use it, preferably with a `curl` command.
      - Document the following endpoints: `POST /api/urls`, `GET /api/urls`, and `GET /{short_code}`.
  - Data Persistence: Mention that URL data is stored in a local `urls.json` file.
  - Running Tests: Provide the command to execute the unit tests from within the `/backend` directory (e.g., `go test ./...`).

### 11.3 Frontend `README.md` (`/frontend/README.md`)

This file should cover the details of the React frontend application.

  - Title: The main heading should be "Frontend Application".
  - Core Technology: State that the application is built with React.
  - Available Scripts: Describe the standard `npm` scripts available in a Create React App project:
      - `npm start`: To run the app in development mode.
      - `npm test`: To launch the test runner.
      - `npm run build`: To build the app for production.
  - API Communication: Briefly explain that the frontend communicates with the backend API using relative paths (e.g., `/api/urls`).

-----

## 12. Agent Development Workflow

This section outlines the precise, sequential process the agent must follow.

### Step 1: Initial Project Setup

1.  Create the complete project structure and empty placeholder files as defined in Section 7.
2.  Create and populate the `.gitignore` file. It must include rules to ignore:
  - The "RA.Aid" temporary directory (`.ra-aid/`).
  - The "Aider" temporary files (`.aider*`).
  - Standard Node.js files and directories (e.g., `node_modules/`, `/frontend/build/`).
  - The compiled Go binary (`che-url-shortener-server`).
  - Common operating system files (e.g., `.DS_Store`).

### Step 2: Backend Implementation & Verification

1.  Implement Code: Write the Go application code in `/backend/` to meet all backend requirements (FR001-FR008).
2.  Implement Tests: Write the unit tests in `/backend/` as specified in Section 9.
3.  Verification Loop:
      - A. Attempt to Run: From the `/backend` directory, execute `go run .`.
      - B. Check for Success: The server is successfully running if the process starts without an immediate error exit.
      - C. Debug on Failure: If the command fails, analyze the error message, modify the code to fix the issue, and return to step 3.A. Continue this loop until the server starts successfully.

### Step 3: Frontend Implementation

1.  Write the React application code in `/frontend/` to meet all frontend requirements (FR009-FR011).

### Step 4: Final System Verification ✅

1.  Run Backend Tests: In the `/backend` directory, run the command `go test ./...`. Ensure that all tests pass. If any test fails, debug the Go code until they all pass.
2.  Build and Run Full Stack: From the project root, create and execute the `start.sh` script according to the requirements in Section 10. Ensure it completes without any build or runtime errors.
  - If necessary, create `stop.sh`. Ensure it completes without any build or runtime errors.
3.  Generate Documentation: Create the three required `README.md` files (in the root, `/backend`, and `/frontend` directories), ensuring each one meets all the requirements specified in Section 11.
4.  Final Review: Perform a final check of the generated code against all functional and non-functional requirements to ensure all conditions have been met.
