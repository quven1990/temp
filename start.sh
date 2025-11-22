#!/bin/bash

cd /Users/xuehao/Desktop/temp/gf_api_auth_log-master/gf_api

echo "=========================================="
echo "  启动 GoFrame HTTP 服务"
echo "=========================================="
echo ""

# 检查PostgreSQL和Redis服务
echo "检查服务状态..."
if brew services list | grep -q "postgresql.*started"; then
    echo "✓ PostgreSQL服务运行中"
else
    echo "⚠️  PostgreSQL服务未运行，正在启动..."
    brew services start postgresql@16
fi

if brew services list | grep -q "redis.*started"; then
    echo "✓ Redis服务运行中"
else
    echo "⚠️  Redis服务未运行，正在启动..."
    brew services start redis
fi

echo ""
echo "确保配置文件存在..."
mkdir -p manifest/config
if [ ! -f manifest/config/config.yaml ]; then
    cp config/config.yaml manifest/config/config.yaml
    echo "✓ 配置文件已复制到 manifest/config/"
fi

echo ""
echo "正在启动HTTP服务..."
echo "服务地址: http://localhost:8001"
echo "API文档: http://localhost:8001/swagger"
echo ""
echo "按 Ctrl+C 停止服务"
echo "=========================================="
echo ""

# 运行服务
/opt/homebrew/bin/go run main.go

