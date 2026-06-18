import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

export async function fetchCages() {
  const res = await api.get('/cages')
  return res.data.data
}

export async function fetchCageById(id) {
  const res = await api.get(`/cages/${id}`)
  return res.data.data
}

export async function sendFeedCommand(cageId, feedAmount, operator = 'web_user') {
  const res = await api.post('/feed', {
    cage_id: cageId,
    feed_amount: feedAmount,
    operator,
  })
  return res.data
}

export async function fetchFeedRecords() {
  const res = await api.get('/feed/records')
  return res.data.data
}

export async function fetchFeedRecordsByCage(cageId) {
  const res = await api.get(`/feed/records/${cageId}`)
  return res.data.data
}
