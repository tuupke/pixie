<template>
  <Sequence
      v-if="el.repeats.length -1 > atRepeats"
      v-for="(coord, i) in coords"

      @hoverElement="(e: ElementEvent) => $emit('hoverElement', e)"

      :x=coord.x
      :y=coord.y
      :rotation=coord.rotation
      :sequence-key="atRepeats < 0 ? sequenceKey.concat(KeyCategory.Repeats) : sequenceKey.concat(i)"

      :atRepeats=atRepeats+1
      :el=el
      :room=room
  />
  <EditableTeamTable
      v-else
      v-for="(c, i) in coords"

      @hoverElement="(e: ElementEvent) => $emit('hoverElement', e)"

      v-bind="c"
      :sequence-key="sequenceKey.concat(i)"
  />
</template>

<script setup lang="ts">

import {computed, inject} from "vue";
import {
  ElementEvent,
  ElementInterface,
  KeyCategory,
  QualifiedKey,
  RoomInterface,
  RotationCoordinateInterface} from "../../types.ts";
import EditableTeamTable from "./EditableTeamTable.vue";

interface SequenceLocal {
  atRepeats: number
  el: ElementInterface,
  room: RoomInterface,
  sequenceKey: QualifiedKey,
}

defineEmits(['hoverElement']);

const props = withDefaults(defineProps<RotationCoordinateInterface & SequenceLocal>(), {
  x: 0,
  y: 0,
  rotation: 0,
});

const elementFromQualifiedKey = inject<(s0: QualifiedKey) => any>("elementFromQualifiedKey")!

const coords = computed(() => {
  const repeats = props.atRepeats
  const index = props.sequenceKey.indexOf(KeyCategory.Repeats)
  if (repeats < 0 || index < 0) {
    return [{x: 0, y: 0, rotation: 0}]
  }

  const simplified = props.sequenceKey.slice(0, index+1)
  const repeatBase = elementFromQualifiedKey(simplified)
  if (repeatBase === null || repeatBase.length <= props.atRepeats) {
    console.log("Repeatbase not found")
    return [{x: 0, y: 0, rotation: 0}]
  }

  return repeatBase[repeats].calculateAllCoordinates(props)
})

</script>

<style scoped>
</style>
