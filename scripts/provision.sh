gcloud compute instances create daren-app-vm \
    --project=dares-app-346910 \
    --zone=us-east4-a \
    --machine-type=e2-micro \
    --image-family=debian-11 \
    --image-project=debian-cloud \
    --boot-disk-size=20GB \
    --tags=http-server,https-server
