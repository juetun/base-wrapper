常用命令

```shell
p=$(supervisorctl pid ad);echo "pid:$p" && lsof -p $p |wc -l

```
