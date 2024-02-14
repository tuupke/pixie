<template>
  <!--  style="display: flex;flex-direction: column;justify-content: space-between;"-->
  <div class="col-9 h-screen flex flex-column" ref="map">
    <svg id="layoutsvg" width="100%" height="100%" @mousemove="maybeTranslateRotate"  @mouseup="resetTranslateRotate">
      <!--      <TeamTable :x="53" :y="50" :rotation="rot" team-id="100"/>-->

      <g :transform="'scale('+translateRotate.scale+')'" class="room" id="room" v-for="(room, roomIndex) in rooms">
        <Sequence
            v-for="(el, elIndex) in room.elements"
            :separation=60
            :axis=true
            :dir=false
            :x=el.base[0]
            :y=el.base[1]
            :rotation=el.base[2]??0
            :num=1
            :atRepeats=-1
            :repeats=el.repeats
            :relevant-room-element=[roomIndex,elIndex]
        />
      </g>
    </svg>
    <div class="grid m-2">

      <div class="col-6">
        Scale {{ translateRotate.scale }}
        <Slider v-model="translateRotate.scale"
                :min=0.01
                :max=2
                :step=0.01
        />
      </div>
      <div class="col-6">
      </div>
      <div class="col-6">
        Area Width {{ settings.areaWidth }}dm
        <Slider v-model="settings.areaWidth"/>
      </div>
      <div class="col-6">
        Area Height {{ settings.areaHeight }}dm
        <Slider v-model="settings.areaHeight"/>
      </div>
      <div class="col-6">
        Area offset-x compared to center {{ settings.areaOffsetX }}%
        <Slider v-model="settings.areaOffsetX" :min=0 :max=100
        />
      </div>
      <div class="col-6">
        Area offset-y compared to center {{ settings.areaOffsetY }}%
        <Slider v-model="settings.areaOffsetY" :min=0 :max=100
        />
      </div>

      <div class="col-6">
        Area paddingX {{ settings.areaPaddingX }}dm
        <Slider v-model="settings.areaPaddingX"/>
      </div>
      <div class="col-6">
        Area paddingY {{ settings.areaPaddingY }}dm
        <Slider v-model="settings.areaPaddingY"/>
      </div>

      <div class="col-6">
        table offset x {{ settings.tableOffsetX }}dm
        <Slider v-model="settings.tableOffsetX" :min="-settings.areaPaddingX" :max="settings.areaPaddingX"/>
      </div>
      <div class="col-6">
        table offset y {{ settings.tableOffsetY }}dm
        <Slider v-model="settings.tableOffsetY" :min="-settings.areaPaddingY" :max="settings.areaPaddingY"/>
      </div>

      <div class="col-6">
        Separation between the seats
        <Slider v-model="settings.seatSep"/>
      </div>
      <div class="col-6">

      </div>
      <div class="col-6">
        Seat Distance to the table
        <Slider v-model="settings.seatDist" :min=0 :max="settings.areaHeight - settings.seatHeight"/>
      </div>

      <div class="col-6">
        Num Seats
        <Slider v-model="settings.seatNum"/>
      </div>
      <div class="col-6">
        Seat height
        <Slider v-model="settings.seatHeight"/>
      </div>
      <div class="col-6">
        Seat padding
        <Slider v-model="settings.seatPadding"/>
      </div>
    </div>
  </div>
  <div class="col-3" v-if="selectedElement">
    <Accordion multiple :activeIndex='[0]'>
      <AccordionTab v-for="(repeat, k) in selectedElement.repeats">
        <template #header>
            <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">{{ repeat.type }}</span>
                <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded
                        aria-label="confirm deletion" v-if="confirmdelete[k]" v-on:click.stop="deleteRepeats(k)"/>
                <Button class="ml-auto" size="small" icon="pi pi-times" severity="secondary" rounded aria-label="Delete"
                        outlined v-else v-on:click.stop="toggleDelete(k)"/>
            </span>
        </template>
        <div class="flex flex-column">
          <div class="flex flex-row">
            <div class="flex flex-shrink-0 flex-column m-2">
              <label :for="'type'+k">Shape type</label>
              <SelectButton v-model="repeat.type" :options="sequenceTypes" :inputId="'type'+k"/>
            </div>

            <div class="flex-grow-1 flex-shrink-1 flex flex-column m-2 w-full" v-if="repeat.type === 'Circle'">
              <label :for="'radius'+k">Radius</label>
              <InputNumber :id="'radius'+k" v-model="repeat.radius" showButtons/>
            </div>
          </div>

          <div class="flex-row field-checkbox m-2">
            <Checkbox v-model="repeat.axis" :inputId="k+'axis'" :binary="true"/>
            <label v-if="repeat.type === 'Line'" :for="k+'axis'"> horizontal</label>
            <label v-else :for="k+'axis'"> clockwise </label>
          </div>

          <div class="flex-row field-checkbox m-2">
            <Checkbox v-model="repeat.dir" :inputId="k+'dir'" :binary="true"/>
            <label v-if="repeat.type === 'Line'" :for="k+'dir'" class="ml-2"> negative</label>
            <label v-else :for="k+'dir'" class="ml-2"> backs facing</label>
          </div>

          <div class="flex-row m-2">
            <label :for="'repeat'+k">Repeat for #time</label>
            <InputNumber
                class="p-inputgroup"
                id="'repeat'+k"
                show-buttons
                v-model="repeat.num"/>
          </div>

          <div class="flex-row m-2">
            <label :for="'separation'+k">Separation</label>

            <div class="p-inputgroup">
              <span v-if="repeat.type === 'Circle'" class="p-inputgroup-addon">
                  <Checkbox v-model="repeat.equivalentSpaced" :binary="true"/>
              </span>
              <InputNumber
                  :inputId="'separation'+k"
                  :disabled="repeat.equivalentSpaced && repeat.type==='Circle'"
                  :modelValue="repeat.equivalentSpaced && repeat.type==='Circle' ? 360/Math.max(1, repeat.num) : repeat.separation"
                  @update:model-value="newValue => repeat.separation = newValue"
                  showButtons
                  :suffix="repeat.type === 'Circle' ? '°':' cm'"/>
            </div>
          </div>


          <div class="m-2">
            <label :for="'extraRotation'+k">Extra rotation</label>
            <InputNumber class="p-inputgroup" :id="'extraRotation'+k" v-model="repeat.extra" suffix="°" :min=-360
                         :max=360 show-buttons/>
          </div>
          <!--        </div>-->
        </div>
      </AccordionTab>
    </Accordion>
  </div>
