<script setup>
import {ref} from "vue";
import {RouterView} from 'vue-router'

const items = ref([
  {
    label: "Pixie",
    icon: 'pi pi-fw pi-home',
    to: '/',
  },
  {
    label: 'Map Designer',
    icon: 'pi pi-fw pi-map',
    to: '/settings/map',
  },
  {
    label: 'Area Designer',
    icon: 'pi pi-fw pi-map',
    to: '/settings/area',
  },
  {
    label: 'Setup',
    to: '/settings/area',
    icon: 'pi pi-fw pi-cog'
  }
])
</script>

<template>
  <div class="flex-column">
    <TabMenu :model="items"/>
    <Toast/>
    <RouterView/>
  </div>
</template>

<style scoped>

</style>
