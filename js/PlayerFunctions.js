var charName = document.getElementById("name").innerHTML
let socket = new WebSocket("ws://192.168.0.12:8082/player/" + charName + "/socket");
socket.binaryType = 'arraybuffer';

//case switch instead of map because js is spaghetti
function route(response) {
  switch (response.type) {
    case 0:

      break;
    case 1:

      break;
    case 2:
      addChatMessage(response)
      break;
    default:
      console.log("ERROR: response type not found. ", response)
      break;
  }
}
