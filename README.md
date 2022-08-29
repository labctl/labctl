# Labctl command line utility

For more information on Labctl, see the user documentation at https://labctl.net

The docs repo can be found at https://github.com/labctl/labctl-docs

The frontend repo can be found at https://github.com/labctl/labctl-frontend - This is automatically built and uploaded to https://github.com/labctl/labctl/tree/main/helpers/frontend/html, where it will be release with labctl

## Config Engine

Labctl uses the Config Engine from containerlab to process magic variables and render templates

The go text/template library is extended with user defined functions from:
- https://docs.gomplate.ca/
- https://pkg.go.dev/github.com/kellerza/template

Config Engine Issues

Feel free to open any Config Engine related issue on labctl. You can also view open issues on containerlab [here](https://github.com/srl-labs/containerlab/labels/config%20engine)
