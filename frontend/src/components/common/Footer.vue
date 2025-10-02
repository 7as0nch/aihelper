<template>
  <footer class="footer">
    <div class="footer-content">
      <p>&copy; 2025 AI Chat. All rights reserved.</p>
      <div class="footer-links">
        <a href="#privacy">隐私政策</a>
        <a href="#terms">服务条款</a>
        <a href="#contact">联系我们</a>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
// Footer组件逻辑
</script>

<style scoped>
.footer {
  background-color: #f8f9fa;
  color: #6c757d;
  padding: 2rem 0;
  border-top: 1px solid #e9ecef;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 2rem;
}

.footer-links a {
  color: #6c757d;
  text-decoration: none;
  margin-left: 1rem;
  transition: color 0.3s;
}

.footer-links a:hover {
  color: #42b983;
}
</style>