<template>
  <Toolbar class="flex-initial flex m-2">
    <template #start>
      <ToggleButton on-label="Snapped" off-label="free-form" class="mr-2" v-model="snapToGrid"/>
      <Button
          label="Reset panzoom"
          class="mr-2"
          @click="rescale"/>
    </template>
    <template #end>
      <slot name="extra-buttons"/>
    </template>
  </Toolbar>
  <div class="settings-view-wrapper">
    <div class="col-3 left-nav">
      <slot name="settings"/>
    </div>
    <div class="col-9">
      <svg style="border: 1px solid red;" id="layoutsvg" ref="svgRef" width="100%" height="100%"
           @wheel.stop="scroll"
           @mousemove="maybeTranslateRotate"
           @mouseup="resetTranslateRotate"
           @mousedown="(e: MouseEvent) => dragStart({coord: topLeft, event: e}, true)">
        <g :transform="'scale('+scale+') translate('+topLeft.x+','+topLeft.y+')'">
          <rect fill="url(#pattern-circles)" :x="background.x" :y="background.y" :width="background.width"
                :height="background.height"/>
        </g>

        <g id="innerSvgRef" ref="innerSvgRef" :transform="'scale('+scale+') translate('+topLeft.x+','+topLeft.y+')'">
          <slot
              name="svg"
              :rotate="rotateElement"
              :moveStart="dragStart"
          />
        </g>

        <pattern
            id="pattern-circles"
            :x="-circleRadius+(((1-settings.areaOffsetX/100)*settings.areaWidth)%dotSpacing)"
            :y="-circleRadius+(((1-settings.areaOffsetY/100)*settings.areaHeight)%dotSpacing)"
            :width="dotSpacing"
            :height="dotSpacing"
            patternUnits="userSpaceOnUse"
            patternContentUnits="userSpaceOnUse">
          <circle id="pattern-circle" :cx="circleRadius" :cy="circleRadius" :r="circleRadius" fill="#000"></circle>
        </pattern>

        <Hatching
            v-for="colors in hatchings"
            :width="settings.areaWidth"
            :colors="colors"
            :id="colors.map((v) => v.replaceAll('#', '')).join('-')+'-hatching'"
        />
      </svg>
    </div>
  </div>
</template>

<style scoped>
svg rect {
  color: lightgray;
}

.settings-view-wrapper {
  display: flex;
  flex: 1;
}

</style>

<script setup lang="ts">

import Hatching from "../components/Layout/Hatching.vue";

import {computed, onMounted, provide, reactive, ref} from "vue";
import {
  CoordinateInterface,
  DragStartEvent,
  KeyCategory,
  RotateEvent,
  RotationCoordinateInterface,
  Vector
} from "../types.ts";
import {teamareaStore} from "../stores/teamarea";

import {mapStore} from "../stores/map";

const map = mapStore()

const settings = teamareaStore()

const topLeft = reactive<CoordinateInterface>({x: 0, y: 0});
const offset = ref<CoordinateInterface>({x: 0, y: 0});
const coordinateBeingTranslated = ref<CoordinateInterface | null>();

// -------------------- Scale --------------------
const scaleLower = 0.03
const scaleUpper = 2
const scale = ref<number>(1)
provide('scale', scale)

const svgRef = ref()
const innerSvgRef = ref()

const snapToGrid = ref(true);

// -------------------- Background --------------------
const circleRadius = computed(() => 0.9 / scale.value)
const rotateClamping = 15;
const translateClamping = 100;
defineExpose({rotateClamping, translateClamping, snapToGrid})


// Spacing calculation uses a heuristic that look acceptable.
const dotSpacing = computed(() => Math.max(Math.round(Math.abs(Math.log(scale.value) / Math.log(4))), 1) * translateClamping);
const background = computed(() => {
  if (!svgRef.value) {
    return {x: 0, y: 0, width: 0, height: 0}
  }

  const val = svgRef.value.getBBox()
  return {
    x: -topLeft.x,
    y: -topLeft.y,
    width: val.width / scale.value,
    height: val.height / scale.value,
  }
})

function scroll(state: WheelEvent) {
  // Before calculating scale, see where the cursor currently is. Needs to be kept 'constant'.
  const innerCoords = toInnerCoordinates(state)

  const dt = 0.02
  const scrollDt = 1 + (-Math.sign(state.deltaY) * dt)

  scale.value = Math.max(Math.min(scale.value * scrollDt, scaleUpper), scaleLower)

  // Calculate the new top-left offset
  const newInner = toInnerCoordinates(state)
  topLeft.x += newInner.x - innerCoords.x
  topLeft.y += newInner.y - innerCoords.y
}

