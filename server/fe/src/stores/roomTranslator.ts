import {defineStore} from 'pinia'

export const roomTranslatorStore = defineStore('roomTranslator', {
    state: () => {
        return {
            translatingRoom: null as string | null,
            rotatingRoom: null as string[] | null,
            firstClick: null as number[] | null,
            scale: 0.5 as number,
            translation: [0, 0] as number[],
            selectedRoom: null as number[] | null,
        }
    },
})
