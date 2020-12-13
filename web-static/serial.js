// TODO: still check things like receive buffer, ordering of output/next command etc
// seems like bytes get lost when we are a bit busy
var writer = null;
var lastDataReceived = 0;

async function sendToEcu(bytes) {
  if (debugDataEnabledElement.checked) {
    debug("> "+bytes);
  }
  if (kLineEcu) {
    sendToKLineEcuRaw(bytes);
    // lastKLineByte = bytes[0];
    // gotKLineEcho = false;
  } else {

    writer = port.writable.getWriter();
    writer.write(Uint8Array.from(bytes));
    writer.releaseLock();
    writer = null;
  }
}
//
// async function readLoop() {
//   while (true) {
//     try {
//       const { value, done } = await reader.read();
//       if (value) {
//
//         for (var i=0; i<value.length; i++) {
//           if (kLineEcu) {
//             if (!doingSlowInit19 && !gotKLineEcho) {
//               // this should be our echo or something has gone wrong
//               if (value[i] == lastKLineByte) {
//                 // debug("Got our kline echo");
//                 // dataBuffer = dataBuffer.slice(1);
//                 gotKLineEcho = true;
//                 // break; // if it was a kline echo then don't even continue, just jump out
//                 continue; // if it was an echo then go check the next byte
//               }
//             }
//           }
//           // (k-line ecus) - actually got some data back here, not just our own echos
//           clearTimeout(kLineRetryTimer);
//           resetTimeout(5000);
//           // TODO: debug here, will show if we are wiping bytes or just ecu not responding or what
//           dataBuffer.push(value[i]);
//         }
//
//         if (debugDataEnabledElement.checked) {
//           debug("< "+value);
//         }
//
//         // might end up with nothing after throwing k-line echo away?
//         if (dataBuffer.length == 0) {
//           continue;
//         }
//
//         lastDataReceived = new Date().getTime();
//
//         if (doingSlowInit19) {
//           parseDataBufferSlowInit();
//         } else if (ecuVersion == "2J") {
//           parseDataBuffer2J();
//         } else if (ecuVersion == "3") {
//           parseDataBuffer3();
//         } else if (ecuVersion == "rc5") {
//           parseDataBufferRc5();
//         } else {
//           parseDataBuffer1x();
//         }
//       }
//       // await sleep(1);
//       if (done) {
//         debug('[readLoop] DONE', done);
//         reader.releaseLock();
//         break;
//       }
//     } catch (e) {
//       debug("readloop caught: "+e);
//       reader = port.readable.getReader();
//       await sleep(1);
//     }
//   }
// }


async function readOnce(parentName) {
  try {
    const { value, done } = await reader.read();
    if (!value) {
      debug("readOnce - no data/cancelled");
      return false;
    }

    var added = 0;

    for (var i=0; i<value.length; i++) {
      if (kLineEcu) {
        if (!doingSlowInit19 && !gotKLineEcho) {
          // this should be our echo or something has gone wrong
          if (value[i] == lastKLineByte) {
            // debug("Got our kline echo");
            // dataBuffer = dataBuffer.slice(1);
            gotKLineEcho = true;
            // debug("< k-line echo");
            return true;
            // break; // if it was a kline echo then don't even continue, just jump out
            // continue; // if it was an echo then go check the next byte
          }
        }
      }
      // (k-line ecus) - actually got some data back here, not just our own echos
      clearTimeout(kLineRetryTimer);
      resetTimeout(5000);
      // TODO: debug here, will show if we are wiping bytes or just ecu not responding or what
      // debug(parentName+"->readOnce added a byte to buffer");
      dataBuffer.push(value[i]);
      added++;
    }

    if (debugDataEnabledElement.checked) {
      debug("< "+value);
    }

    if (dataBuffer.length < added) {
      debug("Not all bytes got added somehow??");
    }

    // might end up with nothing after throwing k-line echo away?
    if (dataBuffer.length == 0) {
      return true;
    }

    lastDataReceived = new Date().getTime();

    // if (doingSlowInit19) {
    //   parseDataBufferSlowInit();
    // } else if (ecuVersion == "2J") {
    //   parseDataBuffer2J();
    // } else if (ecuVersion == "3") {
    //   parseDataBuffer3();
    // } else if (ecuVersion == "rc5") {
    //   parseDataBufferRc5();
    // } else {
    //   parseDataBuffer1x();
    // }


  } catch (e) {
    debug("readOnce caught: "+e);
    console.log("readOnce caught: "+e);
    reader = port.readable.getReader();
  }
  return true;

}

var baud_setting;
var parity_setting;
var kLineEcu;

async function connect() {
  pollAgentDetection = false;
  port = await navigator.serial.requestPort();

  document.getElementById("hideBeforeConnection").style.display = "block";
  document.getElementById("connectionSection").style.display = "none";
  document.getElementById("ecuChosen").innerText = ecuVersion;


  // per ecu settings
  switch (ecuVersion) {
    case "1.2":
    case "1.3-1.6":
      baud_setting = 9600;
      parity_setting = "none"
      kLineEcu = false;
      break;
    case "1.9":
      baud_setting = 9600;
      parity_setting = "none"
      kLineEcu = true;
      break;
    case "2J":
      baud_setting = 10400;
      parity_setting = "none"
      kLineEcu = true;
      break;

    case "3":
      baud_setting = 9600;
      parity_setting = "even";
      kLineEcu = true;
      break;

    case "rc5": // airbag module for JDM/MPI minis
      baud_setting = 2400;
      parity_setting = "none"
      kLineEcu = true;
      kLineByteSendGap = 5;
      break;

    default:
      debug("Unknown ECU selected");
      alert("Unknown ECU selected");
      break;
  }

  showDivs();

  debug("baud_setting:"+baud_setting);
  debug("parity_setting:"+parity_setting);

  // - Wait for the port to open.
  // I think baudRate is the newer version needed by chrome 86 ->
  // leave baudrate for a while for people not updated yet
  await port.open({ baudrate: baud_setting, baudRate: baud_setting, databits: 8, parity: parity_setting, stopbits: 1 });
  reader = port.readable.getReader();
  // readLoop();

  // act like we timed out, so (re)connect
  resetTimeout(1); // 1ms/now
}
