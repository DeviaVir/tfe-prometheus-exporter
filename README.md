# tfe-prometheus-exporter

A prometheus exporter for Terraform Enterprise.

## Docker

```
docker run -p 9112:9112 -e TFE_TOKEN=<some-admin-token> -e TFE_ADDRESS=<some-tfe-address> -e HTTP_LISTENADDR=":9112" -it --rm deviavir/tfe-prometheus-exporter:latest
```

## Example output

```
$ curl localhost:9112/metrics
[...]
# HELP tf_enterprise_runs_applied Runs applied (applied)
# TYPE tf_enterprise_runs_applied gauge
tf_enterprise_runs_applied 4950
# HELP tf_enterprise_runs_applying Runs currently applying (applying)
# TYPE tf_enterprise_runs_applying gauge
tf_enterprise_runs_applying 0
# HELP tf_enterprise_runs_canceled Runs canceled (canceled)
# TYPE tf_enterprise_runs_canceled gauge
tf_enterprise_runs_canceled 5453
# HELP tf_enterprise_runs_confirmed Runs confirmed (confirmed)
# TYPE tf_enterprise_runs_confirmed gauge
tf_enterprise_runs_confirmed 0
# HELP tf_enterprise_runs_discarded Runs discarded (discarded)
# TYPE tf_enterprise_runs_discarded gauge
tf_enterprise_runs_discarded 79
# HELP tf_enterprise_runs_errored Runs errored (errored)
# TYPE tf_enterprise_runs_errored gauge
tf_enterprise_runs_errored 12514
# HELP tf_enterprise_runs_pending Runs currently in the queue (pending)
# TYPE tf_enterprise_runs_pending gauge
tf_enterprise_runs_pending 0
# HELP tf_enterprise_runs_planned Runs planned (planned)
# TYPE tf_enterprise_runs_planned gauge
tf_enterprise_runs_planned 30
# HELP tf_enterprise_runs_planning Runs currently planning (planning)
# TYPE tf_enterprise_runs_planning gauge
tf_enterprise_runs_planning 5
# HELP tf_enterprise_runs_policy_checked Runs with checked policy (policy-checked)
# TYPE tf_enterprise_runs_policy_checked gauge
tf_enterprise_runs_policy_checked 0
# HELP tf_enterprise_runs_policy_checking Runs currently checking policy (policy-checking)
# TYPE tf_enterprise_runs_policy_checking gauge
tf_enterprise_runs_policy_checking 0
# HELP tf_enterprise_runs_policy_override Runs with overriden policy (policy-override)
# TYPE tf_enterprise_runs_policy_override gauge
tf_enterprise_runs_policy_override 0
# HELP tf_enterprise_runs_total Total number of runs with any status (total)
# TYPE tf_enterprise_runs_total gauge
tf_enterprise_runs_total 60060
```
