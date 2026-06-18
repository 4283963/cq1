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
let cageMeshes = []

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

  buildCages(props.cages)

  window.addEventListener('resize', onResize)
  animate()
}

function buildCages(cages) {
  cageMeshes.forEach((m) => scene.remove(m))
  cageMeshes = []

  cages.forEach((cage) => {
    const group = new THREE.Group()

    const cageSize = 2
    const boxGeom = new THREE.BoxGeometry(cageSize, cageSize, cageSize)
    const edgesGeom = new THREE.EdgesGeometry(boxGeom)
    const edgeMat = new THREE.LineBasicMaterial({
      color: statusColors[cage.status] || 0x00d4aa,
    })
    const wireframe = new THREE.LineSegments(edgesGeom, edgeMat)
    group.add(wireframe)

    const faceMat = new THREE.MeshPhongMaterial({
      color: statusColors[cage.status] || 0x00d4aa,
      transparent: true,
      opacity: 0.15,
    })
    const faceMesh = new THREE.Mesh(boxGeom, faceMat)
    group.add(faceMesh)

    const pillarGeom = new THREE.CylinderGeometry(0.05, 0.05, cageSize + 1, 8)
    const pillarMat = new THREE.MeshPhongMaterial({ color: 0x667788 })
    const offsets = [
      [-cageSize / 2, 0, -cageSize / 2],
      [cageSize / 2, 0, -cageSize / 2],
      [-cageSize / 2, 0, cageSize / 2],
      [cageSize / 2, 0, cageSize / 2],
    ]
    offsets.forEach(([ox, oy, oz]) => {
      const pillar = new THREE.Mesh(pillarGeom, pillarMat)
      pillar.position.set(ox, oy, oz)
      group.add(pillar)
    })

    group.position.set(cage.location_x, cage.location_y, cage.location_z)
    group.userData = { cageId: cage.id }

    scene.add(group)
    cageMeshes.push(group)
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
  cageMeshes.forEach((group, i) => {
    group.position.y = Math.sin(time + i * 0.5) * 0.1
  })

  renderer.render(scene, camera)
}

watch(
  () => props.cages,
  (newCages) => {
    if (scene) buildCages(newCages)
  },
  { deep: true }
)

onMounted(() => {
  init()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  if (animationId) cancelAnimationFrame(animationId)
  if (renderer) {
    renderer.dispose()
    containerRef.value?.removeChild(renderer.domElement)
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
