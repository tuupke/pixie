<template>
  <div v-if="assignment">
    <div class="mb-2">
      <SelectButton v-model="modeValue" :options="options"/>
    </div>
    <div class="mb-2">
      <InputNumber v-if="modeValue === 'Automatic'" disabled placeholder="Keyword" v-model="assignment.num"/>
      <InputNumber v-else-if="modeValue === 'Manual'" placeholder="Table number" v-model="number"/>
      <InputNumber v-else disabled placeholder=""/>
    </div>
  </div>
  <div v-else>
    Key not found
  </div>
</template>

<script setup lang="ts">
import {KeyCategory, QualifiedKey} from "../types.ts";
import {computed, toRaw} from "vue";
import {mapStore, tableAssignment} from "../stores/map.ts";

const options = ['Automatic', 'Manual', 'Ignored'];

const modeValue = computed({
  get: () => {
    if (exception.value === undefined) {
      return 'Automatic'
    }

    if (exception.value === null) {
      return 'Ignored'
    }

    if (typeof exception.value === "number") {
      return 'Manual';
    }

    return 'Unknown';
  },
  set: (newValue: string): void => setOveride(newValue, number.value)
})

const map = mapStore()
const props = defineProps<{ sequenceKey: QualifiedKey }>()

const assignment = computed<tableAssignment | undefined>(() => map.assignments.getValue(props.sequenceKey))
const number = computed<number>({
  get: () => !assignment.value || assignment.value.ignored ? 1 : assignment.value.num,
  set: (newValue: number): void => setOveride(modeValue.value, newValue)
})

const exception = computed(() => map.exceptions.getValue(props.sequenceKey))

function setOveride(type: string, override: number | null) {
  const room = map.fromQualifiedKeyUpToIncluding(props.sequenceKey, KeyCategory.Room)

  const suffixIndex = props.sequenceKey.indexOf(KeyCategory.Elements)
  const suffix = props.sequenceKey.slice(suffixIndex)

  // Assume that the value is at the end of the exception list
  for (let i =  room.overrides.length -1; i >= 0; i--) {
    if (map.keyCompare(toRaw(room.overrides[i][0]), suffix) === 0) {
      room.overrides.splice(i, 1)
    }
  }

  if (type === 'Automatic') {
    return
  }

  let newOverride = null;
  if (type === 'Manual') {
    newOverride = override;
  }

  room.overrides.push([suffix, newOverride])
}


</script>
