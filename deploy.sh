set -eu

scripts=$(cat <<EOF
cd /home/ec2-user/snowforecast-twitter-bot/
sudo git pull origin master
sudo /usr/local/bin/docker-compose restart
exit;
EOF
)

# Deploy
echo $scripts | ssh -i $1 "ec2-user@ec2-18-181-161-220.ap-northeast-1.compute.amazonaws.com" bash -l -eux -s
