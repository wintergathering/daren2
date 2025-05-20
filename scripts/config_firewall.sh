gcloud compute firewall-rules create allow-daren-app-tcp8080 \
    --project=dares-app-346910 \
    --allow tcp:8080 \
    --source-ranges=0.0.0.0/0 \
    --target-tags=http-server --description="Allow Daren app on TCP 8080"
