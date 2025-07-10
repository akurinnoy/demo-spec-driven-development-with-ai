_**Instructions:** This is a template for assigning tasks to an AI assistant. Fill out the sections below to describe your project. The more detailed and clear your description, the better the result will be. **Please delete this italicized text and other hints before saving the final version of the file.**_

---

## 🎯 Project Goal

_**Hint:** Describe the main goal of your project in one or two sentences. What is the final product you want to achieve?_

**Example:**
Create a simple web application for managing a To-Do list, with the ability to save data to a local file.

---

## ⚙️ General Requirements & Setup

_**Hint:** Describe the general tools, settings, and structure for the project here._

-   **Programming Language & Environment:** [e.g., JavaScript and Node.js, Python 3.10]
-   **Package Manager:** [e.g., yarn, npm, pip]
-   **Project Structure:** [e.g., A monorepo with `yarn workspaces` for `frontend` and `backend`, or a single simple directory]
-   **Code Quality Tools:** [e.g., Set up ESLint with the `airbnb` configuration. Add Prettier for code formatting.]
-   **Version Control:**
    -   Create a `.gitignore` file.
    -   Add standard directories and files to it (e.g., `node_modules`, `.env`, build outputs).
    -   **Important:** Add AI assistant-specific files: `.ra-aid/` and `.aider.*`.
-   **Documentation:** Create a root `README.md` file that includes:
    -   A short project description.
    -   Instructions for installing dependencies.
    -   Instructions for running the project and tests.

---

## 📦 Component 1: [Name, e.g., Backend]

_**Hint:** Break your project into logical parts (components). For each component, describe its functionality in detail. If you have a simple project, this might be the only component section._

-   **Technology Stack:** [e.g., Express.js, FastAPI, Flask]
-   **Core Functionality:**
    -   [Describe the first feature, e.g., "`GET /tasks` API endpoint should return a list of all tasks"]
    -   [Describe the second feature, e.g., "`POST /tasks` API endpoint should accept a task title, create it, and save it"]
-   **Data Requirements (if any):** [e.g., "Data should be stored in a `database.json` file as an array of objects. Each object should have `id`, `title`, and `isDone` fields."]
-   **Testing:**
    -   **Frameworks:** [e.g., Jest, PyTest]
    -   **What to test:** [e.g., "Write unit tests for the task creation logic. Write integration tests for all API endpoints."]

---

## 📦 Component 2: [Name, e.g., Frontend]

_**Hint:** Copy this block if your project has another logical component, like a client-side application._

-   **Technology Stack:** [e.g., React, Vue, or plain HTML/CSS/JS]
-   **Core Functionality:**
    -   [e.g., "The main screen should display a list of tasks fetched from the backend."]
    -   [e.g., "The user can enter text for a new task in an input field."]
    -   [e.g., "Clicking the 'Add' button sends the task to the backend."]
    -   [e.g., "The user can mark a task as completed."]
-   **Environment Specifics (if any):** [e.g., "Since the project will run in Eclipse Che, the API request URL should not be hardcoded (`localhost`) but should be built dynamically based on `window.location`."]
-   **Testing:**
    -   **Frameworks:** [e.g., Jest and React Testing Library]
    -   **What to test:** [e.g., "Verify that the task list component renders data correctly. Simulate adding a new task and verify that it appears in the list."]

---

## 🏁 Final Deliverable

_**Hint:** Describe what you consider a completed project. This helps the AI know when to stop._

**Example:**
The project is considered complete when:
1.  Both parts (`frontend` and `backend`) run without errors.
2.  A new task can be added via the UI and is successfully saved.
3.  All described tests pass successfully.
4.  The `README.md` file is filled out and contains up-to-date instructions.
