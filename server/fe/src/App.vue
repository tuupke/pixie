<!--<script setup>-->
<!--import {RouterLink, RouterView} from 'vue-router'</script>-->

<script>
export default {
  data() {
    return {
      items: [
        {
          label: 'Home',
          icon: 'pi pi-fw pi-home',
          to: '/',
        },
        {
          label: 'Settings',
          to: 'settings',
          icon: 'pi pi-fw pi-cog'
        }
      ]
    }
  }
}

</script>

<template>
  <div class="grid">
    <div class="col-12 grid flex flex-wrap card-container blue-container" >
      <Toast />
      <RouterView/>
    </div>
  </div>
</template>

<style scoped>

</style>
