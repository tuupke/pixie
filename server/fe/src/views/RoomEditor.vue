<template>
  <SvgEditor v-if="selectedRoom" ref="editor">
    <template #extra-buttons>
      <InputGroup class="ml-2">
        <InputGroupAddon size="small">Name</InputGroupAddon>
        <InputText v-model="selectedRoom!.name"/>
      </InputGroup>
      <InputGroup class="ml-2">
        <InputGroupAddon size="small">Padding</InputGroupAddon>
        <InputNumber
            v-model="outlinePadding"
            :min="0"
            placeholder="Padding"/>
        <ConfirmButton
            icon="pi pi-expand"
            label="outline"
            severity="warn"
            @confirmed="overrideOutline"
        />
      </InputGroup>
      <router-link to="/settings/map">
        <Button
            class="ml-2 mr-3"
            icon="pi pi-backward"
            label="back"
            severity="warning"
            outlined
        />
      </router-link>
      {{ selectedRoom.elements.length }}<span class="pi pi-bullseye ml-2 mr-3"/>
      {{
        selectedRoom!.elements.reduce((carry: number, element: ElementInterface): number => {
          return carry + element.repeats.reduce((p, r) => {
            return p * r.num
          }, 1)
        }, 0)
      }}<span class="pi pi-objects-column ml-2 mr-3"/>
    </template>
    <template #settings>
      <Card class="mt-2">
        <template #title>
          <span class="flex align-items-center gap-2 w-full">
              <span class="font-bold white-space-nowrap">ELEMENT</span>
              <Button
                  class="ml-auto mr-1"
                  size="small"
                  icon="pi pi-plus"
                  severity="success"
                  outlined
                  v-on:click.stop="addRepeats"/>
              <DeleteButton
                  v-if="highlightedKey"
                  class="mr-4"
                  :sequence-key="highlightedKey"
                  :upto="KeyCategory.Elements"
                  size="small"
                  icon="pi pi-trash"
                  severity="danger"
              />
          </span>
        </template>
        <template #content>
          <AssignmentOverview v-if="highlightedKey" :sequenceKey="highlightedKey"></AssignmentOverview>
        </template>
      </Card>
      <Accordion v-if="highlightedKey" class="mt-2">
        <AccordionPanel
            v-if="selectedElement!=null"
            v-for="(repeat, k) in selectedElement!.repeats"
            :value="k.toString()"
        >
          <AccordionHeader>
              <span class="flex align-items-center gap-2 w-full">
                  <span class="font-bold white-space-nowrap">R{{ k }}: {{ sequenceToString(repeat) }}</span>
                  <DeleteButton
                      class="ml-auto mr-2"
                      size="small"

                      icon="pi pi-trash"
                      severity="danger"
                      :sequenceKey="[KeyCategory.Placements, selectedRoomIndex, KeyCategory.Room, KeyCategory.Elements, selectedElementIndex!, KeyCategory.Repeats, k]"
                  />
              </span>
          </AccordionHeader>
          <AccordionContent>
            <div class="flex flex-column">
              <div class="flex flex-row">
                <div class="flex flex-shrink-0 flex-column m-2">
                  <label :for="'type'+k">Shape type</label>
                  <SelectButton :allowEmpty=false
                                v-model="repeat.type" :options="[SequenceType.Line, SequenceType.Circle]"
                                :inputId="'type'+k"/>
                </div>
              </div>

              <div class="flex-row m-2" v-if="repeat.type === SequenceType.Circle">
                <label :for="'repeat'+k">Radius</label>
                <InputNumber :id="'radius'+k" v-model="repeat.radius" showButtons class="p-inputgroup"/>
              </div>

              <div class="flex flex-shrink-0 flex-column m-2">
                <!--              <label v-if="repeat.type === SequenceType.Line" :for="k+'axis'">Orientation</label>-->
                <!--              <label v-else :for="k+'axis'"> clockwise </label>-->
                <label :for="k+'axis'">Axis</label>
                <SelectButton :allowEmpty=false
                              v-model="repeat.axis" :options="[SequenceAxis.Horizontal, SequenceAxis.Vertical]"
                              :option-label="axisName(repeat)"
                              :inputId="k+'axis'"/>
              </div>

              <div class="flex flex-shrink-0 flex-column m-2">
                <label v-if="repeat.type === SequenceType.Line" :for="k+'dir'" class="ml-2">
                  Direction
                </label>
                <label v-else :for="k+'dir'" class="ml-2"> backs facing</label>

                <SelectButton
                    :allowEmpty=false
                    v-model="repeat.dir"
                    :options="[SequenceDirection.Positive, SequenceDirection.Negative]"
                    :option-label="directionName(repeat)"
                    :inputId="k+'dir'"
                />
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
                <span v-if="repeat.type ===  SequenceType.Circle" class="p-inputgroup-addon">
                    <Checkbox v-model="repeat.equivalentSpaced" :binary="true"/>
                </span>
                  <InputNumber
                      :inputId="'separation'+k"
                      :disabled="repeat.equivalentSpaced && repeat.type=== SequenceType.Circle"
                      :modelValue="repeat.equivalentSpaced && repeat.type=== SequenceType.Circle ? 360/Math.max(1, repeat.num) : repeat.separation"
                      @update:model-value="(newValue: number) => repeat.separation = newValue"
                      showButtons
                      :suffix="repeat.type ===  SequenceType.Circle ? '°':' cm'"/>
                </div>
              </div>
              <div class="flex-row m-2">
                <label :for="'bitboxes'+k">Flip when repeats increases:</label>
                <div class="p-inputgroup">
                  <BitBoxes
                      :inputId="'bitboxes'+k"
                      :num="selectedElement!.repeats.length"
                      :options="selectedElement!.repeats.map(sequenceToString)"
                      :skip-bit="k"
                      v-model="repeat.snakeGroup"/>
                </div>
              </div>
            </div>
          </AccordionContent>
        </AccordionPanel>
      </Accordion>
    </template>
    <template #svg="{ rotate, moveStart }">
      <Room :translating=true
            ref="room"
            v-bind="selectedRoom"
            :sequenceKey="[KeyCategory.Placements, selectedRoomIndex]"
            @hoverElement="hoveringElement"
            @scrollElement="rotate"
            @dragStart="moveStart"
      />
    </template>
  </SvgEditor>
