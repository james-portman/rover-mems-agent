const debugPre = document.getElementById('debug');
const debugEnabledElement = document.getElementById('debugEnabled');
const debugDataEnabledElement = document.getElementById('debugDataEnabled');
const debugConsoleEnabledElement = document.getElementById('debugConsoleEnabled');
const ecuConnectedElement = document.getElementById('ecuConnected');
const spinnerElement = document.getElementById('spinner');
const ecuIdElement = document.getElementById('ecuId');
const carInfoElement = document.getElementById('carInfo');
const debugWrapperElement = document.getElementById('debugWrapper');

outputHistory = [];

ecuConnected = false;
ecuId = "";

knownEcus = {
  "9A,00,02,02": "MNE?????? Rover Mini SPI",
  "99,00,06,03": "MNE101351 Rover Mini SPI JDM, Air con, MEMS 1.6, 2 plugs, MPI flywheel/coil pack",
  "99,00,03,03": "MNE101170 Rover Mini SPI JDM, Air con, MEMS 1.6",
  "22,00,00,82": "MNE10077 Rover Metro 1.4 Auto MEMS 1.3",
  "AD,00,05,09": "MKC104052 Rover ???",
  "C7,00,06,CB": "MKC103111 MGF 1.8 MPi",
  "3A,00,00,14": "Rover Mini SPI JDM",
  "10,88,88,36": "MNE10039 Rover Metro? 1x 36 pin plug"
};

initAttempt = 0;

dataBuffer = new Array();

var commandsAlertTimer;
function commandsAlert(message, alert_type, timeout) {
  if (alert_type == undefined) { alert_type = "success"; }
  if (timeout == undefined) { timeout = 2000; }
  $('#commands_alert_placeholder').html('<div class="commands-alert alert alert-'+alert_type+' alert-dismissible fade show"><a href="#" class="close" data-dismiss="alert">×</a><span>'+message+'</span></div>')
  bootstrapAlertTimer = setTimeout(function(){
    $('.commands-alert').alert('close');
  }, timeout);
}

function updateEcuId() {
  ecuIdElement.textContent = "ECU ID: "+ecuId;
  if (knownEcus[ecuId] != undefined) {
    carInfoElement.textContent = knownEcus[ecuId];
  } else {
    carInfoElement.textContent = "Unknown car"
  }
}

function spin() {
  output = "";
  switch (spinnerElement.textContent) {
    case "-": output = "\\"; break;
    case "\\": output = "|"; break;
    case "|": output = "/"; break;
    case "/": output = "-"; break;
  }
  spinnerElement.textContent = output;
}

function log(message) {
  debug(message);
}
function debug(message) {
  // if (debugEnabledElement.checked) {
  //   var newText = message + '\n' + debugPre.textContent;
  //   if (newText.length > 1000) {
  //     newText = newText.slice(0, 999);
  //   }
  //   debugPre.textContent = newText;
  //
  //   if (debugConsoleEnabledElement.checked) {
  //     console.log(message);
  //   }
  // }
}

function setEcuConnected(connected) {
  if (connected) {
    ecuConnected = true;
    ecuConnectedElement.textContent = "Connected";
    ecuConnectedElement.classList.remove("disconnected");
    ecuConnectedElement.classList.add("connected");
  } else {
    ecuConnected = false;
    ecuConnectedElement.textContent = "Disconnected";
    ecuConnectedElement.classList.remove("connected");
    ecuConnectedElement.classList.add("disconnected");
  }
}


async function clickConnect() {
  try {
    debug("User clicked connect");

    getEcuVersionFromForm();

    await connect();
  } catch (e) {
    debug("clickConnect caught: "+e);
    console.log("clickConnect caught: "+e);
  }
}

var port = null;

function resetTimeout(ms) {
  clearTimeout(window.timeoutTimer);
  window.timeoutTimer = window.setTimeout(timedOut, ms); // reset timer
}

var firstrun = true;
async function timedOut() {
  if (firstrun) {
    firstrun = false;
  } else {
    try {
      await reader.cancel();
      if (writer != null) {
        await writer.releaseLock();
        writer = null;
      }
      await port.close();
      // I think baudRate is the newer version needed by chrome 86 ->
      // leave baudrate for a while for people not updated yet
      await port.open({ baudrate: baud_setting, baudRate: baud_setting, databits: 8, parity: parity_setting, stopbits: 1 });
      reader = port.readable.getReader();
    } catch(err) {
      console.log(err);
      resetTimeout(5000);
      return;
    }
  }

  clearTimeout(timerParseDataBufferRc5);
  clearTimeout(timerParseDataBuffer19SlowInit);
  clearTimeout(timerParseDataBuffer1x);
  clearTimeout(timerParseDataBuffer2j);
  clearTimeout(timerParseDataBuffer3);
  await sleep(500);

  log("Timed out with data in dataBuffer:");
  log(dataBuffer);
  resetTimeout(5000);
  dataBuffer = [];
  setEcuConnected(false);
  initAttempt++;

  switch (ecuVersion) {

    case "1.2":
    case "1.3-1.6":
      await initEcu131619();
      break;

    case "1.9":
      // slowInit19();
      if (initAttempt % 2 == 1) {
        await slowInit19();
      } else {
        await initEcu131619();
      }
      break;

    case "2J":
      await initEcu2J();
      break;

    case "3":
      await initEcu3();
      break;

    case "rc5":
      await initRc5();
      break;

    default:
      debug("Unknown ECU/module");
      break;
  }
  spin();
  ecuConnectedElement.textContent = "Connecting...";

}

setEcuConnected(false);
