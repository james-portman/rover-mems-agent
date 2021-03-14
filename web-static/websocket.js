var refreshSpeedSelect = document.getElementById("refreshSpeed");

function newwebsocket() {
  url = 'ws://localhost:8080/ws';
  window.c = new WebSocket(url);

  c.onmessage = function(msg){
    // $("#output").prepend((new Date())+ " <== "+msg.data+"\n")
    // console.log(msg)
    window.counter++;
    try {
      data = JSON.parse(msg.data);
      // console.log(data);
      // wsParse(data);
      parseAgentResponse(data);
      window.parsedCounter++;
    } catch(err) {
      commandsError("Failed parsing data from agent");
      console.log("failed parsing");
      console.log(err);
    }
    // console.log(data);

    // put some delay in before next data so we don't run 100%, for CPU/heat mostly
    // also fastest CPU puts out data is 50hz/20ms anyway
    setTimeout(function() { requestWsData(); }, refreshSpeedSelect.value);
  }

  c.onopen = function(){
    requestWsData();
  }

  c.onerror = function(){
    commandsError("Error while communicating with agent, trying to reconnect...");
    console.log("WEBSOCKET ERROR, will make a new connection");
    c.close()
  }

  c.onclose = function(){
    console.log("WEBSOCKET closed");
    commandsError("Connection to agent closed. You may need to restart it.");
    // setTimeout(function(){ newwebsocket(); }, 1000 );
  }

  // send = function(data){
  //   $("#output").prepend((new Date())+ " ==> "+data+"\n")
  //   c.send(data)
  // }
}

function requestWsData() {
  c.send(".");
}
//
// function wsParse(data) {
//   console.log(data);
//   for (key in data) {
//     var value = data[key];
//
//     if (key == "battery_voltage") {
//       value = value.toFixed(1);
//     } else if (key == "air_temp") {
//       value = value.toFixed(1);
//     } else if (key == "coolant_temp") {
//       value = value.toFixed(1);
//     } else if (key == "map_mbar") {
//       value = value.toFixed(0);
//     } else if (key == "barometric_pressure_bar") {
//       value = (value*1000).toFixed(0);
//     } else if (key == "tps_percent") {
//       value = value.toFixed(1);
//     } else if (key == "final_spark_advance") {
//       value = value.toFixed(2);
//     } else if (key == "last_ecu_data_received") {
//       window.lastEcuDataReceived = value;
//       continue;
//     }
//
//     // dont repeatedly set the same value
//     if (window.lastSeen[key] == value) {
//       continue;
//     }
//
//     window.lastSeen[key] = value;
//
//     // call a dynamic function called update_[key] if it exists (in per dash js probably)
//     // could have base ones then override with custom?
//     if (window["update_"+key] != undefined) {
//       window["update_"+key](value);
//     }
//
//   }
// }

newwebsocket();
