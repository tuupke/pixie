import {defineStore} from 'pinia'

export const roomTranslatorStore = defineStore({
    id: 'roomTranslator',
    state: () => ({
        translatingRoom: null,
        rotatingRoom: null,
        offset: null,
        firstClick: null,
        scale: 0.5,
        translation: [0, 0],
        selectedRoom: null,
    }),
})
