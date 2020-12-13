function xor_all_bytes(bytes) {
  current = 0
  for (var i=0; i<bytes.length; i++) {
    current = current ^ bytes[i];
  }
  return current;
}

function sum_checksum_all_bytes(bytes) {
  sum_checksum = 0;
  for (var i=0; i<dataBuffer.length-1; i++) {
    sum_checksum += dataBuffer[i];
  }
  return sum_checksum & 0xFF;
}

async function sleep(ms) {
  // if (ms < 4) { debug("asked to sleep for "+ms+" but minimum will be 4ms+"); }
  await new Promise(resolve => setTimeout(resolve, ms));
}

function busySleepMs(ms) {
  var start = new Date().getTime();
  var stop = start + ms;
  while (new Date().getTime() < stop) {
    // lol
  }
}

async function sleepUntil(timestampMs) {
  now = new Date().getTime()
  var sleepFor = timestampMs-now;
  // if (ms < 4) { debug("asked to sleep for "+ms+" but minimum will be 4ms+"); }
  await new Promise(resolve => setTimeout(resolve, sleepFor));
}


function showDivs() {
    // show divs
  switch (ecuVersion) {
    case "1.9":
      for (var i=0; i<mems19OnlyDivs.length; i++) {
        mems19OnlyDivs[i].style.display = "block";
      }
    case "1.2":
    case "1.3-1.6":
      for (var i=0; i<mems1xOnlyDivs.length; i++) {
        mems1xOnlyDivs[i].style.display = "block";
      }
      break;

    case "2J":
      for (var i=0; i<mems2jOnlyDivs.length; i++) {
        mems2jOnlyDivs[i].style.display = "block";
      }
      break;

    case "3":
      for (var i=0; i<mems3OnlyDivs.length; i++) {
        mems3OnlyDivs[i].style.display = "block";
      }
      break;

    case "rc5": // airbag module for JDM/MPI minis
      for (var i=0; i<rc5OnlyDivs.length; i++) {
        rc5OnlyDivs[i].style.display = "block";
      }
      break;
  }
}


function getEcuVersionFromForm() {
  var radios = document.getElementsByName('ecuVersion');
  for (var i = 0, length = radios.length; i < length; i++) {
    if (radios[i].checked) {
      ecuVersion = radios[i].value;
      break;
    }
  }
}
