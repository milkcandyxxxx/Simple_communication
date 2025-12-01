<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { GreetMany, SeedMany } from '../../wailsjs/go/main/App'

// 全局状态存储（切换页面不丢失）
if (!window.connectionState) {
  window.connectionState = {
    isConnected: false,
    clientName: '',
    serverUrl: '',
    resultText: '请输入客户端名称和服务器地址进行连接'
  }
}

const isConnected = ref(window.connectionState.isConnected)
const clientName = ref(window.connectionState.clientName)
const serverUrl = ref(window.connectionState.serverUrl)
const resultText = ref(window.connectionState.resultText)
const isLoading = ref(false)
const messageToSend = ref('')
const connectionStatus = ref('未连接')

function handleConnect() {
  if (!clientName.value.trim()) {
    resultText.value = '❌ 请输入客户端名称'
    return
  }
  
  if (!serverUrl.value.trim()) {
    resultText.value = '❌ 请输入服务器地址'
    return
  }
  
  isLoading.value = true
  resultText.value = '⏳ 正在连接服务器...'
  
  GreetMany(clientName.value, serverUrl.value).then(result => {
    isLoading.value = false
    if (result === 0) {
      isConnected.value = true
      resultText.value = `✅ 成功连接到服务器 ws://${serverUrl.value}:8080/ws`
      connectionStatus.value = '已连接'
      
      // 保存到全局状态
      window.connectionState.isConnected = true
      window.connectionState.clientName = clientName.value
      window.connectionState.serverUrl = serverUrl.value
      window.connectionState.resultText = resultText.value
    } else {
      resultText.value = '❌ 连接失败，请查看控制台日志'
      connectionStatus.value = '连接失败'
    }
  }).catch(err => {
    isLoading.value = false
    resultText.value = `❌ 连接异常: ${err}`
    connectionStatus.value = '连接失败'
    console.error('连接异常:', err)
  })
}

function handleDisconnect() {
  isConnected.value = false
  resultText.value = '⛔ 已断开与服务器的连接'
  connectionStatus.value = '未连接'
  
  // 更新全局状态
  window.connectionState.isConnected = false
  window.connectionState.resultText = resultText.value
}

function handleSend() {
  if (!messageToSend.value.trim()) {
    return
  }
  
  if (!isConnected.value) {
    resultText.value = '❌ 请先连接到服务器'
    return
  }
  
  SeedMany(clientName.value, messageToSend.value).then(() => {
    const sendLog = {
      time: new Date().toLocaleString(),
      message: `已发送: ${messageToSend.value}`,
      type: 'system'
    }
    
    // 添加到全局日志
    if (!window.chatLogs) {
      window.chatLogs = []
    }
    window.chatLogs.push(sendLog)
    
    // 触发日志更新事件
    window.dispatchEvent(new CustomEvent('chatLogUpdate', {
      detail: sendLog
    }))
    
    // 清空输入框
    messageToSend.value = ''
  }).catch(err => {
    resultText.value = `❌ 发送失败: ${err}`
    console.error('发送失败:', err)
  })
}

onMounted(() => {
  // 恢复全局状态
  isConnected.value = window.connectionState.isConnected
  clientName.value = window.connectionState.clientName
  serverUrl.value = window.connectionState.serverUrl
  resultText.value = window.connectionState.resultText
  
  // 监听全局聊天日志更新事件，显示连接状态
  const handleChatLogUpdate = (event) => {
    console.log('🚀 ConnectPage 收到更新：', event.detail)
    if (isConnected.value) {
      connectionStatus.value = '已连接'
    }
  }
  
  window.addEventListener('chatLogUpdate', handleChatLogUpdate)
  window.handleChatLogUpdateConnect = handleChatLogUpdate
})

onUnmounted(() => {
  // 移除事件监听
  if (window.handleChatLogUpdateConnect) {
    window.removeEventListener('chatLogUpdate', window.handleChatLogUpdateConnect)
  }
})
</script>

<template>
  <div class="connect-page">
    <div class="header">
      <h1 class="title">🔌 客户端连接控制台</h1>
      <div class="status-badge" :class="{ active: isConnected }">
        <span class="status-dot"></span>
        <span>{{ isConnected ? '已连接' : '未连接' }}</span>
      </div>
      <div class="status-badge" :class="{ active: connectionStatus === '已连接' }">
        <span class="status-dot"></span>
        <span>{{ connectionStatus }}</span>
      </div>
    </div>

    <div class="content">
      <div class="result-box" :class="{ success: isConnected }">
        <pre class="result-text">{{ resultText }}</pre>
      </div>

      <div class="control-panel">
        <div class="input-group">
          <label>客户端名称:</label>
          <input 
            v-model="clientName" 
            :disabled="isConnected"
            placeholder="输入客户端名称"
            class="input-field"
          />
        </div>
        
        <div class="input-group">
          <label>服务器地址:</label>
          <input 
            v-model="serverUrl" 
            :disabled="isConnected"
            placeholder="例如: localhost"
            class="input-field"
          />
        </div>
        
        <div class="button-group">
          <button 
            class="btn btn-primary" 
            @click="handleConnect"
            :disabled="isConnected || isLoading"
          >
            <span v-if="isLoading">连接中...</span>
            <span v-else-if="isConnected">✅ 已连接</span>
            <span v-else>🔗 连接服务器</span>
          </button>
          
          <button 
            class="btn btn-danger" 
            @click="handleDisconnect"
            :disabled="!isConnected || isLoading"
          >
            <span>⛔ 断开连接</span>
          </button>
        </div>
      </div>
      
      <div class="send-panel" v-if="isConnected">
        <div class="input-group">
          <label>发送消息:</label>
          <div class="send-input-group">
            <input 
              v-model="messageToSend" 
              placeholder="输入要发送的消息"
              class="input-field"
              @keyup.enter="handleSend"
            />
            <button 
              class="btn btn-send" 
              @click="handleSend"
              :disabled="!messageToSend.trim()"
            >
              📤 发送
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.connect-page {
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

.control-panel, .send-panel {
  background: #fff;
  border: 3px solid #000;
  box-shadow: 8px 8px 0 #000;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-group label {
  font-weight: 700;
  font-size: 16px;
  color: #000;
  text-align: left;
}

.input-field {
  padding: 12px;
  border: 3px solid #000;
  font-size: 16px;
  font-family: 'Courier New', monospace;
}

.input-field:disabled {
  background: #f5f5f5;
  color: #999;
}

.button-group {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 16px;
}

.send-input-group {
  display: flex;
  gap: 12px;
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
  min-width: 150px;
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

.btn-send {
  background: #e8fff3;
  min-width: 100px;
}

.btn-send:hover:not(:disabled) {
  background: #c8e6c9;
}

.btn:disabled {
  background: #f5f5f5;
  color: #999;
  cursor: not-allowed;
  transform: none;
}
</style>