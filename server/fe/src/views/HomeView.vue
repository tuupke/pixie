<style>
.teamtable rect {
  stroke: #444;
  stroke-width: 1;
  stroke-dasharray: none;
  stroke-linecap: butt;
  stroke-dashoffset: 0;
  stroke-linejoin: miter;
  stroke-miterlimit: 4;
  fill-rule: nonzero;
  opacity: 1;

  fill: #fff;
}

.aisle {
  stroke: #aaa;
  stroke-width: 1;
}

.selectedteam rect {
  stroke-width: 2;
}

.teamtable.found rect {
  fill: orange;
}

.teamtable.found.exists  rect {
  fill: lightgreen;
}

.teamtable.noteam rect {
  opacity: 0.3;
}

.teamtable:hover rect {
  stroke-width: 3;
}

.printer, .problem {
  fill: #aaa;
}

</style>

<template>
  <Toolbar class="col-12">
    <template #start>
      <FileUpload class="p-button-danger mr-2" mode='basic' name="files" url="http://localhost:4000/api/tim-json"
                  :auto="true" :multiple="true"
                  chooseLabel="Tim JSON upload" :fileLimit="4" accept="application/json" @upload="reloadSettings"/>

      <SplitButton label="Reload settings" icon="pi pi-check" :model="reloadItems"
                   class="p-button-warning"></SplitButton>

      <span v-if="settingsStore.map">
      &nbsp;Scanning mode&nbsp;
      <Dropdown v-model="scanningMode" :options="scanningModeTypes" placeholder="Select a scanningmode"/>
        </span>
    </template>
  </Toolbar>

  <Card v-if="teamsStore.hasLocations" class="col-9">
    <template #title>
      Contest layout and registration
    </template>
    <template #content>
      <svg :viewBox="viewBox()" id="layoutsvg" @click="clearTeam(null, true)">
        <g id="sceneroot">
          <g v-for="team in teamsStore.teams"
             :transform="matrixCalc(team.location.x,team.location.y, team.location.rotation)"
             @click.stop="setupTeam" v-on:mouseover.stop="showTeam" v-on:mouseout.stop="clearTeam"
             :id="'team_' + team.guid"
             :class="'teamtable ' + (team.host_id ? 'found ' : ' ')  + (settingsStore.hosts.find(e => e.guid === team.host_id) ? 'exists' : '') + (team.team ? '' : 'noteam')">
            <g v-if="team.location.x + team.location.y + team.location.rotation > 0">
              <rect x="-20" y="-10" width="40" height="15"/>
              <text font-family="sans-serif" font-size="15" font-style="normal" x=0 y=3
                    text-anchor="middle"
                    font-weight="normal">
                {{ team.team_id }}
              </text>
              <rect x="-17" y="7" rx="0" ry="0" width="10" height="3"/>
              <rect x="-5" y="7" rx="0" ry="0" width="10" height="3"/>
              <rect x="07" y="7" rx="0" ry="0" width="10" height="3"/>
            </g>
          </g>
        </g>
        <g v-if="settingsStore.map !== undefined && settingsStore.map.value !== undefined">
          <line v-for="aisle in settingsStore.map.value.aisles" class="aisle"
                :x1="scale*aisle.x1"
                :y1="scale*aisle.y1"
                :x2="scale*aisle.x2"
                :y2="scale*aisle.y2"
          ></line>
          <rect v-if="settingsStore.map.value.printer" class="printer" :width="scale" :height="scale"
                :x="scale*(settingsStore.map.value.printer.x)-scale/2"
                :y="scale*(settingsStore.map.value.printer.y)-scale/2"/>

          <g v-for="problem in settingsStore.problems"
             :transform="'translate(' + (scale*problem.location.x) + ',' + (scale*problem.location.y) +')'">
            <circle :r="scale/2" :style="(problem.rgb) ? 'fill: ' + problem.rgb: ''" class="problem"/>
            <text font-family="sans-serif" font-size="15" font-style="normal" x=0 y=5 text-anchor="middle"
                  font-weight="normal">
              {{ problem.id }}
            </text>
          </g>
        </g>
      </svg>

      <img hidden src="/public/image.png" id="qrimage"/>
    </template>
  </Card>
  <Card v-else class="col-9">
    <template #title>Teamlist</template>
    <template #content>
      <DataTable :scrollable="true" scrollHeight="400px" :value="teamsWithTeam"
                 selectionMode="single"
                 v-on:row-click="() => this.selected = true"
                 v-model:selection="modifyingteam"
                 :rowClass="() => 'teamtable'">
        <Column :sortable="true" field="guid" header="guid"></Column>
        <Column :sortable="true" field="team" header="name"></Column>
        <Column :sortable="true" field="team_id" header="team-id"></Column>
        <Column :sortable="true" field="host_id" header="host-id"></Column>
      </DataTable>
    </template>
  </Card>

  <Card class="col-3" v-if="modifyingteam.guid !== undefined">
    <template #title>
      {{ modifyingteam.team ?? 'No teamname yet' }}
    </template>
    <template #subtitle>
      {{ modifyingteam.guid }}
    </template>
    <template #content>
      <table :set="fh = settingsStore.hosts.find(e => e.guid === modifyingteam.host_id)">
        <tr>
          <td>User:</td>
          <td> {{ modifyingteam.username }} ({{ modifyingteam.id }})</td>
        </tr>
        <tr>
          <td>Team:</td>
          <td> {{ modifyingteam.team }} ({{ modifyingteam.team_id }})</td>
        </tr>
        <tr>
          <td>Location:</td>
          <td> x {{ modifyingteam.location.x }}, y {{ modifyingteam.location.y }}, rot
            {{ modifyingteam.location.rotation }}°
          </td>
        </tr>

        <tr>
          <td>Host:</td>
          <td><Team v-model="modifyingteam.host_id" :options="settingsStore.hosts" :filterFields="['hostname','guid']"
                    @onchange="() => {this.modifyingteam.host_id = null; setHost()}">
            <template v-slot:option="slotProps">{{ slotProps.hostname }} {{ slotProps.guid }}</template>
            <template v-slot:value="slotProps">{{ slotProps.hostname }}</template>
          </Team></td>
        </tr>

        <tr>
          <td>Host guid:</td>
          <td v-if="modifyingteam.host_id">{{ fh.guid }}</td>
          <td v-else>--</td>
        </tr>
        <tr>
        <td>Host name:</td>
        <td v-if="modifyingteam.host_id">{{ fh.hostname }}</td>
        <td v-else>--</td>
      </tr>
        <tr>
          <td>Primary ip:</td>
          <td v-if="modifyingteam.host_id">{{ fh.primary_ip }}</td>
          <td v-else>--</td>
        </tr>
        <tr>
          <td>Primary MAC:</td>
          <td v-if="modifyingteam.host_id">{{ fh.primary_mac }}</td>
          <td v-else>--</td>
        </tr>
      </table>

      <div v-if="selected">
        <Button @click="callScanner" class="p-button-raised p-button-rounded p-button-danger">Scan</Button>
        <Button v-if='scanningMode != "single"' class="p-button-raised p-button-rounded p-button-warning"
                @click="nextTeam">Next
        </Button>
        <Button v-if="modifyingteam.host_id" @click="() => {modifyingteam.host_id = null; setHost();}">Clear Host</Button>
        <video id="qr-scanner" style="width: 100%;"></video>
      </div>
    </template>
  </Card>
  <Card class="col-3" v-else>
    <template #title>
      {{ 'No team selected' }}
    </template>
    <template #content>
      Hover over, or click on, a team to show information in this panel.
    </template>
  </Card>