</template>

<style scoped>

</style>

<script setup lang="ts">

import SvgEditor from "./SvgEditor.vue";
import {
  ElementEvent,
  ElementInterface,
  Key,
  KeyCategory,
  QualifiedKey,
  Repeats,
  RoomInterface,
  RotateEvent,
  SequenceAxis,
  SequenceDirection,
  SequenceInterface,
  SequenceType
} from "../types.ts";
import BitBoxes from "../components/Layout/BitBoxes.vue";
import DeleteButton from "../components/Layout/DeleteButton.vue";
import AssignmentOverview from "../components/AssignmentOverview.vue";
import {computed, provide, ref} from "vue";
import {mapStore} from "../stores/map.ts";
import {onKeyStroke} from "@vueuse/core";
import Room from "../components/Layout/Room.vue";
import {teamareaStore} from "../stores/teamarea.ts";
import ConfirmButton from "../components/Layout/ConfirmButton.vue";

const map = mapStore()
const settings = teamareaStore()
const props = defineProps<{ selectedRoomIndex: number }>()
defineEmits<{
  moveStart: [ElementEvent]
  scrollElement: [RotateEvent]
}>()

const selectedRoom = computed<RoomInterface | undefined>(() => map.fromQualifiedKey([KeyCategory.Placements, props.selectedRoomIndex, KeyCategory.Room]))

// const selectedElementIndex = ref<number | null>(null);
const selectedElementIndex = computed<number | null>(() => {
  if (highlightedKey.value === null) {
    return null
  }

  const index = highlightedKey.value.indexOf(KeyCategory.Elements)
  if (index < 0 || highlightedKey.value?.length < index + 1) {
    return null;
  }

  const num = highlightedKey.value[index + 1]
  return Number.isFinite(num) ? num : null
});

const selectedElement = computed(() => {
  return selectedElementIndex.value === null || selectedRoom.value === null
      ? null
      : selectedRoom.value.elements[selectedElementIndex.value]
})

function hoveringElement(e: ElementEvent) {
  highlightedKey.value = e.key
}

const highlightedKey = ref<QualifiedKey | null>(null)
provide("highlightedKey", highlightedKey)

function resetHover() {
  highlightedKey.value = null
}

onKeyStroke('Escape', resetHover)
onKeyStroke('Tab', (e) => {
  if (nextElement(e.shiftKey ? -1 : 1)) {
    e.preventDefault()
  }
});