function rotateElement(e: RotateEvent) {
  // Before calculating scale, see where the cursor currently is. Needs to be kept 'constant'.
  const mouse = toInnerCoordinates(e.event, true)
  const element = map.fromQualifiedKeyUpTo(e.key, KeyCategory.Elements);

  rotateAroundBy(element.base, mouse, e.event.deltaY)
}

const rotated = ref(0)

function rotateAroundBy(toRotate: RotationCoordinateInterface, around: CoordinateInterface, by: number) {
  by /= 10

  if (snapToGrid.value) {
    rotated.value += by

    if (Math.abs(rotated.value) < rotateClamping) {
      return;
    }

    by = rotated.value
    rotated.value = 0

    const newAngle = Math.round((toRotate.rotation + by) / rotateClamping) * rotateClamping
    by = newAngle - toRotate.rotation
  }

  const newCoord = new Vector(toRotate.x - around.x, toRotate.y - around.y).rotate(by).add(around)

  toRotate.x = newCoord.x
  toRotate.y = newCoord.y
  toRotate.rotation += by
}

const hatchings = ref<string[][]>([])
provide('hatchingForColors', hatchingForColors)

function hatchingForColors(...colors: string[]): string {
  if (colors.length === 0) {
    colors.push('white')
  }

  colors = colors.sort()
  const exists = hatchings.value.reduce((a, v) => a || (colors.length === v.length && colors.every((value, index) => value === v[index])), false)
  if (!exists) {
    // Append
    hatchings.value.push(colors)
  }
  return `url(#${colors.map((v) => v.replaceAll("#", "")).join("-")}-hatching)`
}

provide("elementFromQualifiedKey", map.fromQualifiedKey)

function resetTranslateRotate() {
  temporarilyPreventSnapping.value = false
  coordinateBeingTranslated.value = null
  offset.value = {x: 0, y: 0}
}

function maybeTranslateRotate(e: MouseEvent) {
  if (coordinateBeingTranslated.value === undefined || coordinateBeingTranslated.value == null) {
    return
  }

  const innerCoords = toInnerCoordinates(e)
  let newCoord = {
    x: innerCoords.x - offset.value.x,
    y: innerCoords.y - offset.value.y,
  }

  if (snapToGrid.value && !temporarilyPreventSnapping.value) {
    newCoord.x = Math.round(newCoord.x / translateClamping) * translateClamping
    newCoord.y = Math.round(newCoord.y / translateClamping) * translateClamping
  }

  coordinateBeingTranslated.value.x = newCoord.x
  coordinateBeingTranslated.value.y = newCoord.y
}

const temporarilyPreventSnapping = ref(false)
function dragStart(e: DragStartEvent, preventSnapping: boolean = false) {
  temporarilyPreventSnapping.value = preventSnapping
  const innerCoords = toInnerCoordinates(e.event)
  coordinateBeingTranslated.value = e.coord

  offset.value = {
    x: innerCoords.x - e.coord.x,
    y: innerCoords.y - e.coord.y,
  }
}

provide('toInnerCoordinates', toInnerCoordinates)

function toInnerCoordinates(e: MouseEvent, translate: boolean = false): CoordinateInterface {
  const br = svgRef.value.getBoundingClientRect();
  return {
    x: (e.clientX - br.left) / scale.value - (translate ? topLeft.x : 0),
    y: (e.clientY - br.top) / scale.value - (translate ? topLeft.y : 0),
  }
}

onMounted(rescale)

function rescale() {
  const padding = 40
  const contentBox = innerSvgRef.value.getBBox({
    stroke: true,
  })

  const svgBox = svgRef.value.getBoundingClientRect()

  const width = svgBox.width - 2 * padding
  const height = svgBox.height - 2 * padding

  const scaleX = width / contentBox.width
  const scaleY = height / contentBox.height

  topLeft.x = -contentBox.x
  topLeft.y = -contentBox.y

  // // Align the map to the center
  if (scaleX <= scaleY) {
    scale.value = scaleX;
    topLeft.x += padding / scaleX
    topLeft.y += (svgBox.height / scaleX - contentBox.height) / 2;
  } else {
    scale.value = scaleY;
    topLeft.x += (svgBox.width / scaleY - contentBox.width) / 2;
    topLeft.y += padding / scaleY
  }
}

</script>
