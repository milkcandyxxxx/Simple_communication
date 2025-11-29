<script setup>
import { ref, onMounted } from 'vue'
import { Start, Stop } from '../../wailsjs/go/main/App'

// 全局状态存储（切换页面不丢失）
if (!window.serviceState) {
  window.serviceState = {
    isStarted: false,
    serviceName: '',
    resultText: '请点击启动按钮启动 WebSocket 服务'
  }
}

const isStarted = ref(window.serviceState.isStarted)
const serviceName = ref(window.serviceState.serviceName)
const resultText = ref(window.serviceState.resultText)
const isLoading = ref(false)

function handleStart() {
  isLoading.value = true
  
  Start('WebSocket Server').then(result => {
    isLoading.value = false
    if (result === 0) {
      isStarted.value = true
      serviceName.value = 'WebSocket Server'
      resultText.value = '✅ WebSocket 服务已启动！\n访问地址: ws://localhost:8080/ws'
      
      // 保存到全局状态
      window.serviceState.isStarted = true
      window.serviceState.serviceName = 'WebSocket Server'
      window.serviceState.resultText = resultText.value
    } else {
      resultText.value = '❌ 服务启动失败，请查看控制台日志'
    }
  }).catch(err => {
    isLoading.value = false
    resultText.value = `❌ 服务启动失败: ${err}`
  })
}

function handleStop() {
  isLoading.value = true
  
  Stop().then(result => {
    isLoading.value = false
    if (result === 0) {
      isStarted.value = false
      resultText.value = '⛔ WebSocket 服务已关闭'
      
      // 更新全局状态
      window.serviceState.isStarted = false
      window.serviceState.resultText = resultText.value
    } else {
      resultText.value = '❌ 关闭服务失败，请查看控制台日志'
    }
  }).catch(err => {
    isLoading.value = false
    resultText.value = `❌ 关闭服务失败: ${err}`
  })
}

onMounted(() => {
  // 恢复全局状态
  isStarted.value = window.serviceState.isStarted
  serviceName.value = window.serviceState.serviceName
  resultText.value = window.serviceState.resultText
})
</script>

<template>
  <div class="start-page">
    <div class="header">
      <h1 class="title">🚀 服务启动控制台</h1>
      <div class="status-badge" :class="{ active: isStarted }">
        <span class="status-dot"></span>
        <span>{{ isStarted ? '运行中' : '未启动' }}</span>
      </div>
    </div>

    <div class="content">
      <div class="result-box" :class="{ success: isStarted }">
        <pre class="result-text">{{ resultText }}</pre>
      </div>

      <div class="control-panel">
        <button 
          class="btn btn-primary" 
          @click="handleStart"
          :disabled="isStarted || isLoading"
        >
          <span v-if="isLoading">启动中...</span>
          <span v-else-if="isStarted">✅ 已启动</span>
          <span v-else>🚀 启动服务</span>
        </button>
        
        <button 
          class="btn btn-danger" 
          @click="handleStop"
          :disabled="!isStarted || isLoading"
        >
          <span v-if="isLoading">关闭中...</span>
          <span v-else>⛔ 关闭服务</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.start-page {
  max-width: 800px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding-bottom: 16px;
  border-bottom: 3px solid #000;
}

.title {
  margin: 0;
  font-size: 28px;
  font-weight: 900;
  color: #000;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  border: 3px solid #000;
  box-shadow: 4px 4px 0 #000;
  padding: 8px 16px;
  font-weight: 700;
  font-size: 14px;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #999;
  border: 2px solid #000;
}

.status-badge.active {
  background: #d4f4dd;
}

.status-badge.active .status-dot {
  background: #4caf50;
  box-shadow: 0 0 8px #4caf50;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.result-box {
  background: #fff;
  border: 3px solid #000;
  box-shadow: 8px 8px 0 #000;
  padding: 20px;
  min-height: 120px;
  transition: all 0.3s;
}

.result-box.success {
  background: #d4f4dd;
  border-color: #4caf50;
  box-shadow: 8px 8px 0 #4caf50;
}

.result-text {
  margin: 0;
  font-size: 16px;
  line-height: 1.6;
  color: #000;
  font-family: 'Courier New', monospace;
  white-space: pre-wrap;
  word-break: break-word;
}

.control-panel {
  background: #fff;
  border: 3px solid #000;
  box-shadow: 8px 8px 0 #000;
  padding: 24px;
  display: flex;
  justify-content: center;
  gap: 16px;
}

.btn {
  background: #fff;
  color: #000;
  border: 3px solid #000;
  box-shadow: 4px 4px 0 #000;
  padding: 12px 32px;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 200px;
}

.btn:hover:not(:disabled) {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0 #000;
}

.btn:active:not(:disabled) {
  transform: translate(1px, 1px);
  box-shadow: 2px 2px 0 #000;
}

.btn-primary {
  background: #e8f6ff;
}

.btn-primary:hover:not(:disabled) {
  background: #b3e5fc;
}

.btn-danger {
  background: #ffe5e5;
}

.btn-danger:hover:not(:disabled) {
  background: #ffcdd2;
}

.btn:disabled {
  background: #f5f5f5;
  color: #999;
  cursor: not-allowed;
  transform: none;
}
</style>
