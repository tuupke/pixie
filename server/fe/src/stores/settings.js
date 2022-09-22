import {defineStore} from 'pinia'
// Import axios to make HTTP requests
import axios from "axios"

export const settingsStore = defineStore("settings", {
    state: () => ({
        settings: [],
        problems: [],
        hosts: [],
    }),
    getters: {
        map(state) {
            return state.settings.filter(e => e.key === "map")[0];
        },
    },
    actions: {
        async fetchSettings() {
            try {
                const data = await axios.get('/api/setting')
                this.settings = data.data

                for (let i in this.settings) {
                    try {
                        this.settings[i].JSON = JSON.parse(this.settings[i].value)
                    } catch (e) {
                        // Do nothing
                    }
                }

                const problems = await axios.get('/api/problem')
                this.problems = problems.data

                const hosts = await axios.get('/api/host')
                this.hosts = hosts.data
            } catch (error) {
                alert(error)
                console.log(error)
            }
        },
        // async saveSetting(guid) {
        //     const toSave = this.settings.find(el => el.guid === guid)
        //     if (toSave) {
        //         axios.put("http://localhost:4000/api/setting/" + toSave.key, JSON.stringify(toSave))
        //     }
        // }
    },
})