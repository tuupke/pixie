<template>
  <Toolbar class="flex-initial flex m-2">
    <template #start>
      <Button label="New Room" icon="pi pi-plus" class="mr-2" severity="success" />
      <Button label="New Room element" icon="pi pi-plus" class="mr-2" severity="success" />
    </template>
  </Toolbar>
  <div class="settings-view-wrapper">
    <div class="col-3 left-nav">
      <Accordion v-if="translateRotate.selectedRoom" multiple :activeIndex='[0]'>
        <AccordionTab>
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">Settings</span>
                              <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded
                                      aria-label="confirm deletion"
                                      v-if="confirmdelete[Math.pow(2, translateRotate.selectedRoom[0]) + Math.pow(3, translateRotate.selectedRoom[1]) + Math.pow(5, k)]"
                                      v-on:click.stop="deleteRepeats(k)"/>
                <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded aria-label="Delete"
                        outlined v-else v-on:click.stop="toggleDelete(k)"/>
                <Button class="ml-auto" size="small" icon="pi pi-plus" severity="success" rounded outlined
                        aria-label="add repeats" @click="addRepeats()"/>
            </span>
          </template>
        </AccordionTab>
        <AccordionTab
            v-for="(repeat, k) in rooms[translateRotate.selectedRoom[0]].elements[translateRotate.selectedRoom[1]].repeats">
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">{{ repeat.type }}</span>
                <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded
                        aria-label="confirm deletion"
                        v-if="confirmdelete[Math.pow(2, translateRotate.selectedRoom[0]) + Math.pow(3, translateRotate.selectedRoom[1]) + Math.pow(5, k)]"
                        v-on:click.stop="deleteRepeats(k)"/>
                <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded aria-label="Delete"
                        outlined v-else v-on:click.stop="toggleDelete(k)"/>
            </span>
          </template>
          <div class="flex flex-column">
            <div class="flex flex-row">
              <div class="flex flex-shrink-0 flex-column m-2">
                <label :for="'type'+k">Shape type</label>
                <SelectButton v-model="repeat.type" :options="sequenceTypes" :inputId="'type'+k"/>
              </div>

<!--              <div class="flex-grow-1 flex-shrink-1 flex flex-column m-2" v-if="repeat.type === 'Circle'">-->
<!--                <label :for="'radius'+k">Radius</label>-->
<!--                <InputNumber :id="'radius'+k" v-model="repeat.radius" showButtons/>-->
<!--              </div>-->
            </div>


            <div class="flex-row m-2" v-if="repeat.type === 'Circle'">
              <label :for="'repeat'+k">Radius</label>
              <InputNumber :id="'radius'+k" v-model="repeat.radius" showButtons class="p-inputgroup" />
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
              <label :for="'repeat'+k">Repeat for #times</label>
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
          </div>
        </AccordionTab>
      </Accordion>
      <div v-else>
        Select an element first
      </div>
    </div>
    <div class="col-9">
      <svg id="layoutsvg" @wheel="scroll" width="100%" height="100%" @mousemove="maybeTranslateRotate" @mouseup="resetTranslateRotate">
        <g :transform="'scale('+translateRotate.scale+')'" class="room" id="room" v-for="(room, roomIndex) in rooms">
          <Sequence
              v-for="(el, elIndex) in room.elements"
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
    </div>
  </div>
</template>

<style scoped>
svg rect {
  color: lightgray;
}

.dragable {
  color: black;
}
.settings-view-wrapper {
    display: flex;
    flex: 1;
}
.left-nav {
    max-height: calc(100vh - 9rem);
    overflow-y: auto;
}
</style>

<script setup>

import {teamareaStore} from "@/stores/teamarea";
import {roomTranslatorStore} from "@/stores/roomTranslator";
import Sequence from "@/components/Layout/Sequence.vue";
import {computed, reactive, ref} from "vue";
import {useKeyModifier, useScroll} from '@vueuse/core';

const settings = teamareaStore()
const translateRotate = roomTranslatorStore()

function scroll(state) {
    if ((translateRotate.scale >= 1 && state.deltaY > 0) || (translateRotate.scale <= 0.3 && state.deltaY < 0)) {
        return;
    }
    translateRotate.scale += (state.deltaY / 10000);
}

const map = ref(null)
useScroll(map, {onScroll: (e) => console.log(e)});

function toggleDelete(k) {
  const key = Math.pow(2, translateRotate.selectedRoom[0]) + Math.pow(3, translateRotate.selectedRoom[1]) + Math.pow(5, k)
  confirmdelete.value[key] = true
  window.setTimeout(() => delete confirmdelete.value[key], 1000)
}

function addRepeats() {
  rooms[translateRotate.selectedRoom[0]].elements[translateRotate.selectedRoom[1]].repeats.push({...fallback});

}

function deleteRepeats(k) {
  const key = Math.pow(2, translateRotate.selectedRoom[0]) +
      Math.pow(3, translateRotate.selectedRoom[1]) +
      Math.pow(5, k)
  if (!confirmdelete.value[key]) {
    return
  }

  delete confirmdelete.value[k]
  rooms[translateRotate.selectedRoom[0]].elements[translateRotate.selectedRoom[1]].repeats.splice(k, 1);
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
            separation: 1000
          },
          {
            type: "Circle",
            num: 2,
            axis: false,
            radius: 500,
            dir: false,
          },
        ]
      },
      {
        base: [100, 400],
        repeats: []
      }
    ]
  }
];

// Ensure all repeats have correct props
const fallback = {
  'type': 'Line',
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
const confirmdelete = ref([]);

const shift = useKeyModifier('Shift')
const control = useKeyModifier('Control')

const clamping = computed(() => (shift.value ? 5 : 0) * (control.value ? 3 : 1))

function maybeTranslateRotate(e) {
  if (translateRotate.translatingRoom) {
    const room = translateRotate.translatingRoom[0]
    const el = translateRotate.translatingRoom[1]

    // coord is in svg space
    const coord = [
      (e.clientX + translateRotate.offset[0]) / translateRotate.scale,
      (e.clientY + translateRotate.offset[1]) / translateRotate.scale,
    ];

    if (clamping.value > 1) {
      const cv = clamping.value
      coord[0] = Math.round(coord[0] / cv) * cv
      coord[1] = Math.round(coord[1] / cv) * cv
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
    let angle = (Math.atan2(yDiff, xDiff) - upVect) * 180 / Math.PI

    if (clamping.value > 1) {
      angle = Math.round(angle / clamping.value) * clamping.value
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
