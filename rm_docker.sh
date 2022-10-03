docker stop $(docker ps)
docker rm $(docker ps -a -q)
docker rmi $(docker images -q)