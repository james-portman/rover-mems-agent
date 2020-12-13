var maxHistoryLength = 10000;

numCollectedDatapointsSpan = document.getElementById("numCollectedDatapoints");

allOutputKeys = {};

min_values = {
  'rpm': -100,
  'coolant': -30,
  'lambda': 0,
  'map': 0
};

max_values = {
  'rpm': 6500,
  'coolant': 120,
  'lambda': 1,
  'map': 110
};

function clearFaults() {
  var faultsDiv = document.getElementById("faults");
  faultsDiv.innerHTML = "";
}

function outputData(data) {

  outputHistory.push(data);

  if (outputHistory.length > maxHistoryLength) {
    outputHistory = outputHistory.slice(outputHistory.length-maxHistoryLength);
  }
  numCollectedDatapointsSpan.innerText = outputHistory.length; // output num of collected data points

  var outputsDiv = document.getElementById("outputs");
  var faultsDiv = document.getElementById("faults");

  // each thing to output
  Object.keys(data).sort().forEach(function(key) {

    if (allOutputKeys[key] == undefined) {
      allOutputKeys[key] = true;
    }

    var outputLabel;
    var outputValue;
    if (data[key]["data"] != undefined && data[key]["name"]) {
      outputLabel = data[key]["name"];
      outputValue = data[key]["data"];
    } else {
      outputLabel = key;
      outputValue = data[key];
    }

    if (min_values[key] === undefined) {
      min_values[key] = outputValue;
    } else if (outputValue < min_values[key]) {
      min_values[key] = outputValue;
    }
    if (max_values[key] === undefined) {
      max_values[key] = outputValue;
    } else if (outputValue > max_values[key]) {
      max_values[key] = outputValue;
    }

    var element = document.getElementById("output_"+key);
    var outputValueDiv;
    if (element == null) {
      var element = document.createElement('div');
      element.id = "output_"+key;
      outputsDiv.appendChild(element);

      var outputGraphDiv = document.createElement('span');
      // outputGraphDiv.id = "output_"+key+"_graph";
      outputGraphDiv.classList.add("outputGraph");
      outputGraphDiv.innerHTML = '<i title="Add graph" class="fas fa-chart-line addChartButton" onclick="addGraph(\''+key+'\', \''+outputLabel+'\');"></i> ';
      element.appendChild(outputGraphDiv);

      var outputLabelDiv = document.createElement('span');
      outputLabelDiv.id = "output_"+key+"_label";
      outputLabelDiv.classList.add("outputLabel");
      element.appendChild(outputLabelDiv);

      var outputValueDiv = document.createElement('span');
      outputValueDiv.id = "output_"+key+"_value";
      outputValueDiv.classList.add("outputValue");
      element.appendChild(outputValueDiv);

      outputLabelDiv.innerText = outputLabel+":";
    } else {
      outputValueDiv = document.getElementById("output_"+key+"_value");
    }

    var canvasValueDiv = document.getElementById("canvas_"+key+"_value");
    if (canvasValueDiv != null) {
      canvasValueDiv.innerText = outputValue;
    }

    outputValueDiv.innerText = outputValue;

  });

  // go through and specifically output faults at the top
  Object.keys(data).sort().forEach(function(key) {
    if (key.startsWith("fault_")) { // only look at faults

      var element = document.getElementById("fault_"+key);

      if (data[key]["data"] == true) {
        if (element == null) {
          var element = document.createElement('span');
          element.id = "fault_"+key;
          faultsDiv.appendChild(element);
          element.innerHTML = '<i class="fas fa-exclamation-triangle"></i> '+data[key]["name"];
          element.classList.add("fault");
          element.classList.add("alert");
          element.classList.add("alert-danger");
        }
      } else {
        if (element != null) {
          element.parentNode.removeChild(element);
        }
      }

    }
  });

}



function downloadData() {
  // console.log(allOutputKeys);

  var lastSeenValue = {};

  var rows = [];
  var header = [];
  header.push("row");
  Object.keys(allOutputKeys).sort().forEach(function(key) {
    header.push(key)
  });
  rows.push(header);

  // make a copy so it doesn't keep changing
  data = outputHistory.slice();

  for (var i=0; i<data.length; i++) {
    rowInput = data[i];
    rowOutput = [i];

    // for each header key
    Object.keys(allOutputKeys).sort().forEach(function(key) {
      if (rowInput[key] == undefined) {
        if (lastSeenValue[key] != undefined) {
          rowOutput.push('"'+lastSeenValue[key]+'"');
        } else {
          rowOutput.push("");
        }
      } else {
        lastSeenValue[key] = rowInput[key]["data"];
        rowOutput.push('"'+rowInput[key]["data"]+'"');
      }
    });
    rows.push(rowOutput);
  }

  let csvContent = "data:text/csv;charset=utf-8,";

  rows.forEach(function(rowArray) {
      let row = rowArray.join(",");
      csvContent += row + "\r\n";
  });

  var encodedUri = encodeURI(csvContent);
  // window.open(encodedUri);

  var link = document.createElement("a");
  link.setAttribute("href", encodedUri);
  link.setAttribute("download", "rover_mems_export.csv");
  document.body.appendChild(link); // Required for FF

  link.click(); // This will download the data file

}


function downloadEeprom() {
  var a = document.createElement("a");
  document.body.appendChild(a);
  a.style = "display: none";
  var blob = new Blob([new Uint8Array(eeprom)], {type: "octet/stream"});
  url = window.URL.createObjectURL(blob);
  a.href = url;
  a.download = "ecu_eeprom_dump.bin";
  a.click();
  window.URL.revokeObjectURL(url);
}
