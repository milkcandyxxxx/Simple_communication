a<script setup>
import { ref } from 'vue'
import StartService from './components/StartService.vue'
import DataPage from './components/DataPage.vue'

// 当前选中的页面，默认为'startup'（启动页）
const currentPage = ref('startup')

// 切换页面的方法
function switchPage(pageName) {
  currentPage.value = pageName
}
</script>

<template>
  <div class="app">
    <aside class="sidebar">
      <nav class="nav">
        <div class="nav-item" :class="{active: currentPage === 'startup'}" @click="switchPage('startup')">
          <span>启动</span>
        </div>
        <div class="nav-item" :class="{active: currentPage === 'data'}" @click="switchPage('data')">
          <span>数据</span>
        </div>
      </nav>
    </aside>
    <main class="main">
      <StartService v-if="currentPage === 'startup'" />
      <DataPage v-else-if="currentPage === 'data'" />
    </main>
  </div>
</template>

<style>
.app {
  display: grid;
  grid-template-columns: 220px 1fr;
  height: 100vh;
  background: #fff;
  color: #000;
}

.sidebar {
  background: #fff;
  border-right: 3px solid #000;
  box-shadow: 8px 0 0 #000;
}

.nav {
  padding: 16px;
}

.nav-item {
  background: #fff;
  border: 3px solid #000;
  box-shadow: 6px 6px 0 #000;
  padding: 12px;
  margin-bottom: 12px;
  cursor: pointer;
  user-select: none;
}

.nav-item.active {
  background: #e8f6ff;
}

.main {
  padding: 24px;
}
</style>
