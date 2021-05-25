FROM alpine

# 程序初始目录
WORKDIR /var/run/

# 将当前目录拷贝到指定文件夹下
COPY ./ /var/run/

# 配置运行环境变量
ENV "GO_ENV" "dev"

# 配置启动脚本
CMD [ "./start.sh" ]
