function ecuTypeSelectChanged() {
  document.getElementById("connectionSection").style.display = "none";
  showDivs();
}

function runUserAction(userAction) {
  if (userAction != "") {
    fetch(agentAddress+'/command/'+userAction, {})
      .then(
        function(response) {
          if (response.status !== 200) {
            commandsAlert("Failed to ask the ECU for action"+userAction, "danger");
            debug('Looks like there was a problem ('+userAction+'). Agent status Code: ' +
              response.status);
          }
        }
      )
      .catch(function(err) {
        debug('Fetch Error from agent ('+userAction+'):-S', err);
        console.log('Fetch Error from agent ('+userAction+'):-S', err);
      });
  }

}

var agentLastSeenFaults = [];

var ecuVersionDiv = document.getElementById("ecuVersion");
var serialPortUI = document.getElementById("serialPort");

window.serialPorts = [];

function parseAgentResponse(data) {
  spin();

  if (data.connected !== undefined) {
    setEcuConnected(data.connected);
  }

  // console.log(data)

  if (data.ecuType !== undefined) {
    if (ecuVersionDiv.value == "") {
      ecuVersionDiv.value = data.ecuType;
      ecuVersionDiv.onchange();
    } else if (data.ecuType != ecuVersionDiv.value) {
      debug("ECU set wrong in agent, telling it which to use")
      fetch(agentAddress+'/ecu/'+ecuVersionDiv.value, {})
        .then(
          function(response) {
            if (response.status !== 200) {
              debug('Looks like there was a problem setting agent ECU type. Status Code: ' +
                response.status);
              return;
            }
          }
        )
        .catch(function(err) {
          debug('Fetch Error from agent (ecu type) :-S', err);
          console.log('Fetch Error from agent (ecu type) :-S', err);
        });
    }
  }

  if (data.faults !== undefined) {
    if (data.faults != agentLastSeenFaults) {
      output = {};
      for (var i=0; i < data.faults.length; i++) {
        faultText = data.faults[i];
        key = faultText.replace(" ", "_");
        output["fault_"+key] = {"name": faultText, "data": true};
      }
      clearFaults();
      outputData(output);
    }
  }

  if (data.alert !== undefined) {
    if (data.alert != "") {
      commandsAlert(data.alert);
    }
  }
  if (data.error !== undefined) {
    if (data.error != "") {
      commandsError(data.error);
    }
  }

  if (data.ecuData !== undefined) {
    output = {};
    // console.log(data.ecuData);
    Object.keys(data.ecuData).forEach(function(key) {
      // console.log(key, data.ecuData[key]);
      output[key] = {"name": key, "data": data.ecuData[key]};
    });
    outputData(output);
  }

  if (data.selectedSerialPort !== undefined) {

    if (serialPortUI.value == "") {
      if (data.selectedSerialPort != "") {
        // select the one the ECU said
        for (var i=0; i<serialPortUI.options.length; i++) {
          if (serialPortUI.options[i].value == data.selectedSerialPort) {
            if (serialPortUI.options[i].selected == false) {
              serialPortUI.options[i].selected = true;
            }
          }
        }
      }

    } else if (data.selectedSerialPort != serialPortUI.value) {
      debug("Serial port set wrong in agent, telling it which to use")
      console.log(serialPortUI);
      fetch(agentAddress+'/serialPort/'+serialPortUI.value, {})
        .then(
          function(response) {
            if (response.status !== 200) {
              debug('Looks like there was a problem setting agent serial port. Status Code: ' +
                response.status);
              return;
            }
          }
        )
        .catch(function(err) {
          debug('Fetch Error from agent (serial port setting) :-S', err);
          console.log('Fetch Error from agent (serial port setting) :-S', err);
        });
    }
  }

  if (data.serialPorts !== undefined) {
    if (data.serialPorts.length !== window.serialPorts.length || ! data.serialPorts.every(function(value, index) { return value === window.serialPorts[index]})) {

      console.log("serial ports changed");
      window.serialPorts = data.serialPorts;
      // add them to UI
      // remove old

      // delete any we no longer have
      for (var i=serialPortUI.options.length-1; i>=0; i--) {
        if (serialPortUI.options[i].value != "" && ! data.serialPorts.includes(serialPortUI.options[i].value)) {
          serialPortUI.remove(i);
        }
      }


      for (var i=0; i< data.serialPorts.length; i++) {

        var already = false;
        for (var j=0; j<serialPortUI.options.length; j++) {
          if (serialPortUI.options[j].value == data.serialPorts[i]) {
            already = true;
          }
        }
        if (already) { continue; }

        var option = document.createElement("option");
        option.text = data.serialPorts[i];
        option.value = data.serialPorts[i];
        if (data.selectedSerialPort == option.value) {
          option.selected = true;
        }
        serialPortUI.add(option);
      }

    }
  }

}
