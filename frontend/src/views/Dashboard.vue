<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <div class="title">
        <div class="logo">🐟</div>
        <div>
          <h1>海洋牧场智能网箱监控系统</h1>
          <p class="subtitle">实时监测 · 数据采集 · 智能投喂</p>
        </div>
      </div>
      <div class="stats">
        <div class="stat-card">
          <span class="stat-label">网箱总数</span>
          <span class="stat-value">{{ cages.length }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-label">正常运行</span>
          <span class="stat-value value-normal">{{ normalCount }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-label">预警</span>
          <span class="stat-value value-warning">{{ warningCount }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-label">报警</span>
          <span class="stat-value value-alarm">{{ alarmCount }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-label">平均水温</span>
          <span class="stat-value">{{ avgTemp }} °C</span>
        </div>
      </div>
    </header>

    <div class="refresh-bar">
      <span class="refresh-info">
        最近更新: {{ lastUpdate || '--' }}
      </span>
      <button class="refresh-btn" @click="loadData" :disabled="loading">
        {{ loading ? '加载中...' : '刷新数据' }}
      </button>
    </div>

    <div class="main-content">
      <section class="panel viewer-panel">
        <div class="panel-header">
          <h2>3D 网箱总览</h2>
          <span class="tip">可拖拽旋转 · 滚轮缩放</span>
        </div>
        <div class="viewer-container">
          <CageViewer3D :cages="cages" />
        </div>
        <div class="legend">
          <span><i class="dot dot-normal"></i> 正常</span>
          <span><i class="dot dot-warning"></i> 预警</span>
          <span><i class="dot dot-alarm"></i> 报警</span>
        </div>
      </section>

      <section class="panel table-panel">
        <div class="panel-header">
          <h2>实时数据监测</h2>
        </div>
        <DataTable :cages="cages" @feed-success="onFeedSuccess" />
      </section>
    </div>

    <footer class="dashboard-footer">
      <p>© 2026 海洋牧场智慧养殖平台 · Powered by Vue3 + Three.js + Gin</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import CageViewer3D from '../components/CageViewer3D.vue'
import DataTable from '../components/DataTable.vue'
import { fetchCages } from '../api/cage.js'

const cages = ref([])
const loading = ref(false)
const lastUpdate = ref('')
let timer = null

const normalCount = computed(() => cages.value.filter((c) => c.status === 'normal').length)
const warningCount = computed(() => cages.value.filter((c) => c.status === 'warning').length)
const alarmCount = computed(() => cages.value.filter((c) => c.status === 'alarm').length)
const avgTemp = computed(() => {
  if (!cages.value.length) return '0.0'
  const sum = cages.value.reduce((s, c) => s + (c.water_temp || 0), 0)
  return (sum / cages.value.length).toFixed(1)
})

async function loadData() {
  loading.value = true
  try {
    const data = await fetchCages()
    cages.value = data || []
    lastUpdate.value = new Date().toLocaleString('zh-CN')
  } catch (err) {
    console.error('加载网箱数据失败:', err)
  } finally {
    loading.value = false
  }
}

function onFeedSuccess() {
  setTimeout(loadData, 800)
}

onMounted(() => {
  loadData()
  timer = setInterval(loadData, 30000)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #0a1628 0%, #0f2036 100%);
}

.dashboard-header {
  padding: 24px 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(0, 212, 170, 0.15);
  flex-wrap: wrap;
  gap: 20px;
}

.title {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo {
  font-size: 40px;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(0, 212, 170, 0.2), rgba(0, 168, 133, 0.1));
  border-radius: 12px;
  border: 1px solid rgba(0, 212, 170, 0.3);
}

.title h1 {
  margin: 0;
  font-size: 22px;
  color: #fff;
  font-weight: 700;
  letter-spacing: 1px;
}

.subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.stats {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.stat-card {
  background: rgba(0, 212, 170, 0.08);
  border: 1px solid rgba(0, 212, 170, 0.2);
  border-radius: 8px;
  padding: 10px 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 80px;
}

.stat-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.55);
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  margin-top: 2px;
}

.value-normal { color: #00d4aa; }
.value-warning { color: #ffaa00; }
.value-alarm { color: #ff4444; }

.refresh-bar {
  padding: 12px 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(0, 0, 0, 0.2);
}

.refresh-info {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.refresh-btn {
  padding: 6px 18px;
  background: transparent;
  border: 1px solid rgba(0, 212, 170, 0.4);
  border-radius: 4px;
  color: #00d4aa;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background: rgba(0, 212, 170, 0.15);
}

.refresh-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.main-content {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  padding: 20px 32px;
}

@media (max-width: 1100px) {
  .main-content { grid-template-columns: 1fr; }
}

.panel {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(0, 212, 170, 0.15);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  padding: 14px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(0, 212, 170, 0.12);
  background: rgba(0, 212, 170, 0.04);
}

.panel-header h2 {
  margin: 0;
  font-size: 15px;
  color: #00d4aa;
  font-weight: 600;
}

.tip {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.viewer-container {
  flex: 1;
  min-height: 420px;
  position: relative;
}

.viewer-panel { min-height: 480px; }

.legend {
  display: flex;
  gap: 18px;
  padding: 10px 20px;
  border-top: 1px solid rgba(0, 212, 170, 0.1);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.65);
}

.dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 6px;
  vertical-align: middle;
}

.dot-normal { background: #00d4aa; }
.dot-warning { background: #ffaa00; }
.dot-alarm { background: #ff4444; }

.table-panel {
  min-height: 480px;
  padding: 0;
}

.table-panel .panel-header + * {
  flex: 1;
  padding: 16px 20px;
  overflow: auto;
}

.dashboard-footer {
  padding: 16px 32px;
  text-align: center;
  color: rgba(255, 255, 255, 0.35);
  font-size: 12px;
  border-top: 1px solid rgba(0, 212, 170, 0.08);
}
</style>
