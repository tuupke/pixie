<template>
  <div class="grid">
    <div class="bg">
      r1c1 aaaaa
    </div>
    <div class="bg">
      r1c2
    </div>

    <div class="bg">
      r2c1
    </div>
    <div class="bg">
      r2c2
    </div>
  </div>
</template>

<script setup>
import {teamareaStore} from "@/stores/teamarea";

const settings = teamareaStore()

const scale = 3
</script>

<style>
.grid {
  display: inline-grid;
  grid-template-columns: auto auto;
}

.bg {
  background: red;
  margin: 3px;
}

.full {
  flex-basis: 100%;
}
</style>
