#!/bin/bash
cd /Users/xuehao/Desktop/temp/gf_api_auth_log-master/gf_api

echo "=== 初始化Git仓库 ==="
git init

echo ""
echo "=== 配置Git用户信息 ==="
git config user.name "Developer" || git config --global user.name "Developer"
git config user.email "developer@example.com" || git config --global user.email "developer@example.com"

echo ""
echo "=== 添加文件到Git ==="
git add .

echo ""
echo "=== 查看待提交的文件 ==="
git status --short | head -20

echo ""
echo "=== 创建初始提交 ==="
git commit -m "Initial commit: GoFrame API project

- 基于GoFrame v2框架的HTTP API服务
- 支持PostgreSQL和Redis
- 实现了多个API接口（设备历史、参数变更历史等）
- 统一的路由管理和外部服务调用封装
- 配置了本地开发环境（PostgreSQL和Redis）"

echo ""
echo "=== Git状态 ==="
git status

echo ""
echo "=== 提交历史 ==="
git log --oneline -1

echo ""
echo "✅ Git仓库初始化完成！"

