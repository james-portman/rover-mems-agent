// maybe just keep drawing the lines and make the canvas scroll? possible to just slide everything left constantly?

// start lines at lowest value, not 0
// print min and max next to labels

var canvases = [];//, 'coolant', 'lambdaV'];
var lineColours = ['aqua', 'yellow', 'lime', 'white'];

var canvasLabelX;
var canvasLabelY;

window.addEventListener('resize', windowResized);

// wait for resize to finish
var resizeId;
function windowResized() {
    clearTimeout(resizeId);
    resizeId = setTimeout(redrawCanvases, 400);
}

function redrawCanvases(){
  for (i in canvases) {
    var canvas = document.getElementById('canvas_'+canvases[i]);
    canvas.style.width ='100%';
    canvas.style.height='100%';
    canvas.width  = canvas.offsetWidth;
    canvas.height = canvas.offsetHeight;

    // only needs doing once but fine:
    canvasLabelX = canvas.height/6/4;
    canvasLabelY = canvas.height/6;
    if (i < lineColours.length) {
      lineColour = lineColours[i];
    } else {
      lineColour = lineColours[0];
    }
    drawit(canvases[i], lineColour);
  }
  window.setTimeout(redrawCanvases, 1000);
}

function drawit(dataField, lineColour) {
  var canvas = document.getElementById('canvas_'+dataField);
  var ctx = canvas.getContext('2d');

  ctx.fillRect(0, 0, canvas.width, canvas.height);
  // draw guidelines grey
  var num_guidelines = 3;
  for (var i=0; i<=canvas.height; i+=canvas.height/(num_guidelines+1)) {
    ctx.beginPath();
    ctx.lineWidth = 0.75;
    ctx.strokeStyle = '#666666';
    ctx.moveTo(0, i);
    ctx.lineTo(canvas.width, i);
    ctx.stroke();
  }
  ctx.lineWidth = 1.5;

  // TODO: STOP DOING THIS ON THE CANVAS
  // ctx.beginPath();
  // ctx.font = '1em Arial';
  // ctx.fillStyle = '#eee';
  // ctx.textAlign = 'left';
  // ctx.strokeStyle = '#eee';
  // // default text gets put outside the canvas as its bottom left is 0, 0
  // // if canvas is 6em tall then one line is height/6
  // ctx.fillText(dataField, canvas.height/6/4, canvas.height/6);

  // put this back after text
  ctx.fillStyle = '#000';

  var num_points = 100;

  x_distance = canvas.width / num_points;

  max = max_values[dataField];
  min = min_values[dataField];

  last_x = 0;
  last_y = 0;
  last_y = canvas.height - last_y;

  ctx.beginPath();
  ctx.strokeStyle = lineColour;

  var thisData;
  if (outputHistory.length > 100) {
    thisData = outputHistory.slice(outputHistory.length-100);
  } else {
    thisData = outputHistory;
  }

  var lastValue = null;
  for (var i=0; i < thisData.length; i++) {
    point_x_coord = x_distance * i;

    datapoint = thisData[i][dataField];
    if (datapoint == undefined) {
      if (lastValue == null) {
        continue;
      } else {
        datapoint = lastValue;
      }
    } else {
      datapoint = datapoint["data"];
      lastValue = datapoint;
    }

    point_y_percentage = datapoint / max * 100;
    point_y_coord = canvas.height / 100 * point_y_percentage;
    // draws from top left so have to flip the Y
    point_y_coord = canvas.height - point_y_coord;

    if (last_x == 0) {
      last_y = point_y_coord;
    }
    ctx.moveTo(last_x, last_y);

    ctx.lineTo(point_x_coord, point_y_coord);
    ctx.stroke();

    last_x = point_x_coord;
    last_y = point_y_coord;
  }

  // draw a vertical marker at the end
  ctx.beginPath();
  ctx.strokeStyle = '#aaa';
  ctx.moveTo(last_x, 0);
  ctx.lineTo(last_x, canvas.height);
  ctx.stroke();

}

redrawCanvases();

function removeGraph(dataKey) {

  index = canvases.indexOf(dataKey);
  if (index > -1) {
    canvases.splice(index, 1);
  }

  var search = "graph_top_level_"+dataKey;
  var graphDiv = document.getElementById(search);
  graphDiv.parentNode.removeChild(graphDiv);
}


function addGraph(dataKey, label) {

  if (canvases.indexOf(dataKey) > -1) {
    return;
  }
  var canvasesDiv = document.getElementById("canvases");

  var newGraphDiv = document.createElement("div");
  newGraphDiv.id = "graph_top_level_"+dataKey;
  canvasesDiv.appendChild(newGraphDiv);

  var deleteButton = document.createElement("span");
  deleteButton.innerHTML = '<i class="far fa-trash-alt deleteGraphButton" onclick="removeGraph(\''+dataKey+'\')"></i>&nbsp;';
  newGraphDiv.appendChild(deleteButton);

  var labelDiv = document.createElement("span");
  labelDiv.innerText = label+": ";
  newGraphDiv.appendChild(labelDiv);

  var valueDiv = document.createElement("span");
  valueDiv.id = "canvas_"+dataKey+"_value";
  newGraphDiv.appendChild(valueDiv);

  var canvasWrapper = document.createElement("div");
  canvasWrapper.style.width = "100%";
  canvasWrapper.style.height = "6em";
  canvasWrapper.classList.add("canvasWrapper");
  newGraphDiv.appendChild(canvasWrapper);

  var newCanvas = document.createElement("canvas");
  newCanvas.id = "canvas_"+dataKey;

  canvasWrapper.appendChild(newCanvas);

  canvases.push(dataKey);
}
