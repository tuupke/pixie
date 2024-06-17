<template>
    <MultiSelect
        v-model="selectedBits"
        :options="options"
        optionLabel="name"
        placeholder="Select options"
        class="w-full md:w-20rem"
        @change="selectHandle"
        display="chip"
        :showToggleAll="false"
        :optionDisabled="(el) => el.index == skipBit"
    />
</template>

<script setup lang="ts">

import {onMounted, ref} from "vue";
import {MultiSelectChangeEvent} from "primevue/multiselect";

const selectedBits = ref();
const options = ref();

function selectHandle(e: MultiSelectChangeEvent) {
  m.value = e.value.reduce((a: number, v: {index: number}) => a | (1 << v.index), 0)
}

const m = defineModel<number | null>();
const props = defineProps<{num: number, skipBit: number, options: string[]}>();

onMounted(() => {
  if (m.value === null) {
    m.value = 0;
  }

  let opts = [];
  for (let i in props.options) {
    opts.push({
      name: props.options[i],
      index: i,
    })
  }

  let mv: number = m.value!

  const val = []
  let i: number = 0;
  while (mv > 0) {
    if ((mv & 1) > 0) {
      val.push(opts[i])
    }

    i++;
    mv = mv >> 1;
  }

  selectedBits.value = val;
  options.value = opts;
});

</script>
