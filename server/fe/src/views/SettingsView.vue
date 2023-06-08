<template>
  <div class="col-12">
  <svg id="layoutsvg" width="100%" height="600px">
    <Room :x=0 :y=0 :width=610 :height=60 shape="rectangle"></Room>
    <TeamTable v-for="n in 10" :x="-teamtableStore.areaWidth/2+(10+teamtableStore.areaWidth)*n" :y="30" :team-id="140+n" :rotation="rotation"/>
<!--    <Room :x=80 :y=100 :width=80 :height=20 shape="circular"></Room>-->
  </svg>
  </div>

  <div  class="col-6">
    Area Width
    <Slider v-model="teamtableStore.areaWidth" />
  </div>
  <div  class="col-6">
    Area Height
    <Slider v-model="teamtableStore.areaHeight" />
  </div>
  <div  class="col-6">
    Seat Separation
    <Slider v-model="teamtableStore.seatSep" />
  </div>
  <div  class="col-6">
    Seat Distance
    <Slider v-model="teamtableStore.seatDist" />
  </div>
  <div  class="col-6">
    Margin-x
    <Slider v-model="teamtableStore.marginX" />
  </div>
  <div  class="col-6">
    Margin-y
    <Slider v-model="teamtableStore.marginY" />
  </div>

  <div  class="col-6">
    Num Seats
    <Slider v-model="teamtableStore.seatNum" />
  </div>
  <div  class="col-6">
    Seat height
    <Slider v-model="teamtableStore.seatHeight" />
  </div>

  <div  class="col-6">
    offsetX
    <Slider v-model="teamtableStore.offsetX" :min=-teamtableStore.marginX :max=teamtableStore.marginX />
  </div>
  <div  class="col-6">
    offsetY
    <Slider v-model="teamtableStore.offsetY" :min=-teamtableStore.marginY :max=teamtableStore.marginY />
  </div>
  <div  class="col-6">
    rotation
    <Slider v-model="rotation" :min=0 :max=360 />
  </div>
</template>

<script>
import {mapStores} from 'pinia'
import {teamtableStore} from "@/stores/teamtable";

export default {
  computed: mapStores(teamtableStore),
  props: {
    'rotation': {type: Number, required: false, default: 0},
  }
}

</script>