</template>

<style scoped>
svg rect {
  color: lightgray;
}

.dragable {
  color: black;
}
</style>

<script setup>

import {teamareaStore} from "@/stores/teamarea";
import {roomTranslatorStore} from "@/stores/roomTranslator";
import Sequence from "@/components/Layout/Sequence.vue";
import {computed, onMounted, onUnmounted, reactive, ref} from "vue";
import {useKeyModifier, useScroll} from '@vueuse/core';

const settings = teamareaStore()
const translateRotate = roomTranslatorStore()

function scroll(state) {
  console.log(state) // {x, y, isScrolling, arrivedState, directions}
}

const map = ref(null)
useScroll(map, {onScroll: (e) => console.log(e)});

function toggleDelete(k) {
  console.log("delete", k)
  confirmdelete.value[k] = true
  window.setTimeout(() => delete confirmdelete.value[k], 1000)
}

function deleteRepeats(k) {
  rooms[0].elements[0].repeats.splice(k, 1);
  delete confirmdelete.value[k]
}

const sequenceTypes = ['Line', 'Circle'];

let data = [
  {
    name: "Room",
    type: "Rect",
    coords: [[0, 0], [10, 0], [10, 10], [0, 10]],
    elements: [
      {
        base: [100, 100],
        repeats: [
          {
            type: "Line",
            num: 3,
            axis: true,
            equivalentSpaced: true,
            dir: false,
          },
          {
            type: "Circle",
            num: 2,
            axis: false,
            separation: 50,
            dir: false,
          },
          // {
          //   type: "Line",
          //   axis: false,
          //   dir: true,
          //   separation: 50,
          // },
        ]
      }
    ]
  }
];

// Ensure all repeats have correct props
const fallback = {
  'type': 'Line',
  'extra': 0,
  'num': 1,
  'axis': false,
  'dir': false,
  'radius': 100,
  'separation': 100,
  'equivalentSpaced': true,
}

for (let ri = 0; ri < data.length; ri++) {
  for (let li = 0; li < data[ri].elements.length; li++) {
    data[ri].elements[li].repeats = data[ri].elements[li].repeats.map(e => {
      return {...fallback, ...e}
    })
  }
}

const rooms = reactive(data)
const selectedElement = ref(rooms[0].elements[0]);
const confirmdelete = ref([]);

const shift = useKeyModifier('Shift')
const control = useKeyModifier('Control')

const clamping = computed(() =>  (shift.value ? 5 : 0) * (control.value ? 3 : 1))

function maybeTranslateRotate(e) {
  if (translateRotate.translatingRoom) {

    const room = translateRotate.translatingRoom[0]
    const el = translateRotate.translatingRoom[1]

    // coord is in svg space
    const coord = [
      (e.clientX + translateRotate.offset[0])/translateRotate.scale,
      (e.clientY + translateRotate.offset[1])/translateRotate.scale,
    ];

    if (clamping.value>1) {
      const cv = clamping.value
      coord[0] = Math.round(coord[0]/cv) * cv
      coord[1] = Math.round(coord[1]/cv) * cv
    }

    rooms[room].elements[el].base[0] = coord[0]
    rooms[room].elements[el].base[1] = coord[1]
  } else if (translateRotate.rotatingRoom) {
    const room = translateRotate.rotatingRoom[0]
    const el = translateRotate.rotatingRoom[1]

    // coord is in svg space
    const coord = [
      e.clientX + translateRotate.offset[0],
      e.clientY + translateRotate.offset[1],
    ];

    const xDiff = coord[0] - rooms[room].elements[el].base[0]
    const yDiff = coord[1] - rooms[room].elements[el].base[1]

    const upVect = Math.atan2(-1, 0)
    let angle = (Math.atan2(yDiff, xDiff) - upVect)* 180 / Math.PI

    if (clamping.value>1) {
      angle = Math.round(angle/clamping.value)*clamping.value
    }

    rooms[room].elements[el].base[2] = angle
  }
}

function resetTranslateRotate() {
  translateRotate.translatingRoom = null;
  translateRotate.rotatingRoom = null;
  translateRotate.base = null;
  translateRotate.offset = null;
}

</script>
