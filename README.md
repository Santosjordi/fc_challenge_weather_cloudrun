# FC Weather App

## How to Run it

To run your Go weather application, you'll need to set up the environment, configure API keys, and either run the Go binary directly or use Docker. This guide covers both methods.

### ⚙️ Prerequisites

  * **Go:** You need Go version 1.22 or higher installed.
  * **API Keys:** Obtain a free API key from [WeatherAPI](https://www.weatherapi.com/).

-----

### 🚀 Running the Application

#### **1. Environment Setup**

Create a file named `.env` in the root directory of your project with the following content. Replace `YOUR_WEATHERAPI_KEY` with the key you obtained.

```ini
WEATHER_API_KEY=YOUR_WEATHERAPI_KEY
SERVER_PORT=:8080
```

#### **2. Running with Go**

This is the quickest way to get started.

1.  Navigate to the `cmd` directory of your project.

2.  Run the following command to start the server:

    ```sh
    go run main.go
    ```

3.  The application will start, and you'll see a message like "Server is running on port :8080".

#### **3. Running with Docker**

This method packages your application into a container, which is ideal for deployment.

1.  Ensure you have **Docker** installed and running on your system.

2.  From the project's root directory, build the Docker image using the `Dockerfile` you created:

    ```sh
    docker build -t fc-weather-app .
    ```

3.  Run the container on your machine. Pass the API key and port as an environment variable to the container.

    ```sh
    docker compose up --build
    ```

-----

### 💻 API Usage

Once the server is running, you can test the endpoint using a tool like `curl` or a web browser.

**Endpoint:** `http://localhost:8080/{cep}`

  * Replace `{cep}` with an 8-digit Brazilian ZIP code (e.g., `89068210`).

**Example Request:**

```sh
curl http://localhost:8080/89068210
```

**Example Successful Response:**

```json
{
  "temp_C": 25.5,
  "temp_F": 77.9,
  "temp_K": 298.65
}
```

**Failure Scenarios:**

  * **Invalid format:** Sending a non-8-digit CEP will return a `422 Unprocessable Entity` error.
  * **CEP not found:** A valid-format CEP that doesn't exist will return a `404 Not Found` error.

## ASSIGNMENT

### Objetivo:

Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

### Requisitos:
- [X] O sistema deve receber um CEP válido de 8 digitos
- [X] O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin.

- [X] O sistema deve responder adequadamente nos seguintes cenários:

Em caso de sucesso:

Código HTTP: 200
Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }

Em caso de falha, caso o CEP não seja válido (com formato correto):
Código HTTP: 422
Mensagem: invalid zipcode

​​​Em caso de falha, caso o CEP não seja encontrado:
Código HTTP: 404
Mensagem: can not find zipcode


- [ ] Deverá ser realizado o deploy no Google Cloud Run.


### Dicas:

- Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/
- Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: https://www.weatherapi.com/
- Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
- Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
- Sendo F = Fahrenheit
- Sendo C = Celsius
- Sendo K = Kelvin


Entrega:

- [X] O código-fonte completo da implementação.
- [X] Testes automatizados demonstrando o funcionamento.
- [] Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- [] Deploy realizado no Google Cloud Run (free tier) e endereço ativo para ser acessado.