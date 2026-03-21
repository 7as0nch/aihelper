FROM nginx:stable-alpine AS production

# 设置时区和 JS 模块支持
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    echo "types { application/javascript js mjs; }" > /etc/nginx/conf.d/mjs.conf

# 从构建阶段拷贝产物
# 注意：Vben 默认构建产物可能在应用目录下的 dist
COPY apps/web-antd/dist /usr/share/nginx/html
COPY scripts/deploy/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
