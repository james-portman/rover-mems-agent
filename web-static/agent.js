function agentConnect() {

  pollAgentDetection = false;

  getEcuVersionFromForm();

  document.getElementById("hideBeforeConnection").style.display = "block";
  document.getElementById("connectionSection").style.display = "none";
  document.getElementById("ecuChosen").innerText = ecuVersion;

  showDivs();

  refreshAgent();
}

var agentTimer;
function retryAgent(delay) {
  if (delay == null) { delay = 1000; }
  clearTimeout(agentTimer);
  agentTimer = setTimeout(function(){ refreshAgent(); }, delay);
}

function refreshAgent() {
  var userAction = "";

  if (userActionClearFaults == true) { userActionClearFaults = false; userAction = "clearfaults"; }
  else if (userActionStartTestRpmGauge == true) { userActionStartTestRpmGauge = false; userAction = "startTestRpmGauge"; }
  else if (userActionStartTestLambdaHeater == true) { userActionStartTestLambdaHeater = false; userAction = "startTestLambdaHeater"; }
  else if (userActionStopTestLambdaHeater == true) { userActionStopTestLambdaHeater = false; userAction = "stopTestLambdaHeater"; }
  else if (userActionStartTestACClutch == true) { userActionStartTestACClutch = false; userAction = "startTestACClutch"; }
  else if (userActionStopTestACClutch == true) { userActionStopTestACClutch = false; userAction = "stopTestACClutch"; }
  else if (userActionStartTestFuelPump == true) { userActionStartTestFuelPump = false; userAction = "startTestFuelPump"; }
  else if (userActionStopTestFuelPump == true) { userActionStopTestFuelPump = false; userAction = "stopTestFuelPump"; }
  else if (userActionStartTestFan1 == true) { userActionStartTestFan1 = false; userAction = "startTestFan1"; }
  else if (userActionStopTestFan1 == true) { userActionStopTestFan1 = false; userAction = "stopTestFan1"; }
  else if (userActionStartTestPurgeValve == true) { userActionStartTestPurgeValve = false; userAction = "startTestPurgeValve"; }
  else if (userActionStopTestPurgeValve == true) { userActionStopTestPurgeValve = false; userAction = "stopTestPurgeValve"; }
  else if (userActionIncreaseIdleSpeed == true) { userActionIncreaseIdleSpeed = false; userAction = "increaseIdleSpeed"; }
  else if (userActionDecreaseIdleSpeed == true) { userActionDecreaseIdleSpeed = false; userAction = "decreaseIdleSpeed"; }

  if (userAction != "") {
    fetch(agentAddress+'/command/'+userAction, {})
      .then(
        function(response) {
          if (response.status !== 200) {
            commandsAlert("Failed to ask the ECU for action"+userAction, "danger");
            debug('Looks like there was a problem ('+userAction+'). Agent status Code: ' +
              response.status);
          }
          retryAgent(1); return;
        }
      )
      .catch(function(err) {
        debug('Fetch Error from agent ('+userAction+'):-S', err);
        console.log('Fetch Error from agent ('+userAction+'):-S', err);
        retryAgent(1); return;
      });
  }


  fetch(agentAddress+'/api', {})
    .then(
      function(response) {
        if (response.status !== 200) {
          debug('Looks like there was a problem talking to the agent. Status Code: ' +
            response.status);
          setEcuConnected(false);
          retryAgent(); return;
        }
        response.json().then(function(data) {
          parseAgentResponse(data);
          retryAgent(); return;
        });
      }
    )
    .catch(function(err) {
      debug('Fetch Error from agent, please make sure it is running', err);
      console.log('Fetch Error from agent, please make sure it is running', err);
      setEcuConnected(false);
      retryAgent(); return;
    });

    retryAgent(5000);
    return;
}

var agentLastSeenFaults = [];

function parseAgentResponse(data) {
  spin();

  if (data.connected !== undefined) {
    setEcuConnected(data.connected);
  }

  if (data.ecuType !== undefined) {
    if (data.ecuType != ecuVersion) {
      debug("ECU set wrong in agent, telling it which to use")
      fetch(agentAddress+'/ecu/'+ecuVersion, {})
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
