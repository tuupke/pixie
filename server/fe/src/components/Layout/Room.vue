<template>
  <g ref="roomRef">
    <Drag class="dragProp"
          v-for="(el, i) in room.elements as ElementInterface[]"
          :coord="el.base"

          @hoverElement="(e: ElementEvent) => $emit('hoverElement', e)"
          @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
          @scrollElement="(e: RotateEvent) => $emit('scrollElement', e)">
      <Sequence
          v-if="el.repeats.length>0"

          @hoverElement="(e: ElementEvent) => $emit('hoverElement', e)"
          @scrollElement="(e: RotateEvent) => $emit('scrollElement', e)"

          :x=el.base.x
          :y=el.base.y
          :rotation=el.base.rotation
          :sequence-key="sequenceKey.concat(KeyCategory.Room, KeyCategory.Elements, i, KeyCategory.Repeats)"
          :atRepeats=0
          :el=el
          :room=room
      />
      <EditableTeamTable
          v-else

          @hoverElement="(e: ElementEvent) => $emit('hoverElement', e)"
          @scrollElement="(e: RotateEvent) => $emit('scrollElement', e)"

          :x=el.base.x
          :y=el.base.y
          :rotation=el.base.rotation
          team-id="12"
          :sequence-key="sequenceKey.concat(KeyCategory.Room, KeyCategory.Elements, i)"
          :editing="translating"/>
    </Drag>
  </g>

  <Path
      v-for="p in pathCoordinates"
      :editing=translating
      :start="p.start"
      :end="p.end"
      :drawStart="p.drawStart ?? p.end"
      :drawEnd="p.drawEnd ?? p.start"
      @newCoord="c => {console.log(c)}"
      @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
  />

  <Path
      v-for="i in room.outline.length"
      :start="room.outline[i-1]"
      :end="room.outline[(i)%room.outline.length]"
      :editing=translating
      :allowsSplit="pathStart"
      @deleteCoord="c => deleteOutlineAt(c)"
      @newCoord="c => room.outline.splice(i, 0, c)"
      @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
  />

</template>

<style scoped>
* {
  cursor: crosshair;
}
</style>

<script setup lang="ts">

import Sequence from "./Sequence.vue";
import {
  CoordinateInterface,
  DragStartEvent,
  ElementEvent,
  ElementInterface,
  KeyCategory,
  PathCoordinatesInterface,
  PathInterface,
  QualifiedKey,
  Repeats,
  RoomInterface,
  RotateEvent,
  RotationCoordinateInterface,
} from "../../types.ts";
import {teamareaStore} from "../../stores/teamarea";
import {computed, inject, provide, Ref, ref} from "vue";
import Path from "./Path.vue";
import EditableTeamTable from "./EditableTeamTable.vue";
import Drag from "../../views/Drag.vue";
import {dist, pathInterface, pathIntersect, pathIntersection, pointOnPath} from "../../path_math.ts";

const settings = teamareaStore()
const pathStart = inject<Ref<boolean>>('canAdd', ref<boolean>(false))

function isArbitrary(pi: PathInterface): boolean {
  return pi.start[0] === KeyCategory.Arbitrary || pi.end[0] === KeyCategory.Arbitrary
}

defineEmits(['dragStart', 'rotateStart', 'hoverElement', 'scrollElement'])
const pathCoordinates = computed(() => room.paths.sort((a, b) => {
  const aIsArbitrary = isArbitrary(a)
  const bIsArbitrary = isArbitrary(b)
  if (aIsArbitrary && bIsArbitrary) {
    return 0
  } else if (aIsArbitrary) {
    return 1
  } else {
    return -1
  }
}).map((pi: PathInterface) => {
  const startCoord = keyToCoord(pi.start);
  const endCoord = keyToCoord(pi.end);

  if (startCoord === undefined || endCoord === undefined) {
    keyToCoord(pi.start)
  }

  const dx = startCoord.x - endCoord.x;
  const dy = startCoord.y - endCoord.y;
  const l2 = dx * dx + dy * dy;
  const startArbitrary = pi.start[0] === KeyCategory.Arbitrary
  const endArbitrary = pi.end[0] === KeyCategory.Arbitrary
  return {
    start: startCoord,
    end: endCoord,
    drawStart: startArbitrary? null : startCoord,
    drawEnd: endArbitrary? null : endCoord,
    startArbitrary: startArbitrary,
    endArbitrary: endArbitrary,
    startDist: startArbitrary ? l2 : 0,
    endDist: endArbitrary ? l2 : 0,
    dx: dx,
    dy: dy,
    l2: l2,
    dist: Math.sqrt(l2)
  } as pathInterface
}).map((pi: pathInterface, index: number, paths: pathInterface[]) => {
  if (!pi.startArbitrary && !pi.endArbitrary) {
    return pi
  }

  for (let i = 0; i < index; i++) {
    const intersection = pathIntersection(pi, paths[i])
    if (intersection === null) {
      continue
    }

    const handle = (pi: pathInterface, intersection: CoordinateInterface): pathInterface => {
      if (pi.startArbitrary) {
        const startDistNew = dist(pi.start, intersection)
        if (startDistNew < pi.startDist) {
          pi.startDist = startDistNew
          pi.drawStart = intersection
        }
      }

      if (pi.endArbitrary) {
        const endDistNew = dist(pi.end, intersection)
        if (endDistNew < pi.endDist) {
          pi.endDist = endDistNew
          pi.drawEnd = intersection
        }
      }

      return pi
    };

    pi = handle(pi, intersection)
    paths[i] = handle(paths[i], intersection)
  }

  return pi
}).reverse())

