#sh 批量删除docker文件
docker stop $(docker ps)
docker rm $(docker ps -a -q)
docker rmi $(docker images -q)