function nextElement(dir: number) {
  if (highlightedKey.value === null) {
    return false;
  }

  const k = [...highlightedKey.value];

  // TODO select next element
  let cat: Key = -1;
  let i: number;
  for (i = k.length - 1; i >= 0; i--) {
    if (typeof (k[i]) !== "number") {
      cat = k[i]
      break
    }
  }

  switch (cat) {
    case KeyCategory.Repeats:
      const elementKey = k.toSpliced(k.indexOf(KeyCategory.Elements) + 2)
      const element = map.fromQualifiedKey(elementKey);
      if (!element) {
        return false;
      }

      for (let atRepeat = 0; atRepeat < element.repeats.length; atRepeat++) {
        k[i + 1 + atRepeat] += dir

        // Stop when in bounds
        if (k[i + 1 + atRepeat] >= 0 && k[i + 1 + atRepeat] < element.repeats[atRepeat].num) {
          break
        }

        k[i + 1 + atRepeat] = dir > 0 ? 0 : element.repeats[atRepeat].num - 1;
      }

      highlightedKey.value = k;

      // TODO move to the next element

      return true;
  }
}


function addRepeats() {
  const element = map.fromQualifiedKeyUpTo(highlightedKey.value!, KeyCategory.Elements)
  element.repeats.push(new Repeats());
}


function sequenceToString(r: SequenceInterface): string {
  let direction: string;
  let separation: string;
  let radius: string = "";

  if (r.type === SequenceType.Circle) {
    radius = " ⌀" + (r.radius * 2) + settings.distanceUnit
    separation = r.separation.toString();
    if (r.equivalentSpaced) {
      separation = (360 / r.num).toString()
      if (r.axis === SequenceAxis.Horizontal) {
        direction = "⥁";
      } else {
        direction = "⥀";
      }
    } else {
      if (r.axis === SequenceAxis.Horizontal) {
        direction = "↻"
      } else {
        direction = "↺"
      }
    }
    separation += "°"
  } else {
    separation = "Δ" + r.separation + settings.distanceUnit;
    if (r.axis === SequenceAxis.Horizontal) {
      if (r.dir === SequenceDirection.Positive) {
        direction = "←"
      } else {
        // Negative
        direction = "→"
      }
    } else {
      // Vertical
      if (r.dir === SequenceDirection.Positive) {
        direction = "↓"
      } else {
        // Negative
        direction = "↑"
      }
    }
  }

  return `${r.type === SequenceType.Line ? '━' : "⭘"} ${separation} ${direction}${r.num}${radius}`
}


function axisName(s: SequenceInterface): (t: SequenceAxis) => string {
  return (t: SequenceAxis) => {
    switch (s.type) {
      case SequenceType.Line:
        switch (t) {
          case SequenceAxis.Horizontal:
            return "Horizontal"
          case SequenceAxis.Vertical:
            return "Vertical"
        }
        break
      case SequenceType.Circle:
        switch (t) {
          case SequenceAxis.Horizontal:
            return "Clockwise"
          case SequenceAxis.Vertical:
            return "Counterclockwise"
        }
        break
    }

    return 'unknown'
  }
}

function directionName(s: SequenceInterface): (t: SequenceDirection) => string {
  return (t: SequenceDirection) => {
    switch (s.type) {
      case SequenceType.Line:
        switch (s.axis) {
          case SequenceAxis.Horizontal:
            switch (t) {
              case SequenceDirection.Positive:
                return "Left"
              case SequenceDirection.Negative:
                return "Right"
            }
          case SequenceAxis.Vertical:
            switch (t) {
              case SequenceDirection.Positive:
                return "Behind"
              case SequenceDirection.Negative:
                return "In front"
            }
        }
      case SequenceType.Circle:
        switch (t) {
          case SequenceDirection.Positive:
            return "Backs"
          case SequenceDirection.Negative:
            return "Fronts"
        }
    }

    return 'unknown'
  }
}


const room = ref(null)
const editor = ref(null)
const outlinePadding = ref(0)

function overrideOutline() {
  if (room.value === null) {
    return
  }

  const snapping = editor.value.snapToGrid.value || true ? editor.value.translateClamping : 0

  selectedRoom.value.outline.splice(-selectedRoom.value.outline.length)
  selectedRoom.value.outline.push(...room.value.deriveOutline(outlinePadding.value, snapping))
}

</script>
