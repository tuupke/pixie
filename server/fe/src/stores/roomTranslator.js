import {defineStore} from 'pinia'

export const roomTranslatorStore = defineStore({
    id: 'roomTranslator',
    state: () => ({
        translatingRoom: null,
        rotatingRoom: null,
        offset: null,
    }),
})
