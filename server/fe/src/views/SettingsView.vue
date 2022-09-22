<template>
  <div class="col-12">
  <svg id="layoutsvg">
    <TeamTable
        :seat-sep="seatSep" :seat-dist="seatDist" :seat-num="numSeats"
        :area-width="areaWidth"
        :area-height="areaHeight"
        :margin="margin"
        :x="30" :y="30">13</TeamTable>

    <TeamTable
        :seat-sep="seatSep" :seat-dist="seatDist" :seat-num="numSeats"
        :area-width="areaWidth"
        :area-height="areaHeight"
        :margin="margin"
        :x="150" :y="30">14</TeamTable>

    <TeamTable
        :seat-sep="seatSep" :seat-dist="seatDist" :seat-num="numSeats"
        :area-width="areaWidth"
        :area-height="areaHeight"
        :margin="margin"
        :x="300" :y="30">15</TeamTable>
  </svg>
  </div>

  <div  class="col-6">
    Area Width
    <Slider v-model="areaWidth" />
  </div>
  <div  class="col-6">
    Area Height
    <Slider v-model="areaHeight" />
  </div>
  <div  class="col-6">
    Seat Separation
    <Slider v-model="seatSep" />
  </div>
  <div  class="col-6">
    Margin
    <Slider v-model="margin" />
  </div>

  <div  class="col-6">
    Num Seats
    <Slider v-model="numSeats" />
  </div>

</template>

<script>

export default {
  data() {
    return {
      width: 40,
      height: 15,
      scale: 20,
      seatWidth: 10,
      seatHeight: 3,
      seatSep: 3,
      seatDist: 3,
      margin: 5,

      areaWidth: 50,
      areaHeight: 35,
      numSeats: 3,
    }
  },
}

// area
// path
// row
// table
//


</script>
