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
  <div class="grid col-12 flex flex-wrap card-container blue-container" >
    <Toast />
    <RouterView/>
  </div>
</template>

<style scoped>

</style>
