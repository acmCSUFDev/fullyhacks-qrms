# QR Management System for Fullyhacks

This is the QR Management System for Fullyhacks. It is a service used to render
QR codes for checking in and out of events, ensuring that each user is only
checked in once. It also helps to keep track of the number of people who have
checked in and out of the event.

## Hosting on Production

fullyhacks-qrms is currently hosted via [acm-aws](https://github.com/acmcsufoss/acm-aws).
It is hosted on an on-campus server.

## Development

You should have Nix installed to develop this project. If you don't have Nix
installed, see [Nix's installation guide](https://nixos.org/download.html).

To start developing, run the following commands:

```sh
nix develop
```

Then, re-generate all files and run the server:

```sh
go generate ./...
go run .
```
