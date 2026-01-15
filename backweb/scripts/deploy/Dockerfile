FROM node:20-slim AS builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV NODE_OPTIONS=--max-old-space-size=8192
ENV TZ=Asia/Shanghai

RUN corepack enable

WORKDIR /app

# 在 Monorepo 中，为了利用缓存，我们需要拷贝所有 package.json
# 但由于目录层级较深，直接全量拷贝配置后再安装依赖是更稳妥的做法
# 如果追求极致缓存，建议在构建机上先执行 turbo prune
COPY . .

# 使用缓存挂载加速 pnpm install
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

# 构建指定的管理端应用
RUN pnpm run build --filter=@vben/web-antd

FROM nginx:stable-alpine AS production

# 设置时区和 JS 模块支持
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    echo "types { application/javascript js mjs; }" > /etc/nginx/conf.d/mjs.conf

# 从构建阶段拷贝产物
# 注意：Vben 默认构建产物可能在应用目录下的 dist
COPY --from=builder /app/apps/web-antd/dist /usr/share/nginx/html
COPY --from=builder /app/scripts/deploy/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
