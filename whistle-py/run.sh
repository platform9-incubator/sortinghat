cd /home/ubuntu/pf9-infra/misc/fuzzy-match


#Without nginx
#sudo uwsgi --socket 0.0.0.0:80  --protocol=http --enable-thread -w web_server:app --daemonize logs/access.log

#With nginx
sudo uwsgi --socket 127.0.0.1:3032 -w web_server:app --enable-thread --daemonize logs/access.log
