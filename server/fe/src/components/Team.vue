<template>
  <div class="p-inputgroup">
    <Dropdown v-bind="modelValue" :options="options" :filter="true" optionValue='guid' :filterFields="filterFields"  @change="change">
      <template #option="slotProps">
        <slot name="option" v-bind="slotProps.option"></slot>
      </template>
      <template #value="slotProps">
        <div v-if="slotProps.value" :set="res = options.find(e => e.guid === slotProps.value)">
          <slot name="value" v-bind="res"></slot>
        </div>
        <span v-else>
          Nothing selected
        </span>
      </template>
    </Dropdown>
  </div>
</template>

<script>
export default {
  props: {
    'modelValue': {type: String, required: false},
    'options': {type: Array, required: true},
    'filterFields': {type: Array, required: true},
    // 'nameField': {type: String, required: true},
  },
  methods: {
    change($event) {
      this.$emit('onchange', $event)
      this.$emit('update:modelValue', $event.value)
    }
  },

  emits: ['update:modelValue', 'onchange']
}

</script>