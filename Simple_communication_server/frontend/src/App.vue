a<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import StartService from './components/StartService.vue'
import DataPage from './components/DataPage.vue'
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'

// 当前选中的页面，默认为'startup'（启动页）
const currentPage = ref('startup')

// 切换页面的方法
function switchPage(pageName) {
  currentPage.value = pageName
}

// 全局事件监听 - 应用启动时注册
onMounted(() => {
  // 监听后端发射的 "chat:update" 事件（全局监听）
  EventsOn('chat:update', (data) => {
    console.log('🌐 全局收到更新：', data)
    try {
      const parsedData = JSON.parse(data)
      // 只保留最新一条的信息，添加到日志
      if (parsedData.length > 0) {
        const latestRecord = parsedData[parsedData.length - 1]
        const logEntry = {
          time: latestRecord.time,
          message: latestRecord.date,
          type: 'info'
        }
        
        // 确保全局日志数组存在
        if (!window.chatLogs) {
          window.chatLogs = []
        }
        
        window.chatLogs.push(logEntry)
        
        // 触发所有页面的数据更新
        window.dispatchEvent(new CustomEvent('chatLogUpdate', {
          detail: logEntry
        }))
      }
    } catch (error) {
      console.error('❌ 全局解析 JSON 失败：', error)
      const errorLog = {
        time: new Date().toLocaleString(),
        message: `解析错误: ${error.message}`,
        type: 'error'
      }
      
      if (!window.chatLogs) {
        window.chatLogs = []
      }
      
      window.chatLogs.push(errorLog)
      
      // 触发所有页面的数据更新
      window.dispatchEvent(new CustomEvent('chatLogUpdate', {
        detail: errorLog
      }))
    }
  })
})

onUnmounted(() => {
  // 应用关闭时取消监听
  EventsOff('chat:update')
})
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
