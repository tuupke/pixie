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
  <div class="">
    <TabMenu :model="items">
      <template #item="{ item, props }">
        <router-link v-slot="{ href, navigate }" :to="item.to" custom>
          <a v-ripple :href="href" v-bind="props.action" @click="navigate">
            <span v-bind="props.icon" />
            <span v-bind="props.label">{{ item.label }}</span>
          </a>
        </router-link>
      </template>
    </TabMenu>
<!--    <TabMenu :model="items"/>-->
    <Toast/>
    <RouterView/>
  </div>
</template>

<style scoped>

</style>
