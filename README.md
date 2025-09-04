# Labctl command line utility

For information on ho  to use Labctl, see the user documentation at <https://labctl.net>.

This README file contains info on the labctl internals and should be the starting point if you want to contribute to the project. The following repositories form part of the labctl project:

- [labctl/labctl](https://github.com/labctl/labctl) - the main labctl go command line utility
- [labctl/labctl-frontend](https://github.com/labctl/labctl-frontend) - the web frontend - built &amp; uploaded [/helpers/frontend/html](https://github.com/labctl/labctl/tree/main/helpers/frontend/html)
- [labctl/labctl-docs](https://github.com/labctl/labctl-docs) - the labctl.net website
- [labctl/labctl-examples](https://github.com/labctl/labctl-examples) - example labs, topology files and templates

## Config Engine

Labctl uses the Config Engine from containerlab to process magic variables and render templates

The go text/template library is extended with user defined functions from:

- <https://pkg.go.dev/github.com/srl-labs/containerlab>
- <https://pkg.go.dev/github.com/kellerza/template>

### Config Engine Issues

Feel free to open any Config Engine related issue on labctl. You can also view open issues on [containerlab](https://github.com/srl-labs/containerlab/labels/config%20engine)

## labctl serve Web API

The server paths are defined in [cmd_serve.go](./app/cmd_serve.go) and used by the labctl-frontend

| path       | Description                                     |
|------------|-------------------------------------------------|
| /ws        | Websocket for commands, updates, rendering, etc |
| /wspty     | PTY via Websocket, for SSH/terminal sessions    |
| /topo      | The Topology file for the graph                 |
| /vars      |                                                 |
| /templates |                                                 |
| /files     |                                                 |
| /error     |                                                 |

### Websocket messages

The websocket messages used on the `/ws` API endpoint are defined in [ws_message.go](./helpers/ws_message.go)

| code        | sent by | description                                                                          |
|-------------|---------|--------------------------------------------------------------------------------------|
| error, warn | server  | Display messages in the frontend                                                     |
| template    | both    | Render a template on the server as&when requested by the frontend                    |
| uidata      | both    | The frontend uses this msg to save persistent data in the local labctl yml file file |
| config      | both    | Execute a config command on the server and return the results to the frontend        |
| fschange    | server  | Notify the frontend of file updates (README, topology, etc.)                         |
