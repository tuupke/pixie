import { defineStore } from 'pinia'
import axios from "axios"

export const teamsStore = defineStore("teams",{
    state: () => ({
        teams: [],
    }),
    getters: {
        hasLocations() {
            return this.teams.reduce((p, e) => p + e.location.x + e.location.y + e.location.rotation, 0) > 0;
        }
    },
    actions: {
        async fetchTeams() {
            try {
                const data = await axios.get('http://localhost:4000/api/external_data')
                this.teams = data.data
            }
            catch (error) {
                alert(error)
                console.log(error)
            }
        },
        updateTeams() {
            axios.get('http://localhost:4000/api/djLoad').then(this.fetchTeams)
        }
    },
})