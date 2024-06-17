<template>
  <Button
      v-if="show"
      v-bind="$attrs"
      aria-label="confirm deletion"
      :outlined="!deleting"
      v-on:click.stop="maybeDeleteElement"
  />
<!--  <span-->
<!--      v-if="show"-->
<!--      v-bind="$attrs"-->
<!--      :class="$attrs.icon"-->
<!--      v-on:click.stop="maybeDeleteElement"></span>-->
</template>
<script setup lang="ts">
import {KeyCategory, QualifiedKey} from "../../types.ts";
import {computed, ref, toRaw} from "vue";
import {mapStore} from "../../stores/map.ts";

const props = defineProps<{ sequenceKey: QualifiedKey, upto?: KeyCategory }>()

const deleting = ref<boolean>(false)
const map = mapStore()

const show = computed<boolean>(() => props.sequenceKey && props.sequenceKey.length > 0 && (
    !props.upto || props.sequenceKey.indexOf(props.upto) > -1
));

function maybeDeleteElement() {
  if (!deleting.value) {
    deleting.value = true
    window.setTimeout(() => deleting.value = false, 1000)
    return
  }

  let key = toRaw(props.sequenceKey)
  if (props.upto) {
    key = key.splice(0, key.indexOf(props.upto)+2)
  }

  console.log("Deleting", key)
  map.doDelete(key)
}

</script>
