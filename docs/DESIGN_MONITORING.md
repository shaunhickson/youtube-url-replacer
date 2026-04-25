# Design: Cloud Monitoring & Alerting

**Issue:** #45
**Status:** Approved

## Overview
To ensure the reliability of the LinkLens Cloud Run service, we need automated alerting for critical metrics. Defining these alerts manually in the GCP Console leads to configuration drift. Instead, we will define these monitoring policies declaratively.

## Implementation Details
Since Issue #44 tracks the overall Terraform setup, we will introduce a `terraform/` directory starting with our monitoring configurations. This allows us to apply monitoring policies programmatically.

### Monitored Metrics
We are implementing three specific alerting policies using `google_monitoring_alert_policy` and `google_monitoring_uptime_check_config`:
1.  **Uptime Check:** Periodically pings the `/health` endpoint to ensure the Cloud Run service is reachable.
2.  **Latency Alert:** Triggers if the 99th percentile (P99) request latency exceeds 2 seconds over a 5-minute rolling window.
3.  **Error Rate Alert:** Triggers if the ratio of HTTP 5xx responses to total responses exceeds 1% over a 5-minute rolling window.

### Delivery
- `terraform/monitoring.tf`: Contains the monitoring resources.
- `terraform/variables.tf`: Contains required inputs (like `project_id` and `cloud_run_service_name`).
