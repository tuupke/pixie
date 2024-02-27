<script setup lang="ts">

import {reactive, ref, watch} from "vue";
import {RouterView, useRoute} from 'vue-router'

const items = reactive<{
  label: string
  icon: string
  to: string
}[]>([
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

const route = useRoute();

const activeMenuItemIndex = ref(0);

watch(
    () => route.path,
    (updatedPath) => {
      activeMenuItemIndex.value = items.findIndex((item) => item.to.startsWith(updatedPath));
    },
    { immediate: true }
);
</script>

<template>
  <div class="app-wrapper">
    <TabMenu :model="items" :active-index="activeMenuItemIndex">
      <template #item="{ item, props }">
        <router-link v-slot="{ href, navigate }" :to="item.to" custom>
          <a :href="href" v-bind="props.action" @click="navigate">
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
.app-wrapper {
    display: flex;
    flex-direction: column;
    height: 100vh;
}
</style>
