#!/bin/sh

# 应用名称
app_name="app_test"

# 镜像版本号
app_version="v1.0.0"

#镜像仓库所属项目
harbor_path="/zhaochangjiang"

#镜像仓库地址
image_warehouse_address="repo.xxx.com:8089"
#镜像仓库地址用户名
image_warehouse_username="admin"
#镜像仓库密码
image_warehouse_password="Harbor12345"

#镜像名称
app_image_id="golang_release"

#生成的镜像仓库全路径
target_image_address="${image_warehouse_address}${harbor_path}/${app_image_id}:${app_version}"

echo  "【INFO】build app ${app_name}"
echo  "【INFO】image_warehouse_address ${target_image_address}"
echo  "【INFO】go get ./..."
project_path=$(cd `dirname $0`; pwd)"/go.mod"
echo project_path
if [ ! -f "$project_path" ]; then
  echo  "【INFO】go mod init"
  go mod init
fi
go get ./...
#交叉编译，编译成linux平台可执行文件
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${app_name} ./
echo "【INFO】set 'start.sh' can be executable"

# 启动文件可执行授权
chmod +x ./start.sh

#删除none:none镜像
echo "删除none:none镜像"
docker rmi -f $(docker images -f "dangling=true" -q)
echo "删除${target_image_address}镜像"
docker rmi -f ${target_image_address}
echo  "【INFO】build images"
#构建镜像
docker build -t ${app_image_id} .

echo  "【INFO】build images"
echo "docker login -u=${image_warehouse_username} -p ${image_warehouse_password} ${image_warehouse_address}"
#登录镜像仓库
docker login -u=${image_warehouse_username} -p ${image_warehouse_password} ${image_warehouse_address}
echo  "docker tag ${app_image_id} ${target_image_address}"
#给镜像打上标签
docker tag ${app_image_id} ${target_image_address}
echo  "docker push ${target_image_address}"
#推送镜像到仓库
docker push ${target_image_address}
echo  "【INFO】docker push images finished"
echo  "【INFO】Delete location images"
docker rmi -f ${app_image_id}
echo  "【INFO】start run container"
echo  "【docker pull ${target_image_address}"
docker pull ${target_image_address}
#docker run 后边的执行脚本会替换 Dockerfile中的CMD命令
#docker run -it -p 80:80 ${app_image_id} /bin/sh

docker run -it -p 80:80 ${target_image_address}
echo "【INFO】container stoped"
