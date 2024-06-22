<template>
  <Path
      v-for="p in pathCoordinates"
      :editing=translating
      :start="p.start"
      :end="p.end"
      @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
  />

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
    v-for="(coord, i) in room.outline"
    :start="room.outline[i]"
    :end="room.outline[(i+1)%room.outline.length]"
    :editing=translating
    />
</template>

<style scoped>
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
  QualifiedKey,
  Repeats,
  RoomInterface, RotateEvent,
  RotationCoordinateInterface,
  RotationStartEvent,
} from "../../types.ts";
import {teamareaStore} from "../../stores/teamarea";
import {computed, inject, onMounted, provide, ref, watch} from "vue";
import Path from "./Path.vue";
import EditableTeamTable from "./EditableTeamTable.vue";
import Drag from "../../views/Drag.vue";

const settings = teamareaStore()

defineEmits(['dragStart', 'rotateStart', 'hoverElement', 'scrollElement'])
const pathCoordinates = computed(() => room.paths.map(pi => {
  const startCoord = keyToCoord(pi.start);
  const endCoord = keyToCoord(pi.end);
  const dx = startCoord.x - endCoord.x;
  const dy = startCoord.y - endCoord.y;
  const l2 = dx * dx + dy * dy;
  return {
    start: endCoord,
    end: startCoord,
    dx: dx,
    dy: dy,
    l2: l2,
    dist: Math.sqrt(l2)
  }
}));

function dist(a: CoordinateInterface, b: CoordinateInterface): number {
  const dx: number = b.x - a.x;
  const dy: number = b.y - a.y;
  return Math.sqrt(dx * dx + dy * dy);
}

const epsilon = 0.0001

function pointOnPath(p: CoordinateInterface): boolean {
  const threshold = 40;
  for (let i in pathCoordinates.value) {
    const path = pathCoordinates.value[i];
    if (path.dist < epsilon) {
      return dist(path.start, p) < threshold
    }

    const t = Math.max(0, Math.min(1, ((p.x - path.start.x) * path.dx + (p.y - path.start.y) * path.dy) / path.l2));

    // line.sx + t * dx, line.sy + t * dy
    const d = dist(p, {x: path.start.x + t * path.dx, y: path.start.y + t * path.dy})
    if (d < threshold) {
      return true
    }
  }

  return false;
}


function between(a: number, b: number, c: number): boolean {
  const eps = 3;
  return a - eps <= b && b <= c + eps;
}

function pathIntersection(a: PathCoordinatesInterface, b: PathCoordinatesInterface): CoordinateInterface | null {
  let x = ((a.start.x * a.end.y - a.start.y * a.end.x) * (b.start.x - b.end.x) - (a.start.x - a.end.x) * (b.start.x * b.end.y - b.start.y * b.end.x)) /
      ((a.start.x - a.end.x) * (b.start.y - b.end.y) - (a.start.y - a.end.y) * (b.start.x - b.end.x));
  let y = ((a.start.x * a.end.y - a.start.y * a.end.x) * (b.start.y - b.end.y) - (a.start.y - a.end.y) * (b.start.x * b.end.y - b.start.y * b.end.x)) /
      ((a.start.x - a.end.x) * (b.start.y - b.end.y) - (a.start.y - a.end.y) * (b.start.x - b.end.x));
  if (isNaN(x) || isNaN(y)) {
    return null;
  } else {
    if (a.start.x >= a.end.x) {
      if (!between(a.end.x, x, a.start.x)) {
        return null;
      }
    } else {
      if (!between(a.start.x, x, a.end.x)) {
        return null;
      }
    }
    if (a.start.y >= a.end.y) {
      if (!between(a.end.y, y, a.start.y)) {
        return null;
      }
    } else {
      if (!between(a.start.y, y, a.end.y)) {
        return null;
      }
    }
    if (b.start.x >= b.end.x) {
      if (!between(b.end.x, x, b.start.x)) {
        return null;
      }
    } else {
      if (!between(b.start.x, x, b.end.x)) {
        return null;
      }
    }
    if (b.start.y >= b.end.y) {
      if (!between(b.end.y, y, b.start.y)) {
        return null;
      }
    } else {
      if (!between(b.start.y, y, b.end.y)) {
        return null;
      }
    }
  }

  return {x: x, y: y};
}

function pathIntersect(b: PathCoordinatesInterface): CoordinateInterface | null {
  for (let i in pathCoordinates.value) {
    const coord = pathIntersection(pathCoordinates.value[i], b)
    if (coord !== null) {
      // console.log('intersect', toRaw(b.start), toRaw(b.end))
      return coord
    }
  }

  return null
}

provide("pointOnPath", pointOnPath)
provide("pathIntersect", pathIntersect)
provide("pathIntersection", pathIntersection)

const room = withDefaults(defineProps<RoomInterface & {
  translating?: boolean
  sequenceKey: QualifiedKey
}>(), {
  translating: false,
});

provide("editing", room.translating)
const elementFromQualifiedKey = inject<(s0: QualifiedKey, offset: number) => any>("elementFromQualifiedKey")!

function keyToCoord(sequenceKey: QualifiedKey): CoordinateInterface {
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

  const repeats = sequenceKey.slice(repeatStartsAt + 1)

  // Should not happen
  if (element.repeats.length > repeats.length) {
    console.log(innerKey)
    throw new Error(`Bounds check missed ${repeats.length} ${element.repeats.length}`)
  }

  // Derive the coordinate
  let baseVect: RotationCoordinateInterface = element.base
  for (let i = 0; i < element.repeats.length; i++) {
    baseVect = (element.repeats[i] as Repeats).calculateCoordinate(baseVect, repeats[i] as number)
  }

  if (element.repeats.length == repeats.length) {
    return baseVect
  }

  // Last part should be top or bottom
  const dy = repeats[element.repeats.length] == KeyCategory.Top ? -settings.PathAttachDistance : settings.PathAttachDistance
  return settings.offset(baseVect, 0, dy, true)
}

const roomRef = ref(null)

function deriveOutline() {

  const b = roomRef.value.getBBox()
  const padding = 0

  room.outline.splice(-room.outline.length)
  room.outline.push(
      {x: b.x - padding, y: b.y - padding},
      {x: b.x + b.width + padding, y: b.y - padding},
      {x: b.x + b.width + padding, y: b.y + b.height + padding},
      {x: b.x - padding, y: b.y + b.height + padding},
  )
}

</script>
