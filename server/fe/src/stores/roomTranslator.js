import {defineStore} from 'pinia'

export const roomTranslatorStore = defineStore({
    id: 'roomTranslator',
    state: () => ({
        translatingRoom: null,
        rotatingRoom: null,
        offset: null,
        firstClick: null,
        scale: 1.0,
        translation: [0, 0]
    }),
})
