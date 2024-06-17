<template>
  <Path
      v-if="topIntersect !== null && !assignment.ignored"
      :start="coordinate"
      :end="topIntersect"
  />
  <Path
      v-if="bottomIntersect !== null && !assignment.ignored"
      :start="coordinate"
      :end="bottomIntersect"
  />

  <TeamTable
      @mouseover.stop="(e: MouseEvent) => $emit('hoverElement', {key: sequenceKey, event: e})"
      :team-id="number"
      :x="x" :y="y" :rotation="rotation"
      :highlighted="highlighted"
      :fill="fill"
      :duplicated="assignment.duplicate"
      :hidden="!assignment.ignored"
      :path="attached && !assignment.ignored"
  />

  <PlacedCrossHairs
      v-if="(topIntersect !== null || showHandles) && !assignment.ignored && editing"
      v-bind="topIntersect !== null ? topIntersect : topHandle"

      highlightable
      color="blue"
  />

  <PlacedCrossHairs
      v-if="(bottomIntersect !== null || showHandles) && !assignment.ignored && editing"
      v-bind="bottomIntersect ? bottomIntersect : bottomHandle"
      highlightable
      :scale-multiply=0.9
      color="purple"
  />

</template>

<script setup lang="ts">

import {
  CoordinateInterface,
  KeyCategory,
  PathCoordinatesInterface,
  QualifiedKey,
  RotationCoordinateInterface,
  Vector
} from "../../types.ts";
import TeamTable from "./TeamTable.vue";
import PlacedCrossHairs from "./PlacedCrossHairs.vue";
import {teamareaStore} from "../../stores/teamarea.ts";
import {computed, inject, ref} from "vue";
import Path from "./Path.vue";
import {mapStore, tableAssignment} from "../../stores/map.ts";

const settings = teamareaStore()
const map = mapStore()

const coordinate = withDefaults(defineProps<RotationCoordinateInterface & {

  sequenceKey: QualifiedKey
}>(), {x: 0, y: 0, rotation: 0})

const editing = inject<boolean>("editing")! ?? false
defineEmits(['hoverElement'])

const pathIntersect = inject<(p: PathCoordinatesInterface) => CoordinateInterface>("pathIntersect")!
const highlightedKey = inject<ref<QualifiedKey>>("highlightedKey")!

const hatchingForColors = inject<(...c :string[]) => string>('hatchingForColors') ?? ((...c :string[]) => 'white')
const fill = computed(() => {
  // const hatchingColors = {blue: 'lightblue', gray: 'lightgray', green: 'lightgreen', orange: 'orange', white: 'white'}
  const colors = [];
  if (!attached.value) {

    colors.push('#013370')
  }

  if (assignment.value.duplicate) {
    colors.push('orange')
  }

  return hatchingForColors(...colors)
})

const highlighted = computed(() => map.keyCompare(coordinate.sequenceKey, highlightedKey.value) === 0)

const showHandles = false;
const topHandle = computed<CoordinateInterface>(() => settings.offset(coordinate, settings.PathAttachDistance, true))
const bottomHandle = computed<CoordinateInterface>(() => settings.offset(coordinate, settings.PathAttachDistance, false))

const topIntersect = computed<CoordinateInterface | null>(() => pathIntersect({
  start: coordinate,
  end: settings.offset(coordinate, settings.PathDetectDistance, true)
}))

const bottomIntersect = computed<CoordinateInterface | null>(() => pathIntersect({
  start: coordinate,
  end: settings.offset(coordinate, settings.PathDetectDistance, false)
}))

const attached = computed<boolean>(() => topIntersect.value !== null || bottomIntersect.value !== null)

const rotateableCoord = computed<CoordinateInterface>(() => {
  const index = coordinate.sequenceKey.indexOf(KeyCategory.Repeats)
  if (index >= 0) {
    const repeats = map.fromQualifiedKey(coordinate.sequenceKey, coordinate.sequenceKey.length - index - 1)
    for (let i = 0; i < repeats.length; i++) {
      // console.log(repeats[i].num-1, coordinate.sequenceKey[index+i+1])
      if (repeats[i].num - 1 !== coordinate.sequenceKey[index + i + 1]) {
        return null
      }
    }
  }

  const coords = [
    settings.cornerCoordinate(coordinate, true, false),
    settings.cornerCoordinate(coordinate, true, true),
    settings.cornerCoordinate(coordinate, false, false),
    settings.cornerCoordinate(coordinate, false, true),
  ];

  const coord = map.fromQualifiedKey(coordinate.sequenceKey, coordinate.sequenceKey.length - index).base
  const cvect = new Vector(coord.x, coord.y)

  let maxDist = 0;
  let maxCoord = null;

  for (let i = 0; i < coords.length; i++) {
    const ldist = cvect.distance(coords[i])
    if (ldist - 0.001 > maxDist) {
      maxDist = ldist
      maxCoord = coords[i]
    }
  }

  return maxCoord;
})

const assignment = computed<tableAssignment>(() => map.assignments.getValue(coordinate.sequenceKey) ?? {
  key: coordinate.sequenceKey,
  num: 0,
  ignored: true,
  duplicate: false,
})
const number = computed<string>(() => {
  return assignment.value.ignored ?
      '' : assignment.value?.num.toString();
});

</script>