</template>

<script>

import axios from "axios";
import {mapStores} from 'pinia'

// declare store variable
import {settingsStore} from "../stores/settings";
import {teamsStore} from "../stores/teams";
import QrScanner from "qr-scanner";

import * as flatbuffers from 'flatbuffers';
import {Ping} from "../../packets"

export default {
  data() {
    return {
      'modifyingteam': {},
      'selected': false,
      'scale': 20,
      'scanningMode': 'single',
      'scanningModeTypes': ['single', 'all', 'unlinked'],
      'qrScanner': null,
      'reloadItems': [{
        label: 'Pixie',
        command: () => {
          this.reloadSettings();
          this.$toast.add({severity: 'success', summary: 'Updated', detail: 'Pixie data updated', life: 3000});
        }
      }, {
        label: 'DOMjudge',
        command: () => {
          axios.get("http://localhost:4000/api/djTeam")
          axios.get("http://localhost:4000/api/djProblem")
          this.reloadSettings()
          this.$toast.add({
            severity: 'success',
            summary: 'Updated',
            detail: 'DOMjudge & Pixie data updated',
            life: 3000
          });
        }
      }],
    }
  },
  computed: {
    // note we are not passing an array, just one store after the other
    // each store will be accessible as its id + 'Store'
    ...mapStores(settingsStore, teamsStore), host() {
      return undefined
    }, teamsWithTeam() {
      return this.teamsStore.teams.filter((e) => e.team)
    }
  },
  methods: {
    enable: function (a) {
      document.getElementById(a.target.id + '_save').disabled = false;
    },
    save: function ($event, a) {
      document.getElementById($event.target.id).disabled = true;
      this.settingsStore.saveSetting(a)
    },
    saveTim: function ($event) {
      document.getElementById("timJson_save").disabled = true;
      const val = document.getElementById("timJson").value;

      axios.post("http://localhost:4000/api/tim-json/", val).catch(e => window.alert(e)).finally(this.teamsStore.fetchTeams())
    },
    matrixCalc: function (x, y, rot) {
      return `translate(${this.scale * x}, ${this.scale * y}) rotate(${90 + rot})`
    },
    viewBox: function () {
      let maxX = 0, maxY = 0;
      for (let i in this.teamsStore.teams) {
        maxX = Math.max(maxX, this.teamsStore.teams[i].location.x)
        maxY = Math.max(maxY, this.teamsStore.teams[i].location.y)
      }

      if (this.settingsStore.problems) {
        for (let i in this.settingsStore.problems) {
          maxX = Math.max(maxX, this.settingsStore.problems[i].location.x)
          maxY = Math.max(maxY, this.settingsStore.problems[i].location.y)
        }
      }

      maxX += 3;
      maxY += 3;

      return -1.5 * this.scale + ' ' + -1.5 * this.scale + ' ' + this.scale * maxX + ' ' + this.scale * maxY;
    },
    showTeam: function ($event) {
      this.presentTeam($event.target, false)
    },
    clearTeam: function ($event, override) {
      if (this.selected && !override) {
        return
      }
      this.selected = false
      this.modifyingteam = {}

      if (this.qrScanner) {
        this.qrScanner.stop()
      }
      this.justReset()
    },
    justReset: function () {
      const reset = document.querySelectorAll(".selectedteam")
      for (let i in reset) {
        if (!reset[i].classList) {
          continue;
        }
        reset[i].classList.remove("selectedteam")
      }
    },
    setupTeam: function ($event) {
      this.presentTeam($event.target ?? $event.originalEvent.target, true)
    },
    callScanner: async function () {
      const {data} = await axios.get("http://localhost:5173/public/image.png", {
        responseType: "Uint8Array",
      })

      let that = this

      function handleResult(result) {
        try {
          that.qrScanner.stop()
        } catch (e) {
          // Do nothing
        }

        result = result.data || result

        // Convert the results to a Uint8Array. https://stackoverflow.com/a/21797381
        const binary_string = window.atob(result);
        const len = binary_string.length;
        let bytes = new Uint8Array(len);
        for (var i = 0; i < len; i++) {
          bytes[i] = binary_string.charCodeAt(i);
        }

        const decoded = new flatbuffers.ByteBuffer(bytes)

        let greeter = Ping.getRootAsPing(decoded)

        this.modifyingteam.host_id = greeter.identifier()
        // console.log(greeter.identifier() + ' ' + greeter.hostname())
      }

      // QrScanner.scanImage(document.getElementById("qrimage")).then(handleResult);
      // return
      this.qrScanner = new QrScanner(document.getElementById("qr-scanner"),
          handleResult, {
            returnDetailedScanResult: true,
            maxScansPerSecond: 5,
            highlightScanRegion: true
          },
      );

      this.qrScanner.start();

    },
    setHost: function () {
      const that = this;
      window.setTimeout(function() {
        axios.patch(`http://localhost:4000/api/external_data/${that.modifyingteam.guid}`, {
          "host_id": that.modifyingteam.host_id
        }).then(this.reloadSettings)
      }, 300)

    },
    nextTeam: function () {
      if (!this.modifyingteam) {
        return
      }

      const numTeams = this.teamsStore.teams.filter((e) => e.location.x + e.location.y + e.location.rotation > 0).length

      let attempts = 0;
      while (attempts <= numTeams) {
        attempts++;

        const thisIndex = this.modifyingteam.team_id - 1
        const nextIndex = ((thisIndex + attempts) % numTeams) + 1
        this.justReset()
        const nextTeam = this.teamsStore.teams.find(t => t.team_id == nextIndex)
        // TODO handle existing host
        if (!nextTeam) {
          window.alert("Could not find next team! " + nextIndex)
          this.clearTeam(null, true)
          return
        }

        const teamTableId = "team_" + nextTeam.guid
        const nextTeamElement = document.getElementById(teamTableId)
        if (!nextTeamElement) {
          window.alert("Could not find next team table! " + teamTableId)
          this.clearTeam(null, true)
          return
        }

        this.presentTeam(nextTeamElement, true)
        return;
      }

      window.alert("Tried them all, no next is available.")
    },
    reloadSettings: function () {
      window.setTimeout(() => {
        this.settingsStore.fetchSettings();
        this.teamsStore.fetchTeams();
      }, 500)
    },
    presentTeam: function (target, override) {
      if (!override && this.selected) {
        return
      }

      this.clearTeam(null, override)

      let tries = 5;
      while (tries-- > 0 && target.parentNode) {
        if (target.classList.contains('teamtable')) {
          break
        }

        target = target.parentNode
      }

      if (override) {
        this.selected = true
      }

      if (!target.classList.contains('teamtable')) {
        window.alert("No teamtable found!")
        return
      }

      target.classList.add("selectedteam")

      const teamId = target.id.substring(5)
      const foundTeam = this.teamsStore.teams.find(el => el.guid === teamId)
      if (foundTeam) {
        this.modifyingteam = foundTeam
      } else {
        window.alert("Team " + teamId + " not found")
      }
    }
  },
  async mounted() {
    if (!this.settingsStore.settings.length) {
      await this.settingsStore.fetchSettings();
      await this.teamsStore.fetchTeams();
    }
  }
}

</script>