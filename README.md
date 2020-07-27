# slab_exporter
Prometheus Exporter for Slab metrics. This exporter reads from `/proc/slabinfo` 
and exports the data as prometheus metrics.

# Usage
```
slab_exporter -h
Prometheus exporter for slab metrics

Usage:
  slab_exporter [flags]

Flags:
  -c, --config string           config file (default is $HOME/.perf_exporter.yaml)
  -l, --listen-address string   Server listen address (default "0.0.0.0:8585")
      --metrics-path string     Metrics endpoint (default "/metrics")
  -r, --regex string            collect slabs matching regex

Build Timestamp:
  2020-07-22T22:34:49-04:00
Go Version:
  go1.14.2
```

# Example

```
curl -s localhost:8585/metrics | grep -i active | head
# HELP slab_active_objs slab active_objs
# TYPE slab_active_objs gauge
slab_active_objs{slab="AF_VSOCK"} 0
slab_active_objs{slab="Acpi_Namespace"} 2635
slab_active_objs{slab="Acpi_Operand"} 5768
slab_active_objs{slab="Acpi_Parse"} 803
slab_active_objs{slab="Acpi_ParseExt"} 561
slab_active_objs{slab="Acpi_State"} 561
slab_active_objs{slab="L2TP_IP"} 0
slab_active_objs{slab="L2TP_IPv6"} 0
```
