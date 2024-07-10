#!/bin/bash

# clean up old installs
sudo systemctl stop kumad
sudo systemctl disable kumad
sudo rm -rf /usr/local/bin/kumad

#install fresh copy
sudo cp kumad /usr/local/bin/kumad
sudo touch /etc/kumad.yaml
