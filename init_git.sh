#!/bin/bash

echo "=== 初始化Git仓库 ==="
git init

echo ""
echo "=== 配置Git用户信息 ==="
git config user.name "Developer"
git config user.email "developer@example.com"

echo ""
echo "=== 添加所有文件 ==="
git add .

echo ""
echo "=== 创建初始提交 ==="
git commit -m "Initial commit: GoFrame API project

- 基于GoFrame v2框架的HTTP API服务
- 支持PostgreSQL和Redis
- 实现了多个API接口（设备历史、参数变更历史等）
- 统一的路由管理和外部服务调用封装
- 配置了本地开发环境（PostgreSQL和Redis）"

echo ""
echo "=== 验证Git状态 ==="
git status

echo ""
echo "=== 查看提交历史 ==="
git log --oneline -1

echo ""
echo "✅ Git仓库初始化完成！"
echo ""
echo "现在可以执行: git status"
