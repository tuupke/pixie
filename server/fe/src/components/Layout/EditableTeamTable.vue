<template>
  <Path
      v-if="(showHandles || topIntersect !== null) && !assignment.ignored"
      :start="middleCoordinate"
      :end="topIntersect !== null ? topIntersect : topHandle"
  />
  <Path
      v-if="(showHandles || bottomIntersect !== null) && !assignment.ignored"
      :start="middleCoordinate"
      :end="bottomIntersect !== null ? bottomIntersect : bottomHandle"
  />

  <TeamTable
      @wheel.stop="(e: WheelEvent) => $emit('scrollElement', {key: sequenceKey, event: e})"
      @mouseover.stop="(e: MouseEvent) => $emit('hoverElement', {key: sequenceKey, event: e})"
      :team-id="number"
      :x="x" :y="y" :rotation="rotation"
      :highlighted="highlighted"
      :fill="fill"
      :duplicated="assignment.duplicate"
      :hidden="assignment.ignored"
      :path="attached && !assignment.ignored"
  />

  <PlacedCrossHairs
      v-if="showHandles && !assignment.ignored && editing"
      v-bind="topIntersect !== null ? topIntersect : topHandle"

      highlightable
      color="blue"
  />

  <PlacedCrossHairs
      v-if="showHandles && !assignment.ignored && editing"
      v-bind="bottomIntersect !== null ? bottomIntersect : bottomHandle"
      highlightable
      :scale-multiply=0.9
      color="purple"
  />

<!--  <PlacedCrossHairs v-bind="outline.northEast" color="orange"/>-->
<!--  <PlacedCrossHairs v-bind="outline.northWest" color="green"/>-->
<!--  <PlacedCrossHairs v-bind="outline.southWest" color="blue"/>-->
<!--  <PlacedCrossHairs v-bind="outline.southEast" color="red"/>-->
</template>

<script setup lang="ts">

import {
  CoordinateInterface,
  ElementEvent,
  KeyCategory,
  PathCoordinatesInterface,
  QualifiedKey,
  RotateEvent,
  RotationCoordinateInterface,
} from "../../types.ts";
import TeamTable from "./TeamTable.vue";
import PlacedCrossHairs from "./PlacedCrossHairs.vue";
import {teamareaStore} from "../../stores/teamarea.ts";
import {computed, inject, Ref, ref} from "vue";
import Path from "./Path.vue";
import {mapStore, tableAssignment} from "../../stores/map.ts";

const settings = teamareaStore()
const map = mapStore()

const coordinate = withDefaults(defineProps<RotationCoordinateInterface & {
  sequenceKey: QualifiedKey
}>(), {x: 0, y: 0, rotation: 0})

const middleCoordinate = computed(() => settings.offset(coordinate, 0, 0));

const editing = inject<Ref<boolean>>("editing", ref(false))!

defineEmits<{
  hoverElement: [ElementEvent]
  scrollElement: [RotateEvent]
}>()

const pathIntersect = inject<(p: PathCoordinatesInterface) => CoordinateInterface>("pathIntersect")!
const highlightedKey = inject<Ref<QualifiedKey>>("highlightedKey", ref([]))
const hatchingForColors = inject<(...c: string[]) => string>('hatchingForColors') ?? (() => 'white')

const outline = computed(() => settings.tableOutlineAtCoord(coordinate, topIntersect.value !== null, bottomIntersect.value !== null))

const fill = computed(() => {
  // const hatchingColors = {blue: 'lightblue', gray: 'lightgray', green: 'lightgreen', orange: 'orange', white: 'white'}
  const colors = [];
  if (!attached.value) {
    colors.push('#013370')
  }

  if (assignment.value.duplicate) {
    colors.push('orange')
  }

  if (colors.length === 0) {
    colors.push('white')
  }

  const encapsulated = [
    outline.value.northEast,
    outline.value.southEast,
    outline.value.northWest,
    outline.value.southWest
  ].reduce((carry: boolean, c: CoordinateInterface): boolean => carry && withinOutline(c, roomOutline.value), true)

  if (!encapsulated) {
    colors.push('gray')
  }

  return hatchingForColors(...colors)
})

const highlighted = computed(() => map.keyCompare(coordinate.sequenceKey, highlightedKey.value) === 0)

const showHandles = inject<boolean>("pathStart", false);
const topHandle = computed<CoordinateInterface>(() => settings.offset(coordinate, 0, -settings.PathAttachDistance, true))
const bottomHandle = computed<CoordinateInterface>(() => settings.offset(coordinate, 0, settings.PathAttachDistance, true))

const topIntersect = computed<CoordinateInterface | null>(() => pathIntersect({
  start: middleCoordinate.value,
  end: settings.offset(coordinate, 0, -settings.PathDetectDistance, true)
}))

const bottomIntersect = computed<CoordinateInterface | null>(() => pathIntersect({
  start: middleCoordinate.value,
  end: settings.offset(coordinate, 0, settings.PathDetectDistance, true)
}))

const attached = computed<boolean>(() => topIntersect.value !== null || bottomIntersect.value !== null)

const assignment = computed<tableAssignment>(() => map.assignments.getValue(coordinate.sequenceKey) ?? {
  key: coordinate.sequenceKey,
  num: 0,
  ignored: true,
  duplicate: false,
})

const number = computed<string>(() => assignment.value.ignored ? '' : assignment.value?.num.toString());
const roomOutline = computed(() => map.fromQualifiedKeyUpToIncluding(coordinate.sequenceKey, KeyCategory.Room).outline ?? [])

const pathIntersection = inject<(a: PathCoordinatesInterface, b: PathCoordinatesInterface) => CoordinateInterface | null>("pathIntersection", (): CoordinateInterface | null => null)

function withinOutline(coord: CoordinateInterface, outline: CoordinateInterface[]): boolean {
  if (outline.length < 2) {
    return true
  }

  const splitLine = {
    start: {x: coord.x, y: -Number.MAX_SAFE_INTEGER},
    end: {x: coord.x, y: coord.y},
  }

  return outline.reduce((carry: boolean, c: CoordinateInterface, i: number, arr: CoordinateInterface[]): boolean => {
    const next = arr[(i + 1) % arr.length]
    const pathSegment: PathCoordinatesInterface = {
      start: c,
      end: next,
    }

    const intersected: boolean = pathIntersection(splitLine, pathSegment) !== null
    return carry !== intersected
  }, false);
}


</script>
