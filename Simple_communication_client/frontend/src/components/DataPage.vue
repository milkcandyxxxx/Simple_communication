<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

// 全局日志存储（组件外部，切换页面不会丢失）
if (!window.chatLogs) {
  window.chatLogs = []
}

const chatLogs = ref(window.chatLogs)
const connectionStatus = ref('WAITING')
const consoleRef = ref(null)

// 自动滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (consoleRef.value) {
      consoleRef.value.scrollTop = consoleRef.value.scrollHeight
    }
  })
}

onMounted(() => {
  // 监听全局自定义事件
  const handleChatLogUpdate = (event) => {
    console.log('📊 DataPage 收到更新：', event.detail)
    connectionStatus.value = 'CONNECTED'
    // 触发响应式更新：创建新数组引用
    chatLogs.value = [...window.chatLogs]
    scrollToBottom()
  }
  
  window.addEventListener('chatLogUpdate', handleChatLogUpdate)
  
  // 初始化日志
  if (chatLogs.value.length === 0) {
    const initLog = {
      time: new Date().toLocaleString(),
      message: '控制台已就绪，等待 WebSocket 连接...',
      type: 'system'
    }
    window.chatLogs.push(initLog)
    chatLogs.value = [...window.chatLogs]
  }
  
  scrollToBottom()
  
  // 保存事件处理函数引用，用于卸载时移除
  window.handleChatLogUpdate = handleChatLogUpdate
})

onUnmounted(() => {
  // 组件卸载时移除事件监听
  if (window.handleChatLogUpdate) {
    window.removeEventListener('chatLogUpdate', window.handleChatLogUpdate)
  }
})

// 清空日志
const clearLogs = () => {
  window.chatLogs = []
  const clearLog = {
    time: new Date().toLocaleString(),
    message: '日志已清空',
    type: 'system'
  }
  window.chatLogs.push(clearLog)
  // 触发响应式更新：创建新数组引用
 chatLogs.value = [...window.chatLogs]
scrollToBottom()
}
</script>

<template>
  <div class="console-page">
    <div class="console-header">
      <div class="header-left">
        <span class="console-icon">▶</span>
        <h1 class="title">控制台日志</h1>
      </div>
      <div class="header-right">
        <div class="status-indicator" :class="connectionStatus.toLowerCase()">
          <span class="status-dot"></span>
          <span class="status-text">{{ connectionStatus }}</span>
        </div>
        <button class="clear-btn" @click="clearLogs">
          <span>清空日志</span>
        </button>
      </div>
    </div>

    <div class="console-wrapper" ref="consoleRef">
      <div class="log-list">
        <div
          v-for="(log, index) in chatLogs"
          :key="index"
          class="log-entry"
          :class="log.type"
        >
          <span class="log-time">{{ log.time }}</span>
          <span class="log-type" :class="log.type">[{{ log.type.toUpperCase() }}]</span>
          <span class="log-message">{{ log.message }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.console-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
}

.console-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #252526;
  border: 3px solid #000;
  border-bottom: 4px solid #000;
  box-shadow: 0 4px 0 #000;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.console-icon {
  font-size: 24px;
  color: #4caf50;
  font-weight: 900;
}

.title {
  margin: 0;
  font-size: 24px;
  font-weight: 900;
  color: #fff;
  letter-spacing: 1px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #1e1e1e;
  border: 3px solid #000;
  box-shadow: 4px 4px 0 #000;
  padding: 6px 12px;
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 2px solid #000;
  background: #666;
}

.status-indicator.connected .status-dot {
  background: #4caf50;
  box-shadow: 0 0 8px #4caf50;
}

.status-indicator.waiting .status-dot {
  background: #ffa726;
  box-shadow: 0 0 8px #ffa726;
}

.status-text {
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  font-family: 'Courier New', monospace;
}

.clear-btn {
  background: #ff6b6b;
  border: 3px solid #000;
  box-shadow: 4px 4px 0 #000;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
  cursor: pointer;
  transition: transform 0.1s, box-shadow 0.1s;
}

.clear-btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0 #000;
}

.clear-btn:active {
  transform: translate(1px, 1px);
  box-shadow: 2px 2px 0 #000;
}

.console-wrapper {
  flex: 1;
  overflow-y: auto;
  background: #1e1e1e;
  border: 3px solid #000;
  border-top: none;
  padding: 16px;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-family: 'Courier New', 'Consolas', monospace;
}

.log-entry {
  display: flex;
  gap: 12px;
  padding: 6px 8px;
  border-left: 3px solid transparent;
  transition: background 0.15s;
}

.log-entry:hover {
  background: #2d2d30;
}

.log-entry.info {
  border-left-color: #4caf50;
}

.log-entry.error {
  border-left-color: #f44336;
  background: rgba(244, 67, 54, 0.1);
}

.log-entry.system {
  border-left-color: #2196f3;
}

.log-time {
  color: #858585;
  font-size: 12px;
  white-space: nowrap;
  min-width: 180px;
}

.log-type {
  font-weight: 700;
  font-size: 12px;
  white-space: nowrap;
  min-width: 60px;
}

.log-type.info {
  color: #4caf50;
}

.log-type.error {
  color: #f44336;
}

.log-type.system {
  color: #2196f3;
}

.log-message {
  color: #d4d4d4;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
  flex: 1;
}

/* 滚动条样式 - 深色主题 */
.console-wrapper::-webkit-scrollbar {
  width: 14px;
}

.console-wrapper::-webkit-scrollbar-track {
  background: #1e1e1e;
  border-left: 3px solid #000;
}

.console-wrapper::-webkit-scrollbar-thumb {
  background: #424242;
  border: 3px solid #000;
}

.console-wrapper::-webkit-scrollbar-thumb:hover {
  background: #4caf50;
}
</style>