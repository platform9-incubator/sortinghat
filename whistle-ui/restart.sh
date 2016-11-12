#Without nginx
#sudo kill -15 $(pgrep -f "sudo uwsgi --socket 0.0.0.0:80 --protocol=http --enable-thread -w web_server:app")
#sudo kill -15 $(pgrep -f "uwsgi --socket 0.0.0.0:80 --protocol=http --enable-thread -w web_server:app")

#With nginx
sudo kill -15 $(pgrep -f "sudo uwsgi --socket 127.0.0.1:3032 -w web_server:app --enable-thread --daemonize logs/access.log")
sudo kill -15 $(pgrep -f "uwsgi --socket 127.0.0.1:3032 -w web_server:app --enable-thread --daemonize logs/access.log")
#sudo service fuzzy stop
sudo python /home/ubuntu/pf9-infra/misc/fuzzy-match/cleanup.py
#sudo service fuzzy start
sudo /home/ubuntu/pf9-infra/misc/fuzzy-match/run.sh &
