import { ref, watch, onBeforeUnmount } from 'vue'
import { sendFeedCommand } from '../api/cage.js'

const OXYGEN_THRESHOLD = 5.0
const SCHEDULED_HOUR = 16
const SCHEDULED_MINUTE = 0
const COOLDOWN_MS = 10 * 60 * 1000
const DEFAULT_FEED_AMOUNT = 2.0
const SCHEDULED_FEED_AMOUNT = 3.0
const TIMER_INTERVAL_MS = 30 * 1000

export function useAutoFeed(cages) {
  const enabled = ref(false)
  const autoFeedLog = ref([])
  const feedingCageIds = ref(new Set())

  const lastFeedTime = new Map()
  let scheduleTimer = null

  function canFeed(cageId) {
    const last = lastFeedTime.get(cageId)
    if (!last) return true
    return Date.now() - last >= COOLDOWN_MS
  }

  function addLog(cageId, cageName, reason, amount, success) {
    autoFeedLog.value.unshift({
      time: new Date().toLocaleString('zh-CN'),
      cageId,
      cageName,
      reason,
      amount,
      success,
    })
    if (autoFeedLog.value.length > 20) {
      autoFeedLog.value = autoFeedLog.value.slice(0, 20)
    }
  }

  async function triggerFeed(cageId, cageName, reason, amount) {
    if (feedingCageIds.value.has(cageId)) return
    if (!canFeed(cageId)) return

    feedingCageIds.value.add(cageId)
    try {
      await sendFeedCommand(cageId, amount, 'auto_feed')
      lastFeedTime.set(cageId, Date.now())
      addLog(cageId, cageName, reason, amount, true)
    } catch (err) {
      addLog(cageId, cageName, reason, amount, false)
    } finally {
      feedingCageIds.value.delete(cageId)
    }
  }

  function checkLowOxygen(cagesData) {
    if (!enabled.value) return
    for (const cage of cagesData) {
      if (cage.dissolved_oxygen < OXYGEN_THRESHOLD) {
        triggerFeed(
          cage.id,
          cage.name,
          `溶氧量过低 (${cage.dissolved_oxygen.toFixed(1)} mg/L)`,
          DEFAULT_FEED_AMOUNT
        )
      }
    }
  }

  function checkScheduledFeed() {
    if (!enabled.value) return
    const now = new Date()
    if (now.getHours() === SCHEDULED_HOUR && now.getMinutes() === SCHEDULED_MINUTE) {
      for (const cage of cages.value) {
        triggerFeed(
          cage.id,
          cage.name,
          `定时投喂 (${SCHEDULED_HOUR}:${String(SCHEDULED_MINUTE).padStart(2, '0')})`,
          SCHEDULED_FEED_AMOUNT
        )
      }
    }
  }

  watch(
    () => cages.value,
    (newCages) => {
      if (newCages && newCages.length) {
        checkLowOxygen(newCages)
      }
    },
    { deep: true }
  )

  function startScheduler() {
    stopScheduler()
    scheduleTimer = setInterval(checkScheduledFeed, TIMER_INTERVAL_MS)
  }

  function stopScheduler() {
    if (scheduleTimer) {
      clearInterval(scheduleTimer)
      scheduleTimer = null
    }
  }

  watch(enabled, (val) => {
    if (val) {
      startScheduler()
      if (cages.value && cages.value.length) {
        checkLowOxygen(cages.value)
      }
    } else {
      stopScheduler()
    }
  })

  onBeforeUnmount(() => {
    stopScheduler()
  })

  return {
    enabled,
    autoFeedLog,
    OXYGEN_THRESHOLD,
    SCHEDULED_HOUR,
    SCHEDULED_MINUTE,
  }
}
