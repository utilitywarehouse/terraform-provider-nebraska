# Terraform Provider for Nebraska

A terraform provider for configuring [Kinvolk's Nebraska update
manager](https://github.com/kinvolk/nebraska).

## Status

Currently only a subset of the possible resources are implemented.

### Data sources

- `nebraska_channel`
- `nebraska_group`
- `nebraska_package`

### Resources

- `nebraska_channel`
- `nebraska_group`
- `nebraska_package`

## Usage

By default, the provider will attempt to connect to a Nebraska server at
`http://localhost:8000`. You can change this by setting the `endpoint`
parameter in the provider configuration or the  `NEBRASKA_ENDPOINT`
environment variable.

The provider doesn't currently support any authentication methods.

Most resources in Nebraska belong to an 'application'. You can optionally set a
default application for the provider to target with the `application_id`
(`NEBRAKSA_APPLICATION_ID`) parameter.

Tip: [the default Flatcar application is pre-created with the id
`e96281a6-d1af-4bde-9a0a-97b76e56dc57`](https://github.com/flatcar/nebraska/blob/2.8.6/backend/pkg/api/applications.go#L32).

```hcl
provider "nebraska" {
  application_id = "e96281a6-d1af-4bde-9a0a-97b76e56dc57"
  endpoint       = "http://nebraska:8000"
}
```

## Development

You can run the acceptance tests with `make testacc` (requires `docker-compose`).
