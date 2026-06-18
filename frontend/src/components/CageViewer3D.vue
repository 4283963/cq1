<template>
  <div class="cage-viewer-3d" ref="containerRef"></div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import * as THREE from 'three'
import { OrbitControls } from 'three/addons/controls/OrbitControls.js'

const props = defineProps({
  cages: { type: Array, default: () => [] },
})

const emit = defineEmits(['select-cage'])

const containerRef = ref(null)

let scene, camera, renderer, controls, animationId
let cageMap = new Map()
let sharedGeometries = {}
let sharedMaterials = {}

const CAGE_SIZE = 2

const statusColors = {
  normal: 0x00d4aa,
  warning: 0xffaa00,
  alarm: 0xff4444,
}

function init() {
  const container = containerRef.value
  if (!container) return

  scene = new THREE.Scene()
  scene.background = new THREE.Color(0x0a1628)
  scene.fog = new THREE.FogExp2(0x0a1628, 0.03)

  camera = new THREE.PerspectiveCamera(
    60,
    container.clientWidth / container.clientHeight,
    0.1,
    1000
  )
  camera.position.set(10, 8, 12)

  renderer = new THREE.WebGLRenderer({ antialias: true })
  renderer.setSize(container.clientWidth, container.clientHeight)
  renderer.setPixelRatio(window.devicePixelRatio)
  container.appendChild(renderer.domElement)

  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05
  controls.target.set(0, 0, 2)

  const ambientLight = new THREE.AmbientLight(0x4488cc, 0.6)
  scene.add(ambientLight)

  const dirLight = new THREE.DirectionalLight(0xffffff, 0.8)
  dirLight.position.set(10, 15, 10)
  scene.add(dirLight)

  createSharedResources()
  createWater()

  updateCages(props.cages)

  window.addEventListener('resize', onResize)
  animate()
}

function createSharedResources() {
  sharedGeometries.box = new THREE.BoxGeometry(CAGE_SIZE, CAGE_SIZE, CAGE_SIZE)
  sharedGeometries.edges = new THREE.EdgesGeometry(sharedGeometries.box)
  sharedGeometries.pillar = new THREE.CylinderGeometry(0.05, 0.05, CAGE_SIZE + 1, 8)

  sharedMaterials.pillar = new THREE.MeshPhongMaterial({ color: 0x667788 })

  sharedMaterials.edge = {
    normal: new THREE.LineBasicMaterial({ color: statusColors.normal }),
    warning: new THREE.LineBasicMaterial({ color: statusColors.warning }),
    alarm: new THREE.LineBasicMaterial({ color: statusColors.alarm }),
  }

  sharedMaterials.face = {
    normal: new THREE.MeshPhongMaterial({
      color: statusColors.normal,
      transparent: true,
      opacity: 0.15,
    }),
    warning: new THREE.MeshPhongMaterial({
      color: statusColors.warning,
      transparent: true,
      opacity: 0.15,
    }),
    alarm: new THREE.MeshPhongMaterial({
      color: statusColors.alarm,
      transparent: true,
      opacity: 0.15,
    }),
  }
}

function createWater() {
  const waterGeometry = new THREE.PlaneGeometry(40, 40, 32, 32)
  const waterMaterial = new THREE.MeshPhongMaterial({
    color: 0x006994,
    transparent: true,
    opacity: 0.7,
    side: THREE.DoubleSide,
  })
  const water = new THREE.Mesh(waterGeometry, waterMaterial)
  water.rotation.x = -Math.PI / 2
  water.position.y = -0.5
  scene.add(water)
}

function createCageMesh(cage) {
  const group = new THREE.Group()

  const faceMesh = new THREE.Mesh(
    sharedGeometries.box,
    sharedMaterials.face[cage.status] || sharedMaterials.face.normal
  )
  group.add(faceMesh)

  const wireframe = new THREE.LineSegments(
    sharedGeometries.edges,
    sharedMaterials.edge[cage.status] || sharedMaterials.edge.normal
  )
  group.add(wireframe)

  const offsets = [
    [-CAGE_SIZE / 2, 0, -CAGE_SIZE / 2],
    [CAGE_SIZE / 2, 0, -CAGE_SIZE / 2],
    [-CAGE_SIZE / 2, 0, CAGE_SIZE / 2],
    [CAGE_SIZE / 2, 0, CAGE_SIZE / 2],
  ]
  const pillars = []
  offsets.forEach(([ox, oy, oz]) => {
    const pillar = new THREE.Mesh(sharedGeometries.pillar, sharedMaterials.pillar)
    pillar.position.set(ox, oy, oz)
    group.add(pillar)
    pillars.push(pillar)
  })

  group.position.set(cage.location_x, cage.location_y, cage.location_z)
  group.userData = { cageId: cage.id }

  scene.add(group)

  return {
    group,
    faceMesh,
    wireframe,
    pillars,
    phase: Math.random() * Math.PI * 2,
    baseY: cage.location_y,
    status: cage.status,
  }
}

function updateCageMesh(entry, cage) {
  entry.group.position.x = cage.location_x
  entry.group.position.z = cage.location_z
  entry.baseY = cage.location_y

  if (entry.status !== cage.status) {
    const status = cage.status || 'normal'
    entry.faceMesh.material = sharedMaterials.face[status] || sharedMaterials.face.normal
    entry.wireframe.material = sharedMaterials.edge[status] || sharedMaterials.edge.normal
    entry.status = status
  }
}

function disposeCageEntry(entry) {
  scene.remove(entry.group)
}

function updateCages(cages) {
  const incomingIds = new Set(cages.map((c) => c.id))

  for (const [id, entry] of cageMap) {
    if (!incomingIds.has(id)) {
      disposeCageEntry(entry)
      cageMap.delete(id)
    }
  }

  cages.forEach((cage) => {
    const entry = cageMap.get(cage.id)
    if (entry) {
      updateCageMesh(entry, cage)
    } else {
      cageMap.set(cage.id, createCageMesh(cage))
    }
  })
}

function onResize() {
  const container = containerRef.value
  if (!container || !camera || !renderer) return
  camera.aspect = container.clientWidth / container.clientHeight
  camera.updateProjectionMatrix()
  renderer.setSize(container.clientWidth, container.clientHeight)
}

function animate() {
  animationId = requestAnimationFrame(animate)
  controls.update()

  const time = Date.now() * 0.001
  let index = 0
  for (const entry of cageMap.values()) {
    const floatY = Math.sin(time * 1.2 + entry.phase) * 0.08
    entry.group.position.y = entry.baseY + floatY
    index++
  }

  renderer.render(scene, camera)
}

watch(
  () => props.cages,
  (newCages) => {
    if (scene) updateCages(newCages)
  },
  { deep: true }
)

function disposeAll() {
  for (const entry of cageMap.values()) {
    disposeCageEntry(entry)
  }
  cageMap.clear()

  Object.values(sharedGeometries).forEach((g) => g.dispose())
  Object.values(sharedMaterials).forEach((m) => {
    if (m && typeof m.dispose === 'function') {
      m.dispose()
    } else if (m && typeof m === 'object') {
      Object.values(m).forEach((mm) => mm.dispose())
    }
  })
  sharedGeometries = {}
  sharedMaterials = {}
}

onMounted(() => {
  init()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  if (animationId) cancelAnimationFrame(animationId)
  disposeAll()
  if (renderer) {
    renderer.dispose()
    if (renderer.domElement && containerRef.value) {
      containerRef.value.removeChild(renderer.domElement)
    }
  }
  if (scene) {
    scene.clear()
  }
})
</script>

<style scoped>
.cage-viewer-3d {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}
</style>
