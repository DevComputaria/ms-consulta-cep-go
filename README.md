# Go CEP Service

This project implements a service to query Brazilian address information based on CEP (postal code) using the ViaCEP API. It's structured following Clean Architecture principles.

## Project Structure

-   `/cmd`: Main application entry point.
-   `/domain`: Core domain entities (e.g., `Address`).
-   `/usecase`: Application-specific business logic (services).
-   `/interfaces`: Adapters to external systems.
    -   `/interfaces/services`: Clients for external services (e.g., ViaCEP client).
    -   `/interfaces/http`: HTTP handlers for exposing the API.

## Prerequisites

-   Go (version 1.16 or higher)

## How to Run

1.  **Clone the repository.**
2.  **Navigate to the project directory.**
3.  **Build the application:**
    ```bash
    go build -o cep-service ./cmd/main.go
    ```
4.  **Run the application:**
    ```bash
    ./cep-service
    ```
    The server will start on port 8080.

## API Endpoint

### Get Address by CEP

-   **URL:** `/cep/{cep_value}`
-   **Method:** `GET`
-   **Description:** Retrieves address information for the given CEP.
-   **Example:**
    ```
    GET /cep/01001000
    ```
-   **Success Response (200 OK):**
    ```json
    {
        "cep": "01001-000",
        "logradouro": "Praça da Sé",
        "complemento": "lado ímpar",
        "bairro": "Sé",
        "localidade": "São Paulo",
        "uf": "SP",
        "ibge": "3550308",
        "gia": "1004",
        "ddd": "11",
        "siafi": "7107"
    }
    ```
-   **Error Responses:**
    -   `400 Bad Request`: If the CEP is missing or invalid in the request.
        ```json
        { "error": "CEP must be provided in the URL path, e.g., /cep/01001000" }
        ```
    -   `404 Not Found`: If the CEP is not found.
        ```json
        { "error": "Address not found for CEP: <cep_value>" }
        ```
    -   `500 Internal Server Error`: For other server-side errors.
        ```json
        { "error": "Internal server error" }
        ```

## How to Run Tests

Navigate to the project directory and run:
```bash
go test ./...
```

This will execute all tests in the project.
