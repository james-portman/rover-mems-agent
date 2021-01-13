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

  if (data.ecuData !== undefined) {
    output = {};
    // console.log(data.ecuData);
    Object.keys(data.ecuData).forEach(function(key) {
      // console.log(key, data.ecuData[key]);
      output[key] = {"name": key, "data": data.ecuData[key]};
    });
    outputData(output);
  }


}
