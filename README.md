# unifi-list-devices

A small CLI that connects to a UniFi controller and prints the list of
clients seen on a site as a table.

## Build

```sh
go build
```

## Usage

```
Usage of ./unifi-list-devices:
  -f, --fields strings    Comma-separated fields to display
  -p, --password string   Unifi controller password
  -s, --site string       Site name (default "default")
      --sort string       Sort by a field (default "name")
  -U, --url string        Unifi controller URL (default "http://127.0.0.1:8443")
  -u, --username string   Unifi controller username
```

### Fields

Both `--fields` and `--sort` accept the same vocabulary (with `num` and `id`
being display-only):

| Name         | Description                                  |
|--------------|----------------------------------------------|
| `num`        | Row number (display only)                    |
| `id`         | Client ID                                    |
| `ip`         | IP address                                   |
| `hostname`   | Hostname reported by the client              |
| `name`       | Name assigned in the controller              |
| `mac`        | MAC address                                  |
| `last-seen`  | Timestamp of the most recent activity        |
| `uptime`     | How long the client has been connected       |
| `first-seen` | Timestamp of the first time the client appeared |
| `assoc-time` | Timestamp of the current association         |

The default field set is
`num,id,ip,hostname,name,mac,last-seen,uptime`, and the default sort is
by `name`.

## Examples

List all clients on the default site:

```sh
./unifi-list-devices -u admin -p secret -U https://unifi.example.com:8443
```

Show only a few columns, sorted by how long each client has been
connected:

```sh
./unifi-list-devices -u admin -p secret \
    -f name,ip,uptime,last-seen --sort uptime
```
