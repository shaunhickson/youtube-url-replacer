resource "google_monitoring_notification_channel" "email_alert" {
  display_name = "Email Alert Channel"
  type         = "email"
  labels = {
    email_address = var.alert_email
  }
}

resource "google_monitoring_uptime_check_config" "https_uptime_check" {
  display_name = "LinkLens API Health Check"
  timeout      = "10s"
  period       = "60s"

  http_check {
    path           = "/health"
    port           = 443
    use_ssl        = true
    validate_ssl   = true
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = var.project_id
      host       = "${var.cloud_run_service_name}-${var.project_id}.a.run.app"
    }
  }
}

resource "google_monitoring_alert_policy" "uptime_alert" {
  display_name = "LinkLens Uptime Alert"
  combiner     = "OR"
  conditions {
    display_name = "Uptime Check Failure"
    condition_threshold {
      filter     = "metric.type=\"monitoring.googleapis.com/uptime_check/check_passed\" AND resource.type=\"uptime_url\" AND metric.labels.check_id=\"${google_monitoring_uptime_check_config.https_uptime_check.uptime_check_id}\""
      duration   = "300s"
      comparison = "COMPARISON_LT"
      threshold_value = 1
      aggregations {
        alignment_period   = "60s"
        cross_series_reducer = "REDUCE_FRACTION_TRUE"
        per_series_aligner   = "ALIGN_NEXT_OLDER"
      }
    }
  }
  notification_channels = [google_monitoring_notification_channel.email_alert.name]
}

resource "google_monitoring_alert_policy" "latency_alert" {
  display_name = "LinkLens High Latency (P99 > 2s)"
  combiner     = "OR"
  conditions {
    display_name = "P99 Latency > 2s"
    condition_threshold {
      filter     = "metric.type=\"run.googleapis.com/request_latencies\" AND resource.type=\"cloud_run_revision\" AND resource.labels.service_name=\"${var.cloud_run_service_name}\""
      duration   = "60s"
      comparison = "COMPARISON_GT"
      threshold_value = 2000
      aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_PERCENTILE_99"
        per_series_aligner   = "ALIGN_DELTA"
      }
    }
  }
  notification_channels = [google_monitoring_notification_channel.email_alert.name]
}

resource "google_monitoring_alert_policy" "error_rate_alert" {
  display_name = "LinkLens High 5xx Error Rate (>1%)"
  combiner     = "OR"
  conditions {
    display_name = "5xx Error Rate > 1%"
    condition_threshold {
      filter     = "metric.type=\"run.googleapis.com/request_count\" AND resource.type=\"cloud_run_revision\" AND resource.labels.service_name=\"${var.cloud_run_service_name}\" AND metric.labels.response_code_class=\"5xx\""
      duration   = "60s"
      comparison = "COMPARISON_GT"
      threshold_value = 0.01
      
      denominator_filter = "metric.type=\"run.googleapis.com/request_count\" AND resource.type=\"cloud_run_revision\" AND resource.labels.service_name=\"${var.cloud_run_service_name}\""
      
      aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }

      denominator_aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }
  notification_channels = [google_monitoring_notification_channel.email_alert.name]
}
