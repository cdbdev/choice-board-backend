# Choice-board-backend

This repo contains the source code for the **Choice Board** BACKEND application.

It consists of 3 folders:
- internal (source code of the application)
- db (database used by application)
- web (javascript bundle and static pages)

Only the **'internal'** folder contains source code, the other 2 are necessary to run the application.

# How-to

To build in DEV-mode, go to the root dir and run: `go run main.go`.

To build an executable, run: `go build`. Make sure you distribute both the **'web'** and **'db'** folder along with your executable.
