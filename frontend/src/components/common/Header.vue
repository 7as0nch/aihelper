<template>
  <header class="header">
    <div class="header-content">
      <h1 class="logo">AI Chat</h1>
      <nav class="nav-menu">
        <ul>
          <li><a href="#home">首页</a></li>
          <li><a href="#about">关于</a></li>
          <li><a href="#contact">联系</a></li>
        </ul>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
// Header组件逻辑
</script>

<style scoped>
.header {
  background-color: #42b983;
  color: white;
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 2rem;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  margin: 0;
}

.nav-menu ul {
  display: flex;
  list-style: none;
}

.nav-menu li {
  margin-left: 2rem;
}

.nav-menu a {
  color: white;
  text-decoration: none;
  transition: opacity 0.3s;
}

.nav-menu a:hover {
  opacity: 0.8;
}
</style>