provide("pointOnPath", pointOnPath)
provide("pathIntersect", pathIntersect)
provide("pathIntersection", pathIntersection)

const room = withDefaults(defineProps<RoomInterface & {
  translating?: boolean
  coord?: RotationCoordinateInterface
  sequenceKey: QualifiedKey
}>(), {
  translating: false,
  coord: {x: 0, y: 0, rotation: 0},
});

provide("editing", room.translating)
const elementFromQualifiedKey = inject<(s0: QualifiedKey, offset: number) => any>("elementFromQualifiedKey")!

function keyToCoord(sequenceKey: QualifiedKey): CoordinateInterface {
  if (sequenceKey[0] === KeyCategory.Arbitrary) {
    return {x: sequenceKey[1] as number, y: sequenceKey[2] as number}
  }

  let innerKey = [...sequenceKey];
  if (innerKey[0] != KeyCategory.Placement) {
    innerKey = [...room.sequenceKey, KeyCategory.Room, ...sequenceKey];
  }

  const index = innerKey.indexOf(KeyCategory.Elements)
  if (index < 0 || innerKey.length < index + 2) {
    throw new Error("incorrect key")
  }

  // Find element
  const element = elementFromQualifiedKey(innerKey, innerKey.length - index - 2)
  if (element === null) {
    throw new Error("incorrect element, cannot be found")
  }

  const repeatStartsAt = sequenceKey.indexOf(KeyCategory.Repeats)
  if (repeatStartsAt < 0) {
    throw new Error("incorrect element, has no repeats")
  }

  const lastKey = sequenceKey[sequenceKey.length - 1]
  const topOffset = [KeyCategory.Top, KeyCategory.Bottom].indexOf(lastKey) >= 0 ? -1 : undefined

  const repeats = sequenceKey.slice(repeatStartsAt + 1, topOffset)

  // Derive the coordinate
  let baseVect: RotationCoordinateInterface = element.base
  for (let i = 0; i < Math.min(element.repeats.length, repeats.length); i++) {
    baseVect = (element.repeats[i] as Repeats).calculateCoordinate(baseVect, repeats[i] as number)
  }

  if (topOffset === undefined) {
    return baseVect
  }

  // Last part should be top or bottom
  const dy = lastKey == KeyCategory.Top ? -settings.PathAttachDistance : settings.PathAttachDistance
  return settings.offset(baseVect, 0, dy, true)
}

const roomRef = ref(null)

defineExpose({deriveOutline})

function deriveOutline(padding: number = 0, snapTo: number = 0) {
  if (roomRef.value === null) {
    return []
  }

  const b = roomRef.value.getBBox()

  let minx = b.x - padding
  let maxx = b.x + b.width + padding
  let miny = b.y - padding
  let maxy = b.y + b.height + padding

  // Skip snapping for low values, '3' arbitrary.
  if (snapTo > 3) {
    minx = Math.floor(minx / snapTo) * snapTo
    maxx = Math.ceil(maxx / snapTo) * snapTo
    miny = Math.floor(miny / snapTo) * snapTo
    maxy = Math.ceil(maxy / snapTo) * snapTo
  }

  return [
    {x: minx, y: miny},
    {x: maxx, y: miny},
    {x: maxx, y: maxy},
    {x: minx, y: maxy},
  ]
}

function deleteOutlineAt(coord: CoordinateInterface) {
  room.outline.splice(room.outline.indexOf(coord), 1)
  if (room.outline.length <= 2) {
    // Remove the entire outline, it has now become a line.
    room.outline.splice(0, room.outline.length)
  }
}

</script>
