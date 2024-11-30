<template>
  <PlacedCrossHairs
      :highlightable="false"
      v-if="allowsSplit && hoverCoord !== null"

      :x="hoverCoord.x"
      :y="hoverCoord.y"
      :scale-multiply="1.5*scaleMultiply"
  />
  <line :x1="start.x" :y1="start.y" :x2="drawnStart.x" :y2="drawnStart.y"
        stroke="lightGray" :stroke-width="1/scale*size" stroke-linecap="round"/>
  <line :x1="end.x" :y1="end.y" :x2="drawnEnd.x" :y2="drawnEnd.y"
        stroke="lightGray" :stroke-width="1/scale*size" stroke-linecap="round"/>

  <line :x1="drawnStart.x"
        :y1="drawnStart.y"
        :x2="drawnEnd.x"
        :y2="drawnEnd.y"

        @mousemove="calcCoord"
        @mouseout="hoverCoord=null"
        @mousedown="(e) => {if (props.allowsSplit && hoverCoord !== null) e.stopPropagation()}"
        @click="emitCoord"
        stroke="gray" :stroke-width="1/scale*size" stroke-linecap="round" />

  <Drag v-if="editing" :coord="start"
        @dragStart="e => $emit('dragStart', e)">
    <PlacedCrossHairs
        v-bind="start"
        @click.alt.stop="$emit('deleteCoord', start)"
        :highlightable="true"
        :scale-multiply="scaleMultiply"
        :color="color"
    />
  </Drag>
    <Drag v-if="editing" :coord="end"
          @dragStart="e => $emit('dragStart', e)">
      <PlacedCrossHairs
          v-bind="end"
          @click.alt.stop="$emit('deleteCoord', end)"
          :highlightable="true"
          :scale-multiply="scaleMultiply"
          :color="color"
      />
  </Drag>

</template>

<script lang="ts" setup>

import {CoordinateInterface, DragStartEvent, PathCoordinatesInterface, Vector} from "../../types.ts";
import PlacedCrossHairs from "./PlacedCrossHairs.vue";
import Drag from "../../views/Drag.vue";
import {computed, inject, ref} from "vue";

const size = 5;

const props = withDefaults(defineProps<{
  start: CoordinateInterface,
  end: CoordinateInterface,
  drawStart?: CoordinateInterface,
  drawEnd?: CoordinateInterface,
  editing?: boolean
  scale?: number
  scaleMultiply?: number
  color?: string,
  allowsSplit?: boolean,
}>(), {
  editing: false,
  scale: 1,
  scaleMultiply: 0.7,
  color: "green",
  allowsSplit: false,
});

const hoverCoord = ref<CoordinateInterface | null>(null)

const scale = inject<number>("scale", 1)!
const toInnerCoordinates = inject<(e: MouseEvent) => CoordinateInterface>('toInnerCoordinates', () => { return {x: 0, y: 0}})!
const pathIntersection = inject<(a: PathCoordinatesInterface, b: PathCoordinatesInterface) => CoordinateInterface | null>("pathIntersection", () => null)

const drawnStart = computed(() => props.drawStart ?? props.start)
const drawnEnd = computed(() => props.drawEnd ?? props.end)

function calcCoord(e: MouseEvent) {
  if (!props.allowsSplit) {
    return
  }

  const val = toInnerCoordinates(e)

  const vect = new Vector(props.end.y - props.start.y, props.start.x - props.end.x).normalize().multiply(size*10)

  const otherPath: PathCoordinatesInterface = {
    start: vect.copy().add(val),
    end: vect.copy().multiply(-1).add(val),
  }

  hoverCoord.value = pathIntersection(props, otherPath)
}

const emits = defineEmits<{
  newCoord: [CoordinateInterface]
  deleteCoord: [CoordinateInterface]
  dragStart: [DragStartEvent]
}>()

function emitCoord() {
  if (!props.allowsSplit || !hoverCoord.value) {
    return
  }

  emits('newCoord', hoverCoord.value)
  hoverCoord.value = null
}

</script>
