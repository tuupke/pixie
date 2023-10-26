<template>
<!--  style="display: flex;flex-direction: column;justify-content: space-between;"-->
  <div class="col-9 h-screen flex flex-column" style="justify-content: space-between;">
    <svg id="layoutsvg" width="100%" height="100%">
      <g class="room" id="room" v-for="room in rooms">
        <Sequence
            v-for="el in room.elements"
            :separation=60
            :axis="true"
            :dir="false"
            :x="el.base[0]"
            :y="el.base[1]"
            :rotation="el.base[2] ?? 0"
            :num=1
            :atRepeats=-1
            :repeats="el.repeats"/>
      </g>
    </svg>
    <div class="grid m-2">

      <div class="col-6">
        Area Width
        <Slider v-model="teamtableStore.areaWidth"/>
      </div>
      <div class="col-6">
        Area Height
        <Slider v-model="teamtableStore.areaHeight"/>
      </div>
      <div class="col-6">
        Seat Separation
        <Slider v-model="teamtableStore.seatSep"/>
      </div>
      <div class="col-6">
        Seat Distance
        <Slider v-model="teamtableStore.seatDist"/>
      </div>
      <div class="col-6">
        Margin-x
        <Slider v-model="teamtableStore.marginX"/>
      </div>
      <div class="col-6">
        Margin-y
        <Slider v-model="teamtableStore.marginY"/>
      </div>

      <div class="col-6">
        Num Seats
        <Slider v-model="teamtableStore.seatNum"/>
      </div>
      <div class="col-6">
        Seat height
        <Slider v-model="teamtableStore.seatHeight"/>
      </div>

      <div class="col-6">
        offsetX
        <Slider v-model="teamtableStore.offsetX" :min="-teamtableStore.marginX" :max="teamtableStore.marginX"/>
      </div>
      <div class="col-6">
        offsetY
        <Slider v-model="teamtableStore.offsetY" :min="-teamtableStore.marginY" :max="teamtableStore.marginY"/>
      </div>
    </div>
  </div>
  <div class="col-3" v-if="selectedElement">
    <Accordion multiple :activeIndex='[0]'  >
      <AccordionTab v-for="(repeat, k) in selectedElement.repeats">
        <template #header>
            <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">{{ repeat.type }}</span>
                <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded aria-label="confirm deletion" v-if="confirmdelete[k]"  v-on:click.stop="deleteRepeats(k)"/>
                <Button class="ml-auto" size="small" icon="pi pi-times" severity="secondary" rounded aria-label="Delete" outlined v-else v-on:click.stop="toggleDelete(k)"/>
            </span>
        </template>
        <div class="flex flex-column">
          <div class="flex flex-row">
            <div class="flex flex-shrink-0 flex-column m-2">
              <label :for="'type'+k">Shape type</label>
              <SelectButton v-model="repeat.type" :options="sequenceTypes" :inputId="'type'+k" />
            </div>

            <div class="flex-grow-1 flex-shrink-1 flex flex-column m-2 w-full" v-if="repeat.type === 'Circle'">
              <label :for="'radius'+k">Radius</label>
              <InputNumber :id="'radius'+k" v-model="repeat.radius" showButtons/>
            </div>
          </div>

          <div class="flex-row field-checkbox m-2">
            <Checkbox v-model="repeat.axis" :inputId="k+'axis'" :binary="true"/>
            <label v-if="repeat.type === 'Line'" :for="k+'axis'"> horizontal</label>
            <label v-else :for="k+'axis'" > clockwise </label>
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
                  <Checkbox v-model="repeat.equivalentSpaced" :binary="true" />
              </span>
              <InputNumber
                :inputId="'separation'+k"
                :disabled="repeat.equivalentSpaced && repeat.type==='Circle'"
                :modelValue="repeat.equivalentSpaced && repeat.type==='Circle' ? 360/Math.max(1, repeat.num) : repeat.separation"
                @update:model-value="newValue => repeat.separation = newValue"
                showButtons
                :suffix="repeat.type === 'Circle' ? '°':' cm'" />
            </div>
          </div>


          <div class="m-2">
            <label :for="'extraRotation'+k">Extra rotation</label>
            <InputNumber class="p-inputgroup" :id="'extraRotation'+k" v-model="repeat.extra" suffix="°" :min=-360 :max=360 show-buttons />
          </div>
<!--        </div>-->
        </div>
      </AccordionTab>
    </Accordion>
  </div>
</template>

<style scoped>
svg {border: 1px solid red;}
</style>

<script>
import {mapStores} from 'pinia'
import {teamtableStore} from "@/stores/teamtable";
import Sequence from "@/components/Layout/Sequence.vue";
import TeamTable from "@/components/Layout/TeamTable.vue";

export default {
  components: {TeamTable, Sequence},
  computed: {
    ...mapStores(teamtableStore),
  },
  methods: {
    toggleDelete(k) {
      this.confirmdelete[k] = true
      const that=this;
      window.setTimeout(function(){delete that.confirmdelete[k]}, 1000)
    },
    deleteRepeats(k) {
      this.rooms[0].elements[0].repeats.splice(k, 1);
      delete this.confirmdelete[k]
    }
  },
  data() {
    let data = {
      ShapeTypes: [
        {name: 'Line'},
        {name: 'Circle'},
      ],

      sequenceTypes: ['Line', 'Circle'],
      selectedElement: null,
      negative: false,
      isLine: true,
      inX: true,
      isEqui: true,
      confirmdelete: [],
      extra: 0,
      rooms: [
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
                {
                  type: "Line",
                  axis: false,
                  dir: true,
                  separation: 50,
                },
              ]
            }
          ]
        }
      ]
    }

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

    for (let ri = 0; ri < data.rooms.length; ri++) {
      for (let li = 0; li < data.rooms[ri].elements.length; li++) {
        data.rooms[ri].elements[li].repeats = data.rooms[ri].elements[li].repeats.map(e => {
          return {...fallback, ...e}
        })
      }
    }

    data.selectedElement = data.rooms[0].elements[0];

    // console.log(data);

    return data

  }
}

</script>
