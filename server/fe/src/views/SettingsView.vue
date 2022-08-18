<template>
  <div v-for='setting in store.settings' :key='setting.guid' class="col-8">
    <h5>{{ setting.key }}</h5>
    <div class="p-inputgroup">
      <Button label="Save" @click="store.save(setting.guid)"/>
      <InputText type="text" v-model="setting.value"/>
    </div>
  </div>
  <div class="col-8">
    <Team v-model="team"/>
  </div>

</template>

<script setup>

import Team from '@/components/Team.vue'

import {settingsStore} from "../stores/settings";
import {onMounted} from 'vue';

// declare store variable
const store = settingsStore();

onMounted(() => {
  if (!store.settings.length) {
    store.fetchSettings();
  }
})
</script>

<script>
export default {
  props: {'team': {type: String, required: true}},
}
</script>
