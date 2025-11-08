FROM alpine:3.8	


LABEL maintainer="john@goframe.org"

###############################################################################
#                                INSTALLATION
###############################################################################

# 使用国内alpine源,并设置时区
RUN echo "http://mirrors.aliyun.com/alpine/v3.8/main/" > /etc/apk/repositories && \
    apk update && apk add --no-cache tzdata ca-certificates bash && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone


# 设置固定的项目路径
ENV WORKDIR /app/main

# 添加应用可执行文件，并设置执行权限
ADD ./main   $WORKDIR/main
RUN chmod +x $WORKDIR/main

# 添加静态资源文件
ADD resource $WORKDIR/resource
ADD config $WORKDIR/config

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR

# 显式暴露端口
EXPOSE 80
#EXPOSE 443

CMD ./main