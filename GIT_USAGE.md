# Git仓库使用说明

## 当前状态

✅ Git仓库已初始化
✅ 所有文件已提交到本地仓库
✅ 当前分支：master

## 查看Git状态

```bash
git status
```

## 查看提交历史

```bash
git log --oneline
```

## 连接到远程仓库（可选）

如果你有GitHub、GitLab等远程仓库，可以这样连接：

```bash
# 添加远程仓库
git remote add origin <你的仓库URL>

# 推送到远程仓库
git push -u origin master
```

## 常用Git命令

```bash
# 查看状态
git status

# 添加文件
git add .

# 提交更改
git commit -m "提交说明"

# 查看提交历史
git log --oneline

# 查看分支
git branch

# 创建新分支
git checkout -b feature/新功能

# 切换分支
git checkout master
```

## 注意事项

- `.gitignore` 文件已配置，会排除敏感配置文件（如 `config/config.yaml`）
- 如果修改了配置文件，记得不要提交敏感信息
- 建议定期提交代码，保持提交信息清晰

