<template>
  <svg id="drawnMap" width='650' height='600' viewBox='0 0 650 600'>
    <desc>
      You may click and drag the ellipse.
      Make the ellipse rotate by grabbing the small circle and dragging.
    </desc>
    <defs>
    </defs>

  </svg>
</template>

<style scoped>



path,
line,
polygon{
  stroke: dodgerblue;
}

path{cursor:move;}

circle{
  fill: dodgerblue;
  cursor: alias;
}


</style>

<script setup>
const SVGLink_NS = 'http://www.w3.org/1999/xlink';
const SVG_NS = 'http://www.w3.org/2000/svg';
const svg = document.querySelector("#drawnMap");
const deg = 180 / Math.PI;
let rotating = false;
let dragging = false;
let impact = {
  x: 0,
  y: 0
};
let m = { //mouse
  x: 0,
  y: 0
};
let delta = {
  x: 0,
  y: 0
};
let ry = []; // elements array
let objectsRy = [];

let ellipsePath = {
  properties: {
    d: 'M-90,0 a90,50 0 1, 0 0,-1  z',
    fill: 'url(#trama)'
  },
  //parent: this.g,
  tagName: 'path',
  pos: {
    x: 370,
    y: 150
  }
}
objectsRy.push(ellipsePath);

function Element(o, index) {
  this.g = document.createElementNS(SVG_NS, 'g');
  this.g.setAttributeNS(null, 'id', index);
  svg.appendChild(this.g);

  o.parent = this.g;

  this.el = drawElement(o);
  this.a = 0;
  this.tagName = o.tagName;
  this.elRect = this.el.getBoundingClientRect();
  this.svgRect = svg.getBoundingClientRect();
  this.Left = this.elRect.left - this.svgRect.left;
  this.Right = this.elRect.right - this.svgRect.left;
  this.Top = this.elRect.top - this.svgRect.top;
  this.Bottom = this.elRect.bottom - this.svgRect.top;

  this.LT = {
    x: this.Left,
    y: this.Top
  };
  this.RT = {
    x: this.Right,
    y: this.Top
  };
  this.LB = {
    x: this.Left,
    y: this.Bottom
  };
  this.RB = {
    x: this.Right,
    y: this.Bottom
  };
  this.c = {
    x: 0, //(this.elRect.width / 2) + this.Left,
    y: 0 //(this.elRect.height / 2) + this.Top
  };
  this.o = {
    x: o.pos.x,
    y: o.pos.y
  };

  this.A = Math.atan2(this.elRect.height / 2, this.elRect.width / 2);
  console.log(this.A );
  this.pointsValue = function() { // points for the box
    return (this.Left + "," + this.Top + " " + this.Right + "," + this.Top + " " + this.Right + "," + this.Bottom + " " + this.Left + "," + this.Bottom + " " + this.Left + "," + this.Top);
  }

  let box = {
    properties: {
      points: this.pointsValue(),
      fill: 'none',
      stroke: 'dodgerblue',
      'stroke-dasharray': '5,5'
    },
    parent: this.g,
    tagName: 'polyline'
  }
  this.box = drawElement(box);

  let leftTop = {
    properties: {
      cx: this.LT.x,
      cy: this.LT.y,
      r: 6,
      fill: "blue"
    },
    parent: this.g,
    tagName: 'circle'
  }

  this.lt = drawElement(leftTop);

  this.update = function() {
    let transf = 'translate(' + this.o.x + ', ' + this.o.y + ')' + ' rotate(' + (this.a * deg) + ')';
    this.el.setAttributeNS(null, 'transform', transf);
    this.box.setAttributeNS(null, 'transform', transf);
    this.lt.setAttributeNS(null, 'transform', transf);
  }

}

for (let i = 0; i < objectsRy.length; i++) {
  let el = new Element(objectsRy[i], i + 1);
  el.update();
  ry.push(el)
}

// EVENTS

svg.addEventListener("mousedown", function(evt) {

  let index = parseInt(evt.target.parentElement.id) - 1;
  if (evt.target.tagName == ry[index].tagName) {
    dragging = index + 1;
    impact = oMousePos(svg, evt);
    delta.x = ry[index].o.x - impact.x;
    delta.y = ry[index].o.y - impact.y;
  }

  if (evt.target.tagName == "circle") {
    rotating = parseInt(evt.target.parentElement.id);
  }

}, false);

svg.addEventListener("mouseup", function(evt) {
  rotating = false;
  dragging = false;
}, false);

svg.addEventListener("mouseleave", function(evt) {
  rotating = false;
  dragging = false;
}, false);

svg.addEventListener("mousemove", function(evt) {
  m = oMousePos(svg, evt);

  if (dragging) {
    let index = dragging - 1;
    ry[index].o.x = m.x + delta.x;
    ry[index].o.y = m.y + delta.y;
    ry[index].update();
  }

  if (rotating) {
    let index = rotating - 1;
    console.log(ry[index].A);
    ry[index].a = Math.atan2(ry[index].o.y - m.y, ry[index].o.x - m.x) - ry[index].A;
    ry[index].update();
  }
}, false);

// HELPERS

function oMousePos(svg, evt) {
  let ClientRect = svg.getBoundingClientRect();
  return { //objeto
    x: Math.round(evt.clientX - ClientRect.left),
    y: Math.round(evt.clientY - ClientRect.top)
  }
}

function drawElement(o) {
  /*
  let o = {
    properties : {
    x1:100, y1:220, x2:220, y2:70},
    parent:document.querySelector("svg"),
    tagName:'line'
  }
  */
  let el = document.createElementNS(SVG_NS, o.tagName);
  for (let name in o.properties) {
    if (o.properties.hasOwnProperty(name)) {
      el.setAttributeNS(null, name, o.properties[name]);
    }
  }
  o.parent.appendChild(el);
  return el;
}
</script>