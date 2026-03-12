# Jenkinsfile CI/CD 部署教程

本文档介绍了如何使用本项目中的 `jenkinsfile` 来实现前端、后端、管理端服务的自动化构建与部署。

## 1. Jenkinsfile 简介

`jenkinsfile` 是 Jenkins 的流水线（Pipeline）脚本文件，基于 Groovy 语法编写。它定义了整个持续集成和持续部署（CI/CD）的流程。
本项目的 Jenkinsfile 支持**参数化构建**，允许你在触发构建时手动选择要部署的服务。

## 2. 支持的构建参数

在 Jenkins 界面点击“Build with Parameters”（参数化构建）时，你将看到以下选项：

*   **SERVICE_TO_DEPLOY** (选择框): 选择需要部署的服务。
    *   `all`: 同时部署前端、后端和管理端。
    *   `frontend`: 仅部署前端服务。
    *   `backend`: 仅部署后端服务。
    *   `admin`: 仅部署管理端服务。
*   **BRANCH_NAME** (文本框): 指定要拉取并构建的 Git 代码分支（默认 `main`）。
*   **IMAGE_TAG** (文本框): 指定本次构建生成的 Docker 镜像版本号（默认 `latest`，建议生产环境使用如 `v1.0.0` 等具体版本号）。

## 3. 如何在 Jenkins 中配置

请按照以下步骤在 Jenkins 中创建一个使用此 Jenkinsfile 的任务：

1.  **新建任务 (New Item)**：
    *   在 Jenkins 首页点击左侧的“新建任务”。
    *   输入任务名称（例如：`aichat-deploy-pipeline`）。
    *   选择 **Pipeline (流水线)**，然后点击“确定”。
2.  **配置流水线 (Pipeline)**：
    *   在任务配置页面，向下滚动到 **Pipeline** 区域。
    *   **Definition** 选择 `Pipeline script from SCM`（从代码仓库获取流水线脚本）。
    *   **SCM** 选择 `Git`。
    *   **Repository URL**: 输入你的 Git 仓库地址（例如 `https://github.com/your-repo/aichat.git`）。
    *   **Credentials**: 如果仓库是私有的，请添加并选择相应的凭据。
    *   **Branches to build**: 填写 `*/main` 或你默认的主分支。
    *   **Script Path**: 填写 `jenkinsfile`（与代码仓库中该文件的路径和名称保持一致）。
3.  **保存并初始化参数**：
    *   点击“保存”。
    *   **重要提示**：第一次运行时，Jenkins 可能还不知道有哪些参数。你需要先点击一次 **Build Now**（立即构建）。这次构建可能会失败或只是读取参数，但之后左侧菜单会变成 **Build with Parameters**（参数化构建）。

## 4. 阶段 (Stages) 说明

流水线包含以下几个主要阶段，使用了 `when` 条件语句来判断是否需要跳过：

*   **Checkout Code**: 从 Git 仓库拉取最新代码。
*   **Build & Deploy Frontend**: 仅当选择 `frontend` 或 `all` 时执行。包含前端的依赖安装、打包、构建 Docker 镜像并推送到仓库，最后更新 Kubernetes 部署。
*   **Build & Deploy Backend**: 仅当选择 `backend` 或 `all` 时执行。包含 Go 后端代码的编译、构建镜像并更新部署。
*   **Build & Deploy Admin**: 仅当选择 `admin` 或 `all` 时执行。流程与前端类似。

## 5. 如何自定义修改

为了让脚本真正在你的环境中运行，你需要修改 `jenkinsfile` 中的一些占位符：

1.  **镜像仓库地址**：修改 `environment` 块中的 `DOCKER_REGISTRY` 为你实际的 Docker 镜像仓库地址。
2.  **构建与部署命令**：在各个 stage 的 `steps` 中，取消注释 `sh ''' ... '''` 块，并根据你实际的项目结构、Dockerfile 路径以及 Kubernetes 部署文件名（如 `k8s-deployment.yaml`）调整 Shell 命令。
3.  **通知配置**：在 `post` 块中，你可以根据需要集成钉钉、飞书或邮件通知插件，以便在部署成功或失败时收到提醒。