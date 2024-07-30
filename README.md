# USD to BRL Quotation Service

This project consists of a server that fetches the USD to BRL exchange rate from an external API and a client that retrieves the quotation from the server and saves the `bid` value to a file.

## Server

The server provides an endpoint to get the current USD to BRL exchange rate. It fetches the data from the [AwesomeAPI](https://docs.awesomeapi.com.br/api-de-moedas) and saves the quotation to a SQLite database.

### How to Run the Server

1. Ensure you have Go and SQLite3 installed.
2. Save the server code to a file named `server.go`.
3. Run the server:

    ```sh
    go run server.go
    ```

### Endpoint

- **GET /cotacao**
  - Description: Fetches the current USD to BRL exchange rate.
  - Response: JSON object containing the exchange rate data.
  - Example response:

    ```json
    {
      "USDBRL": {
        "code": "USD",
        "codein": "BRL",
        "name": "Dólar Americano/Real Brasileiro",
        "high": "5.6136",
        "low": "5.6125",
        "varBid": "0.0015",
        "pctChange": "0.03",
        "bid": "5.6131",
        "ask": "5.6141",
        "timestamp": "1722373206",
        "create_date": "2024-07-30 18:00:06"
      }
    }
    ```

## Client

The client sends a request to the server's `/cotacao` endpoint to get the USD to BRL exchange rate. If the request is successful, it extracts the `bid` value from the response and saves it to a file named `cotacao.txt`.

### How to Run the Client

1. Ensure the server is running.
2. Save the client code to a file named `client.go`.
3. Run the client:

    ```sh
    go run client.go
    ```

### Output

The client creates a file named `cotacao.txt` with the following content format:
