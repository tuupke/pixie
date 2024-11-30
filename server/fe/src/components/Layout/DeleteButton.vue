<template>
  <ConfirmButton
      v-if="show"
      v-bind="$attrs"
      aria-label="confirm deletion"
      v-on:confirmed="deleteElement"
  />
</template>
<script setup lang="ts">
import {KeyCategory, QualifiedKey} from "../../types.ts";
import {computed, toRaw} from "vue";
import {mapStore} from "../../stores/map.ts";
import ConfirmButton from "./ConfirmButton.vue";

const props = defineProps<{ sequenceKey: QualifiedKey, upto?: KeyCategory }>()

const map = mapStore()

const show = computed<boolean>(() => props.sequenceKey && props.sequenceKey.length > 0 && (
    !props.upto || props.sequenceKey.indexOf(props.upto) > -1
));

function deleteElement() {
  let key = toRaw(props.sequenceKey)
  if (props.upto) {
    key = key.splice(0, key.indexOf(props.upto)+2)
  }

  console.log("Deleting sequenceKey", props.sequenceKey)
  map.doDelete(key)
}

</script>
