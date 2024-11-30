<template>
  <SvgEditor>
    <template #extra-buttons>
      <ConfirmButton label="Room" icon="pi pi-plus" class="mr-3" severity="success" @confirmed="map.placements.push({})"/>
      {{ map.placements.length }}<span class="pi pi-map ml-1 mr-3" />
      {{
        map.placements.reduce((totalTables, roomPlacement): number => totalTables + roomPlacement.room.elements.length ?? 0, 0)
      }}<span class="pi pi-bullseye ml-2 mr-3" />
      {{
        map.placements.reduce((totalTables, roomPlacement): number => {
          return totalTables + roomPlacement.room.elements.reduce((previous, element): number => {
            return previous + element.repeats.reduce((p, r) => {
              return p * r.num
            }, 1)
          }, 0)
        }, 0)
      }}<span class="pi pi-objects-column ml-2 mr-3" />
    </template>
    <template #settings>
      <Accordion class="mt-2" :active-index="0">
        <AccordionPanel v-for="(placement, index) in map.placements" :value="index.toString()">
          <AccordionHeader>
            <span class="flex align-items-center gap-2 w-full">
              <span class="flex align-items-center gap-2 w-full">
                      <InputText v-model="placement.room.name" size="small"/>
              </span>
            </span>

            <router-link :to="'/settings/map/'+index">
              <Button icon="pi pi-pencil" size="small" outlined class="mr-2 " severity="warn" />
            </router-link>
            <DeleteButton
                class="mr-2"
                type="button"
                size="small"
                icon="pi pi-trash"
                severity="danger"
                :sequenceKey="[KeyCategory.Placements, index]"
            />

<!--            <Button icon="pi pi-trash" size="small" outlined class="mr-2" severity="danger" />-->
          </AccordionHeader>
          <AccordionContent>
            <div>
              Number of placed elements: {{ placement.room.elements.length }}
            </div>
            <div>
              Number of tables/areas: {{
                placement.room.elements.reduce((previous, element): number => {
                  return previous + element.repeats.reduce((p, r) => {
                    return p * r.num
                  }, 1)
                }, 0)
              }}
            </div>
          </AccordionContent>
        </AccordionPanel>
      </Accordion>
    </template>
    <template #svg="{ rotate, moveStart }">
      <Drag
          v-for="(placement, i) in map.placements"
          :coord=placement.coord
          @dragStart="console.log('start')"
      >
        <Room v-bind="placement.room"
              :coord="placement.coord"

              @scrollElement="(e) => rotate({key: [KeyCategory.Placements, i], event: e.event})"
              :sequenceKey="[KeyCategory.Placements, i]"
        />
      </Drag>
    </template>
  </SvgEditor>
</template>

<style scoped>

</style>

<script setup lang="ts">

import SvgEditor from "./SvgEditor.vue";
import {KeyCategory} from "../types.ts";
import Drag from "./Drag.vue";
import Room from "../components/Layout/Room.vue";
import {mapStore} from "../stores/map.ts";
import DeleteButton from "../components/Layout/DeleteButton.vue";

const map = mapStore()

</script>
