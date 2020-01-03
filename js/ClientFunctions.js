
chatButton = document.getElementById("chatButton")
chatButton.addEventListener("click", sendChat)
chatInput = document.getElementById("chatInput")
chatbox = document.getElementById("chatbox")

var decoder = new TextDecoder("utf-8")
var toType = function (obj) {//from: http://javascriptweblog.wordpress.com/2011/08/08/fixing-the-javascript-typeof-operator/
  return ({}).toString.call(obj).match(/\s([a-zA-Z]+)/)[1].toLowerCase()
}

function encodeRequest(requestType, message) {//number, string
  if (requestType > 255) {
    console.log("ERROR: Request type number too large, upper bound is 255")
    return null
  }
  return String.fromCharCode(requestType) + message
}

function decodeResponse(event) {
  let data = new Uint8Array(event.data)
  return { "type": data[0], "message": decoder.decode(data.slice(1)) }
}

function addChatMessage(response) {
  console.log("RECEIVING CHAT")
  var chatMsgElement = document.createElement("div")
  chatMsgElement.classList.add("border")
  chatMsgElement.innerHTML = response.message
  chatbox.appendChild(chatMsgElement)
}

function sendChat() {
  console.log("SENDING CHAT")
  message = chatInput.value
  chatInput.value = ""
  encReq = encodeRequest(2, message)
  if (encReq != null)
    socket.send(encReq)
}
//https://developer.mozilla.org/en-US/docs/Web/API/MessageEvent
socket.onmessage = function (event) {
  response = decodeResponse(event)
  route(response)
};

socket.onopen = function (event) {
  //SOCKET-DEPENDENT SETUP
  //chatButton.onclick = sendChat

  console.log("Connected");
  //encReq = encodeRequest(2, "MESSAGE")

};

socket.onclose = function (event) {
  if (event.wasClean) {
    console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
  } else {
    console.log('[close] Connection died');
  }
};

socket.onerror = function (error) {
  console.log(`[error] ${error.message}`);
};