resource "google_secret_manager_secret" "youtube_api_key" {
  secret_id = "YOUTUBE_API_KEY"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "youtube_api_key_version" {
  secret      = google_secret_manager_secret.youtube_api_key.id
  secret_data = var.youtube_api_key
}
