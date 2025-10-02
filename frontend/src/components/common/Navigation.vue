<template>
  <nav class="navigation">
    <ul class="nav-list">
      <li class="nav-item">
        <a href="#dashboard" class="nav-link">
          <i class="icon-dashboard"></i>
          <span>仪表盘</span>
        </a>
      </li>
      <li class="nav-item">
        <a href="#users" class="nav-link">
          <i class="icon-users"></i>
          <span>用户管理</span>
        </a>
      </li>
      <li class="nav-item">
        <a href="#content" class="nav-link">
          <i class="icon-content"></i>
          <span>内容管理</span>
        </a>
      </li>
      <li class="nav-item">
        <a href="#settings" class="nav-link">
          <i class="icon-settings"></i>
          <span>系统设置</span>
        </a>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
// Navigation组件逻辑
</script>

<style scoped>
.navigation {
  background-color: #ffffff;
  box-shadow: 2px 0 4px rgba(0,0,0,0.1);
  height: 100vh;
  width: 250px;
  position: fixed;
  top: 0;
  left: 0;
  overflow-y: auto;
}

.nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  border-bottom: 1px solid #e9ecef;
}

.nav-link {
  display: flex;
  align-items: center;
  padding: 1rem 1.5rem;
  text-decoration: none;
  color: #495057;
  transition: all 0.3s;
}

.nav-link:hover {
  background-color: #f8f9fa;
  color: #42b983;
}

.nav-link i {
  margin-right: 0.75rem;
  font-size: 1.2rem;
}
</style>