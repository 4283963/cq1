<template>
  <div class="data-table">
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>网箱名称</th>
          <th>水温 (°C)</th>
          <th>溶氧量 (mg/L)</th>
          <th>状态</th>
          <th>投喂量 (kg)</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="cage in cages" :key="cage.id">
          <td>{{ cage.id }}</td>
          <td>{{ cage.name }}</td>
          <td>
            <span :class="tempClass(cage.water_temp)">{{ cage.water_temp }}</span>
          </td>
          <td>
            <span :class="oxygenClass(cage.dissolved_oxygen)">{{ cage.dissolved_oxygen }}</span>
          </td>
          <td>
            <span :class="['status-badge', cage.status]">{{ statusLabel(cage.status) }}</span>
          </td>
          <td>
            <input
              type="number"
              v-model.number="feedAmounts[cage.id]"
              min="0.1"
              step="0.1"
              class="feed-input"
              placeholder="投喂量"
            />
          </td>
          <td>
            <button
              class="feed-btn"
              :disabled="feeding === cage.id"
              @click="handleFeed(cage.id)"
            >
              {{ feeding === cage.id ? '投喂中...' : '投喂' }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="feedMessage" class="feed-message" :class="{ success: feedSuccess, error: !feedSuccess }">
      {{ feedMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { sendFeedCommand } from '../api/cage.js'

const props = defineProps({
  cages: { type: Array, default: () => [] },
})

const feedAmounts = reactive({})
const feeding = ref(null)
const feedMessage = ref('')
const feedSuccess = ref(true)

function tempClass(temp) {
  if (temp > 30) return 'value-alarm'
  if (temp > 26) return 'value-warning'
  return 'value-normal'
}

function oxygenClass(oxygen) {
  if (oxygen < 5) return 'value-alarm'
  if (oxygen < 7) return 'value-warning'
  return 'value-normal'
}

function statusLabel(status) {
  const map = { normal: '正常', warning: '预警', alarm: '报警' }
  return map[status] || status
}

async function handleFeed(cageId) {
  const amount = feedAmounts[cageId]
  if (!amount || amount <= 0) {
    feedMessage.value = '请输入有效的投喂量'
    feedSuccess.value = false
    return
  }

  feeding.value = cageId
  feedMessage.value = ''

  try {
    const res = await sendFeedCommand(cageId, amount)
    feedMessage.value = res.message || '投喂指令已发送'
    feedSuccess.value = true
    feedAmounts[cageId] = null
  } catch (err) {
    feedMessage.value = '投喂失败: ' + (err.message || '未知错误')
    feedSuccess.value = false
  } finally {
    feeding.value = null
  }
}
</script>

<style scoped>
.data-table {
  width: 100%;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

thead th {
  background: rgba(0, 212, 170, 0.1);
  color: #00d4aa;
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  border-bottom: 1px solid rgba(0, 212, 170, 0.2);
  white-space: nowrap;
}

tbody td {
  padding: 10px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  color: #c0d0e0;
}

tbody tr:hover {
  background: rgba(0, 212, 170, 0.05);
}

.value-normal { color: #00d4aa; }
.value-warning { color: #ffaa00; }
.value-alarm { color: #ff4444; font-weight: bold; }

.status-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}
.status-badge.normal { background: rgba(0, 212, 170, 0.15); color: #00d4aa; }
.status-badge.warning { background: rgba(255, 170, 0, 0.15); color: #ffaa00; }
.status-badge.alarm { background: rgba(255, 68, 68, 0.15); color: #ff4444; }

.feed-input {
  width: 80px;
  padding: 4px 8px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(0, 212, 170, 0.3);
  border-radius: 4px;
  color: #c0d0e0;
  font-size: 13px;
  outline: none;
}
.feed-input:focus {
  border-color: #00d4aa;
}

.feed-btn {
  padding: 5px 16px;
  background: linear-gradient(135deg, #00d4aa, #00a885);
  border: none;
  border-radius: 4px;
  color: #0a1628;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}
.feed-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 212, 170, 0.4);
}
.feed-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.feed-message {
  margin-top: 12px;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 13px;
}
.feed-message.success { background: rgba(0, 212, 170, 0.15); color: #00d4aa; }
.feed-message.error { background: rgba(255, 68, 68, 0.15); color: #ff4444; }
</style>
