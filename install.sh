#!/bin/bash

# download latest binary
curl -O https://github.com/snowskeleton/kumad/blob/main/kumad

# clean up old installs
sudo systemctl stop kumad
sudo systemctl disable kumad
sudo rm -rf /usr/local/bin/kumad

#install fresh copy
sudo mv kumad /usr/local/bin/kumad
sudo touch /etc/kumad.yaml
echo "Download complete!"
echo 'Start with '
echo -e "\tsudo kumad up --push_url '<Uptime Kuma URL>'"
