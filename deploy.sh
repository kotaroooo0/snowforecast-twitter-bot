#!/bin/sh
cd /home/ec2-user/snowforecast-twitter-bot/
sudo git pull origin master
sudo /usr/local/bin/docker-compose restart
