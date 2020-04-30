# gonotify

A web service that lets you send yourself WhatsApp messages serving as notifications

> **Demo:** [https://gonotify.xyz](https://gonotify.xyz)

The basic idea is to provide the user with a service that he/she can use to send himself WhatsApp messages programmatically without much hassle. With GoNotify you can create groups of various numbers and send messages to the group with just a simple API call. I have explained the idea in detail in [a dev.to post](https://dev.to/prmsrswt/whatsapp-messages-as-a-service-3kc).

## Installation

- You would need `go` to build the binary.
- The WebUI is built using ReactJS, and we use `yarn` as package manager. So you would need `node` and `yarn` too.
- Run `make build` in project root. This will build a binary at `build/gonotify`

## Configuration

- `gonotify` requires a configuration file to start. You can use the `-c` flag to pass a path for the config file.
- By default `gonotify` uses `config/config.yml`.
- A sample config file is provided in `config/config.example.yml`. Edit the file according to your needs.

## Features

- Add and verify your Phone numbers once.
- Create multiple groups with your phone numbers.
- Send notification to all numbers in a group.
- Use the API to do all of this programmatically.

## Usage

### Concepts

- **Number:** You can add multiple Phone numbers to your account. You need to verify each phone number once.
- **Group:** A group is a collection of phone numbers. The notifications are always targeted towards a Group.

### Example API usage

The examples use [HTTPie](https://httpie.org) as client.

GoNotify server issues JWT token on successfull authentication. `Authorization` header is used to supply the token to the backend. The `Authorization` header's value must be in the form of `Bearer <jwt token>`.

#### Login

```bash
$ http :8080/api/v1/login phone="9912312345" password="password"
```

Response:

```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUjhggdgnR5cCI6IkpXVCJ9.eyJleHjkgJFgfhvj3ODcsImlkIjoyfQ.TsxdOsxf0cOt5cNNSgOx5CH4oxxtGogPKcA0XPPyhnTaKhc4xpmcsJV_GY56bkghfhgdh0jO1TtSolOw8GT3TGtQyA"
}
```

#### Send a message

```bash
$ http :8080/api/v1/send authorization:"Bearer eyJhbGciOiJIUjhggdgnR5cCI6IkpXVCJ9.eyJleHjkgJFgfhvj3OgggttmlkIjoyfQ.TsxdOsxf0cOt5cNNSgOx5CH4oxxtGogPKcA0XPPyhnTaKhc4xpmcsJV_GY56bkghfhgdh0jO1TtSolOw8GT3TjghyyA" body="test message" group="default"
```

Response:

```json
{
  "message": "Message sent successfully"
}
```

## Roadmap

- A CLI to easiliy access the API. (In Progress)
- Expand to services other than WhatsApp (SMS, Telegram, Slack, Email etc.).

## Contributing

Any kind of contributions are welcome!
